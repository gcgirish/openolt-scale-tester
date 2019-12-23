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

type WorkFlow interface {
	ProvisionScheds() error
	ProvisionQueues() error
	ProvisionNNiTrapFlow() error
	ProvisionEapFlow() error
	ProvisionDhcpFlow() error
	ProvisionIgmpFlow() error
	ProvisionHsiaFlow() error
	// Add others here
}

func DeployWorkflow(wf WorkFlow) {
	wf.ProvisionScheds()
	wf.ProvisionQueues()
	wf.ProvisionNNiTrapFlow()
	wf.ProvisionEapFlow()
	wf.ProvisionDhcpFlow()
	wf.ProvisionIgmpFlow()
	wf.ProvisionHsiaFlow()
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

