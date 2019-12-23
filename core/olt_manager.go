/*
 * Copyright 2018-present Open Networking Foundation

 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at

 * http://www.apache.org/licenses/LICENSE-2.0

 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package core

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/cenkalti/backoff/v3"
	"github.com/opencord/openolt-scale-tester/config"
	"github.com/opencord/voltha-lib-go/v2/pkg/log"
	oop "github.com/opencord/voltha-protos/v2/go/openolt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"sync"
	"time"
)

const (
	ReasonOk = "OK"
)

type OnuDeviceKey struct {
	onuID    uint32
	ponInfID uint32
}

type OpenOltManager struct {
	ipPort        string
	deviceInfo    *oop.DeviceInfo
	OnuDeviceMap  map[OnuDeviceKey]*OnuDevice `json:"onuDeviceMap"`
	clientConn    *grpc.ClientConn
	openOltClient oop.OpenoltClient
	testConfig    *config.OpenOltScaleTesterConfig
	rsrMgr        *OpenOltResourceMgr
	lockRsrAlloc  sync.RWMutex
}

func init() {
	_, _ = log.AddPackage(log.JSON, log.DebugLevel, nil)
}

func NewOpenOltManager(ipPort string) *OpenOltManager {
	log.Infow("initialized openolt manager with ipPort", log.Fields{"ipPort": ipPort})
	return &OpenOltManager{
		ipPort:       ipPort,
		OnuDeviceMap: make(map[OnuDeviceKey]*OnuDevice),
		lockRsrAlloc: sync.RWMutex{},
	}
}

func (om *OpenOltManager) Start(testConfig *config.OpenOltScaleTesterConfig) error {
	var err error
	om.testConfig = testConfig

	// Establish gRPC connection with the device
	if om.clientConn, err = grpc.Dial(om.ipPort, grpc.WithInsecure(), grpc.WithBlock()); err != nil {
		log.Errorw("Failed to dial device", log.Fields{"ipPort": om.ipPort, "err": err})
		return err
	}
	om.openOltClient = oop.NewOpenoltClient(om.clientConn)

	// Populate Device Info
	if deviceInfo, err := om.populateDeviceInfo(); err != nil {
		log.Error("error fetching device info", log.Fields{"err": err, "deviceInfo": deviceInfo})
		return err
	}

	if om.rsrMgr = NewResourceMgr("ABCD", "127.0.0.1:2379", "etcd", "openolt", om.deviceInfo); om.rsrMgr == nil {
		log.Error("Error while instantiating resource manager")
		return errors.New("instantiating resource manager failed")
	}

	// Start reading indications
	go om.readIndications()

	// Provision OLT NNI Trap flows as needed by the Workflow
	if err = ProvisionNniTrapFlow(om.openOltClient, om.testConfig, om.rsrMgr); err != nil {
		log.Error("failed-to-add-nni-trap-flow", log.Fields{"err": err})
	}

	// Provision ONUs one by one
	go om.provisionONUs()

	return nil

}

func (om *OpenOltManager) populateDeviceInfo() (*oop.DeviceInfo, error) {
	var err error

	if om.deviceInfo, err = om.openOltClient.GetDeviceInfo(context.Background(), new(oop.Empty)); err != nil {
		log.Errorw("Failed to fetch device info", log.Fields{"err": err})
		return nil, err
	}

	if om.deviceInfo == nil {
		log.Errorw("Device info is nil", log.Fields{})
		return nil, errors.New("failed to get device info from OLT")
	}

	log.Debugw("Fetched device info", log.Fields{"deviceInfo": om.deviceInfo})

	return om.deviceInfo, nil
}

func (om *OpenOltManager) provisionONUs() {
	var numOfONUsPerPon uint
	var i, j, onuID uint32
	var err error
	numOfONUsPerPon = om.testConfig.NumOfOnu / uint(om.deviceInfo.PonPorts)
	if oddONUs := om.testConfig.NumOfOnu % uint(om.deviceInfo.PonPorts); oddONUs > 0 {
		log.Warnw("Odd number ONUs left out of provisioning", log.Fields{"oddONUs": oddONUs})
	}
	for i = 0; i < om.deviceInfo.PonPorts; i++ {
		for j = 0; j < uint32(numOfONUsPerPon); j++ {
			// TODO: More work with ONU provisioning
			om.lockRsrAlloc.Lock()
			sn := GenerateNextONUSerialNumber()
			log.Debugw("provisioning onu", log.Fields{"onuID": j, "ponPort": i, "serialNum": sn})
			if onuID, err = om.rsrMgr.GetONUID(i); err != nil {
				log.Errorw("error getting onu id", log.Fields{"err": err})
				om.lockRsrAlloc.Unlock()
				continue
			}
			om.activateONU(i, onuID, sn, om.stringifySerialNumber(sn))
			om.lockRsrAlloc.Unlock()
		}
	}

	// TODO: We need to dump the results at the end. But below json marshall does not work
	// We will need custom Marshal function.
	/*
		e, err := json.Marshal(om)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(e))
	*/
}

func (om *OpenOltManager) activateONU(intfID uint32, onuID uint32, serialNum *oop.SerialNumber, serialNumber string) {
	log.Debugw("activate-onu", log.Fields{"intfID": intfID, "onuID": onuID, "serialNum": serialNum, "serialNumber": serialNumber})
	// TODO: need resource manager
	var pir uint32 = 1000000
	var onuDevice OnuDevice
	var err error
	onuDeviceKey := OnuDeviceKey{onuID: onuID, ponInfID: intfID}
	Onu := oop.Onu{IntfId: intfID, OnuId: onuID, SerialNumber: serialNum, Pir: pir}
	now := time.Now()
	nanos := now.UnixNano()
	milliStart := nanos / 1000000
	onuDevice.ProvisionStartTime = time.Unix(0, nanos)
	if _, err = om.openOltClient.ActivateOnu(context.Background(), &Onu); err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.AlreadyExists {
			log.Debug("ONU activation is in progress", log.Fields{"SerialNumber": serialNumber})
		} else {
			log.Errorw("activate-onu-failed", log.Fields{"Onu": Onu, "err ": err})
		}
	} else {
		log.Infow("activated-onu", log.Fields{"SerialNumber": serialNumber})
	}
	nanos = now.UnixNano()
	milliEnd := nanos / 1000000
	onuDevice.ProvisionEndTime = time.Unix(0, nanos)
	onuDevice.ProvisionDurationInMs = milliEnd - milliStart
	onuDevice.Reason = ReasonOk
	if err != nil {
		onuDevice.Reason = err.Error()
	}

	om.OnuDeviceMap[onuDeviceKey] = &onuDevice
}

func (om *OpenOltManager) stringifySerialNumber(serialNum *oop.SerialNumber) string {
	if serialNum != nil {
		return string(serialNum.VendorId) + om.stringifyVendorSpecific(serialNum.VendorSpecific)
	}
	return ""
}

func (om *OpenOltManager) stringifyVendorSpecific(vendorSpecific []byte) string {
	tmp := fmt.Sprintf("%x", (uint32(vendorSpecific[0])>>4)&0x0f) +
		fmt.Sprintf("%x", uint32(vendorSpecific[0]&0x0f)) +
		fmt.Sprintf("%x", (uint32(vendorSpecific[1])>>4)&0x0f) +
		fmt.Sprintf("%x", (uint32(vendorSpecific[1]))&0x0f) +
		fmt.Sprintf("%x", (uint32(vendorSpecific[2])>>4)&0x0f) +
		fmt.Sprintf("%x", (uint32(vendorSpecific[2]))&0x0f) +
		fmt.Sprintf("%x", (uint32(vendorSpecific[3])>>4)&0x0f) +
		fmt.Sprintf("%x", (uint32(vendorSpecific[3]))&0x0f)
	return tmp
}

// readIndications to read the indications from the OLT device
func (om *OpenOltManager) readIndications() {
	defer log.Errorw("Indications ended", log.Fields{})
	indications, err := om.openOltClient.EnableIndication(context.Background(), new(oop.Empty))
	if err != nil {
		log.Errorw("Failed to read indications", log.Fields{"err": err})
		return
	}
	if indications == nil {
		log.Errorw("Indications is nil", log.Fields{})
		return
	}

	// Create an exponential backoff around re-enabling indications. The
	// maximum elapsed time for the back off is set to 0 so that we will
	// continue to retry. The max interval defaults to 1m, but is set
	// here for code clarity
	indicationBackoff := backoff.NewExponentialBackOff()
	indicationBackoff.MaxElapsedTime = 0
	indicationBackoff.MaxInterval = 1 * time.Minute
	for {
		indication, err := indications.Recv()
		if err == io.EOF {
			log.Infow("EOF for  indications", log.Fields{"err": err})
			// Use an exponential back off to prevent getting into a tight loop
			duration := indicationBackoff.NextBackOff()
			if duration == backoff.Stop {
				// If we reach a maximum then warn and reset the backoff
				// timer and keep attempting.
				log.Warnw("Maximum indication backoff reached, resetting backoff timer",
					log.Fields{"max_indication_backoff": indicationBackoff.MaxElapsedTime})
				indicationBackoff.Reset()
			}
			time.Sleep(indicationBackoff.NextBackOff())
			indications, err = om.openOltClient.EnableIndication(context.Background(), new(oop.Empty))
			if err != nil {
				log.Errorw("Failed to read indications", log.Fields{"err": err})
				return
			}
			continue
		}
		if err != nil {
			log.Infow("Failed to read from indications", log.Fields{"err": err})
			break
		}
		// Reset backoff if we have a successful receive
		indicationBackoff.Reset()
		om.handleIndication(indication)

	}
}

func (om *OpenOltManager) handleIndication(indication *oop.Indication) {
	switch indication.Data.(type) {
	case *oop.Indication_OltInd:
		log.Info("received olt indication")
	case *oop.Indication_IntfInd:
		intfInd := indication.GetIntfInd()
		log.Infow("Received interface indication ", log.Fields{"InterfaceInd": intfInd})
	case *oop.Indication_IntfOperInd:
		intfOperInd := indication.GetIntfOperInd()
		if intfOperInd.GetType() == "nni" {
			log.Info("received interface oper indication for nni port")
		} else if intfOperInd.GetType() == "pon" {
			log.Info("received interface oper indication for pon port")
		}
		/*
			case *oop.Indication_OnuDiscInd:
				onuDiscInd := indication.GetOnuDiscInd()
				log.Infow("Received Onu discovery indication ", log.Fields{"OnuDiscInd": onuDiscInd})
		*/
	case *oop.Indication_OnuInd:
		onuInd := indication.GetOnuInd()
		log.Infow("Received Onu indication ", log.Fields{"OnuInd": onuInd})
	case *oop.Indication_OmciInd:
		omciInd := indication.GetOmciInd()
		log.Debugw("Received Omci indication ", log.Fields{"IntfId": omciInd.IntfId, "OnuId": omciInd.OnuId, "pkt": hex.EncodeToString(omciInd.Pkt)})
	case *oop.Indication_PktInd:
		pktInd := indication.GetPktInd()
		log.Infow("Received pakcet indication ", log.Fields{"PktInd": pktInd})
		/*
				case *oop.Indication_PortStats:
				portStats := indication.GetPortStats()
				log.Infow("Received port stats", log.Fields{"portStats": portStats})
			case *oop.Indication_FlowStats:
				flowStats := indication.GetFlowStats()
				log.Infow("Received flow stats", log.Fields{"FlowStats": flowStats})
		*/
	case *oop.Indication_AlarmInd:
		alarmInd := indication.GetAlarmInd()
		log.Infow("Received alarm indication ", log.Fields{"AlarmInd": alarmInd})
	}
}
