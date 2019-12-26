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
	"errors"
	"github.com/opencord/openolt-scale-tester/config"
	"github.com/opencord/openolt-scale-tester/core"
	"github.com/opencord/voltha-lib-go/v2/pkg/log"
	oop "github.com/opencord/voltha-protos/v2/go/openolt"
)


func init() {
	_, _ = log.AddPackage(log.JSON, log.DebugLevel, nil)
}

type WorkFlow interface {
	ProvisionScheds(subs *core.Subscriber) error
	ProvisionQueues(subs *core.Subscriber) error
	ProvisionEapFlow(subs *core.Subscriber) error
	ProvisionDhcpFlow(subs *core.Subscriber) error
	ProvisionIgmpFlow(subs *core.Subscriber) error
	ProvisionHsiaFlow(subs *core.Subscriber) error
	// Add others here
}

func DeployWorkflow(subs *core.Subscriber) {
	var wf = getWorkFlow(subs)
	if wf == nil {
		log.Error("could-not-find-workflow")
		return
	}
	// TODO: Catch and log errors
	_ = wf.ProvisionScheds(subs)
	_ = wf.ProvisionQueues(subs)
	_ = wf.ProvisionEapFlow(subs)
	_ = wf.ProvisionDhcpFlow(subs)
	_ = wf.ProvisionIgmpFlow(subs)
	_ = wf.ProvisionHsiaFlow(subs)
}

func getWorkFlow(subs *core.Subscriber) WorkFlow {
	switch subs.TestConfig.WorkflowName {
	case "ATT":
		log.Info("chosen-att-workflow")
		return AttWorkFlow{}
	default:
		log.Errorw("operator-workflow-not-supported-yet", log.Fields{"workflowName": subs.TestConfig.WorkflowName})
	}
	return nil
}

func ProvisionNniTrapFlow(oo oop.OpenoltClient, config *config.OpenOltScaleTesterConfig, rsrMgr *core.OpenOltResourceMgr) error {
	switch config.WorkflowName {
	case "ATT":
		if err := ProvisionAttNniTrapFlow(oo, config, rsrMgr); err != nil {
			log.Error("error-installing-flow", log.Fields{"err": err})
			return err
		}
	default:
		log.Errorw("operator-workflow-not-supported-yet", log.Fields{"workflowName": config.WorkflowName})
		return errors.New("workflow-not-supported")
	}
	return nil
}

