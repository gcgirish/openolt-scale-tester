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
	"github.com/opencord/openolt-scale-tester/config"
	oop "github.com/opencord/voltha-protos/v2/go/openolt"
)

type Subscriber struct {
	SubscriberName string   `json:"subscriberName"`
	UniPortNo      uint32   `json:"uniPortNo"`
	Ctag           uint32   `json:"ctag"`
	Stag           uint32   `json:"stag"`
	GemPortIDs     []uint32 `json:"gemPortIds"`
	AllocIDs       []uint32 `json:"allocIds"`
	FlowIDs        []uint32 `json:"flowIds"`
	openOltClient  oop.OpenoltClient
	testConfig     *config.OpenOltScaleTesterConfig
	rsrMgr         *OpenOltResourceMgr
}

func (subs *Subscriber) Start(onuCh chan bool) {

	onuCh <- true
}
