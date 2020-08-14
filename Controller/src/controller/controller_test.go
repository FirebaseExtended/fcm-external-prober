/*
 * Copyright 2020 Google LLC
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

package controller

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/FirebaseExtended/fcm-external-prober/Probe/src/utils"
	"github.com/golang/protobuf/proto"
)

type fakeControllerLogger struct {
}

func (f *fakeControllerLogger) LogFatal(desc string)                       {}
func (f *fakeControllerLogger) LogFatalf(desc string, args ...interface{}) {}
func (f *fakeControllerLogger) LogError(desc string)                       {}
func (f *fakeControllerLogger) LogErrorf(desc string, args ...interface{}) {}

func TestGetPossibleZones(t *testing.T) {
	testStrings := []string{"REGION-a\nREGION-b\nREGION2-a\nREGION2-B",
		"INFORMATION\nMIN_CPU\nOTHER_INFORMATION", "INFORMATION"}
	maker = utils.NewFakeCommandMaker(testStrings, []bool{false, false, false}, false)
	clock = utils.NewFakeClock([]time.Time{time.Unix(0, 0)}, false)
	vms = make(map[string]*regionalVM)
	cfg, err := getTestConfig("testConfig.txt")
	config = cfg
	if err != nil {
		t.Logf("TestGetPossibleZones: unable to parse test configuration file: %v", err)
		t.FailNow()
	}

	getPossibleZones()

	if len(vms) != 1 {
		t.Logf("TestGetPossibleZones: incorrect number of resulting zones: actual: %d, expected: %d", len(vms), 1)
		t.FailNow()
	}
	for _, v := range config.Probes.Probe {
		if vms[v.GetRegion()] != nil {
			t.Logf("TestGetPossibleZones: incorrect value in resulting vm object")
			t.Fail()
		}
	}
}

func TestController(t *testing.T) {
	testStrings := []string{"REGION-a\nREGION-b\nREGION2-a\nREGION2-b\nREGION3-a",
		"INFORMATION\nMIN_CPU\nOTHER_INFORMATION", "MIN_CPU", "INFORMATION", "", "", "", ""}
	maker := utils.NewFakeCommandMaker(testStrings, make([]bool, 8), false)
	timer := utils.NewFakeClock([]time.Time{time.Unix(0, 0), time.Unix(1, 0), time.Unix(1, 0), time.Unix(1, 0), time.Unix(2, 0)}, false)
	cfg, err := getTestConfig("testConfig.txt")
	if err != nil {
		t.Logf("TestGetPossibleZones: unable to parse test configuration file: %v", err)
		t.FailNow()
	}
	ctrl := NewController(cfg, maker, timer, new(fakeControllerLogger))
	stopping = true

	ctrl.InitProbes()
	ctrl.MonitorProbes()

	if stoppedVMs != 2 {
		t.Logf("TestControl: VMs not stopped correctly")
		t.Fail()
	}
	if vms["REGION-a"].state != stopped || vms["REGION2-a"].state != stopped {
		t.Logf("TestControl: Incorrect VM stopped")
		t.Fail()
	}
}

func getTestConfig(filename string) (*ControllerConfig, error) {
	c, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg := new(ControllerConfig)
	err = proto.UnmarshalText(string(c), cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
