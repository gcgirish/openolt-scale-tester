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
	"fmt"
	"github.com/opencord/voltha-lib-go/v2/pkg/log"
	"github.com/opencord/voltha-protos/v2/go/openolt"
)

func init() {
	_, _ = log.AddPackage(log.JSON, log.DebugLevel, nil)
}

const (
	vendorName = "ABCD"
)

var vendorSpecificId = 1000

func GenerateNextONUSerialNumber() *openolt.SerialNumber {

	vi := []byte(vendorName)

	vendorSpecificId += 1
	vs := []byte(fmt.Sprint(vendorSpecificId))
	// log.Infow("vendor-id-and-vendor-specific", log.Fields{"vi":vi, "vs":vs})
	sn := &openolt.SerialNumber{VendorId: vi, VendorSpecific: vs}
	// log.Infow("serial-num", log.Fields{"sn":sn})

	return sn
}
