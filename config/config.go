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

package config

import (
	"flag"
	"fmt"
	"github.com/opencord/openolt-scale-tester/core"
	"github.com/opencord/voltha-lib-go/v2/pkg/log"
	"strconv"
)

// Open OLT default constants
const (
	defaultOpenOltAgentIp            = "10.90.0.114"
	defaultOpenOltAgentPort          = 9191
	defaultNumOfOnu   				 = 128
	defaultNumOfSubscribersPerOnu    = 1
	defaultWorkFlowName  			 = "ATT"
	defaultTimeIntervalBetweenSubs   = 5 // in seconds
)

// OpenOltScaleTester represents the set of configurations used by the read-write adaptercore service
type OpenOltScaleTester struct {
	// Command line parameters
	OpenOltAgentAddress 	string
	OpenOltAgentIP			string
	OpenOltAgentPort        uint32
	NumOfOnu	        	uint32
	SubscribersPerOnu   	uint32
	WorkflowName   			string
	TimeIntervalBetweenSubs uint32 // in seconds
	OpenOltManager          *core.OpenOltManager
}

func init() {
	_, _ = log.AddPackage(log.JSON, log.WarnLevel, nil)
}

// NewOpenOltScaleTester returns a new RWCore config
func NewOpenOltScaleTester() *OpenOltScaleTester {
	var OpenOltScaleTester = OpenOltScaleTester{ // Default values
		OpenOltAgentAddress:        defaultOpenOltAgentIp + ":" + strconv.Itoa(defaultOpenOltAgentPort),
		NumOfOnu:   				defaultNumOfOnu,
		SubscribersPerOnu:   		defaultNumOfSubscribersPerOnu,
		WorkflowName:  			    defaultWorkFlowName,
		TimeIntervalBetweenSubs:    defaultTimeIntervalBetweenSubs,

	}
	return &OpenOltScaleTester
}

// ParseCommandArguments parses the arguments for OpenOltScale Tester
func (st *OpenOltScaleTester) ParseCommandArguments() {

	help := fmt.Sprintf("OpenOLT Agent IP Address")
	flag.StringVar(&(st.OpenOltAgentIP), "openolt_agent_ip_address", defaultOpenOltAgentIp, help)

	help = fmt.Sprintf("OpenOLT Agent gRPC port")
	flag.IntVar(&(st.OpenOltAgentPort), "openolt_agent_port", defaultOpenOltAgentPort, help)

	help = fmt.Sprintf("Number of ONU")
	flag.IntVar(&(st.NumOfOnu), "num_of_onu", defaultNumOfOnu, help)

	help = fmt.Sprintf("Kafka - Cluster messaging port")
	flag.IntVar(&(st.SubscribersPerOnu), "subscribers_per_onu", defaultNumOfSubscribersPerOnu, help)

	help = fmt.Sprintf("Workflow name")
	flag.StringVar(&(st.WorkflowName), "workflow_name", defaultWorkFlowName, help)

	help = fmt.Sprintf("Time Interval Between provisioning each subscriber")
	flag.IntVar(&(st.TimeIntervalBetweenSubs), "time_interval_between_subs", defaultTimeIntervalBetweenSubs, help)

	flag.Parse()

}
