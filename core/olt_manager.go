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
	backoff "github.com/cenkalti/backoff/v3"
	"github.com/opencord/openolt-scale-tester/config"
	"github.com/opencord/voltha-lib-go/v2/pkg/log"
	oop "github.com/opencord/voltha-protos/v2/go/openolt"
	"google.golang.org/grpc"
	"io"
	"time"
)

type OnuDeviceKey struct {
	onuID    uint16
	ponInfID uint16
}

type OpenOltManager struct {
	ipPort        string
	deviceInfo    *oop.DeviceInfo
	onuDeviceMap  map[OnuDeviceKey]*OnuDevice
	clientConn    *grpc.ClientConn
	openOltClient oop.OpenoltClient
	testConfig    *config.OpenOltScaleTesterConfig
}

func init() {
	_, _ = log.AddPackage(log.JSON, log.DebugLevel, nil)
}

func NewOpenOltManager(ipPort string) *OpenOltManager {
	log.Infow("initialized openolt manager with ipPort", log.Fields{"ipPort": ipPort})
	return &OpenOltManager{
		ipPort:       ipPort,
		onuDeviceMap: make(map[OnuDeviceKey]*OnuDevice),
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

	// Start reading indications
	go om.readIndications()

	go om.provisionONUs()

	return nil

}

///// Utility functions

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
	numOfONUsPerPon = om.testConfig.NumOfOnu / uint(om.deviceInfo.PonPorts)
	oddONUs := om.testConfig.NumOfOnu % uint(om.deviceInfo.PonPorts)
	log.Warnw("Odd number ONUs left out of provisioning", log.Fields{"oddONUs": oddONUs})
	for i := 0; i < int(om.deviceInfo.PonPorts); i++ {
		for j := 0; j < int(numOfONUsPerPon); j++ {
			// TODO: More work with ONU provisioning
			log.Debugw("provisioning onu", log.Fields{"onuID": j, "ponPort": i})

		}
	}
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
	case *oop.Indication_OnuDiscInd:
		onuDiscInd := indication.GetOnuDiscInd()
		log.Infow("Received Onu discovery indication ", log.Fields{"OnuDiscInd": onuDiscInd})
	case *oop.Indication_OnuInd:
		onuInd := indication.GetOnuInd()
		log.Infow("Received Onu indication ", log.Fields{"OnuInd": onuInd})
	case *oop.Indication_OmciInd:
		omciInd := indication.GetOmciInd()
		log.Debugw("Received Omci indication ", log.Fields{"IntfId": omciInd.IntfId, "OnuId": omciInd.OnuId, "pkt": hex.EncodeToString(omciInd.Pkt)})
	case *oop.Indication_PktInd:
		pktInd := indication.GetPktInd()
		log.Infow("Received pakcet indication ", log.Fields{"PktInd": pktInd})
	case *oop.Indication_PortStats:
		portStats := indication.GetPortStats()
		log.Infow("Received port stats", log.Fields{"portStats": portStats})
	case *oop.Indication_FlowStats:
		flowStats := indication.GetFlowStats()
		log.Infow("Received flow stats", log.Fields{"FlowStats": flowStats})
	case *oop.Indication_AlarmInd:
		alarmInd := indication.GetAlarmInd()
		log.Infow("Received alarm indication ", log.Fields{"AlarmInd": alarmInd})
	}
}
