/*
 *  Copyright 2020 Google LLC
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package probe

import (
	"testing"
	"utils"
)

func TestFindDevice(t *testing.T) {
	maker = utils.NewFakeCommandMaker([]string{"TEST_DEVICE_1\nTEST_DEVICE_2\nTEST_DEVICE_3"}, []bool{false}, false)

	str, err := findDevice()

	if err != nil {
		t.Log("TestFindDevice: error on valid input")
		t.Fail()
	}
	if str != "TEST_DEVICE_1" {
		t.Log("TestFindDevice: incorrect device found")
		t.Fail()
	}
}
