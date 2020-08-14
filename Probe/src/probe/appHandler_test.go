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
	"time"

	"github.com/FirebaseExtended/fcm-external-prober/Probe/src/utils"
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

func TestConvertTimeExpected(t *testing.T) {
	tim, err := convertTime("123.456000")
	if err != nil {
		t.Logf("TestConvertTimeExpected: returned error on valid input: %v", err)
		t.FailNow()
	}
	expected := time.Unix(123, 456000000)
	if !tim.Equal(expected) {
		t.Logf("TestConvertTimeExpected: incorrect time returned: actual: %v, expected %v", tim, expected)
		t.FailNow()
	}
}

func TestConvertTimeError(t *testing.T) {
	_, err := convertTime("1.1.1.1")
	if err == nil {
		t.Logf("TestConvertTimeExpected: no error returned on invalid input")
		t.Fail()
	}
}

func TestFindTimeOffset(t *testing.T) {
	clock = utils.NewFakeClock([]time.Time{time.Unix(0, 0), time.Unix(1, 0)}, false)
	maker = utils.NewFakeCommandMaker([]string{"000.500000"}, []bool{false}, false)
	offset, err := findTimeOffset()
	if err != nil {
		t.Logf("TestFindTimeOffset: error returned on valid input: %v", err)
		t.FailNow()
	}
	if offset != 0 {
		t.Logf("TestFindTimeOffset: incorrect time offset: actual: %d, expected: 0", offset)
		t.Fail()
	}
}
