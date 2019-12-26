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

package workflow

import (
	"github.com/opencord/openolt-scale-tester/config"
	"github.com/opencord/openolt-scale-tester/core"
	"github.com/opencord/voltha-lib-go/v2/pkg/log"
	"github.com/opencord/voltha-lib-go/v2/pkg/ponresourcemanager"
	oop "github.com/opencord/voltha-protos/v2/go/openolt"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func init() {
	_, _ = log.AddPackage(log.JSON, log.DebugLevel, nil)
}

// A dummy struct to comply with the WorkFlow interface.
type AttWorkFlow struct {
}

func ProvisionAttNniTrapFlow(oo oop.OpenoltClient, config *config.OpenOltScaleTesterConfig, rsrMgr *core.OpenOltResourceMgr) error {
	var flowID []uint32
	var err error

	if flowID, err = rsrMgr.ResourceMgrs[uint32(config.NniIntfID)].GetResourceID(uint32(config.NniIntfID),
		ponresourcemanager.FLOW_ID, 1); err != nil {
		return err
	}

	flowClassifier := &oop.Classifier{EthType: 2048, IpProto: 17, SrcPort: 67, DstPort: 68, PktTagType: "double_tag"}
	actionCmd := &oop.ActionCmd{TrapToHost: true}
	actionInfo := &oop.Action{Cmd: actionCmd}

	flow := oop.Flow{AccessIntfId: -1, OnuId: -1, UniId: -1, FlowId: flowID[0],
		FlowType: "downstream", AllocId: -1, GemportId: -1,
		Classifier: flowClassifier, Action: actionInfo,
		Priority: 1000, PortNo: 65536}

	_, err = oo.FlowAdd(context.Background(), &flow)

	st, _ := status.FromError(err)
	if st.Code() == codes.AlreadyExists {
		log.Debugw("Flow already exists", log.Fields{"err": err, "deviceFlow": flow})
		return nil
	}

	if err != nil {
		log.Errorw("Failed to Add flow to device", log.Fields{"err": err, "deviceFlow": flow})
		rsrMgr.ResourceMgrs[uint32(config.NniIntfID)].FreeResourceID(uint32(config.NniIntfID), ponresourcemanager.FLOW_ID, flowID)
		return err
	}
	log.Debugw("Flow added to device successfully ", log.Fields{"flow": flow})

	return nil
}

func (att AttWorkFlow) ProvisionScheds(subs *core.Subscriber) error {
	log.Info("provisioning-scheds")
	return nil
}

func (att AttWorkFlow) ProvisionQueues(subs *core.Subscriber) error {
	log.Info("provisioning-queues")
	return nil
}

func (att AttWorkFlow) ProvisionEapFlow(subs *core.Subscriber) error {
	log.Info("provisioning-eap")
	return nil
}

func (att AttWorkFlow) ProvisionDhcpFlow(subs *core.Subscriber) error {
	log.Info("provisioning-dhcp")
	return nil
}

func (att AttWorkFlow) ProvisionIgmpFlow(subs *core.Subscriber) error {
	log.Info("att-workflow-does-not-support-igmp-yet--nothing-to-do")
	return nil
}

func (att AttWorkFlow) ProvisionHsiaFlow(subs *core.Subscriber) error {
	log.Info("provisioning-hsia")
	return nil
}
