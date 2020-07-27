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
	"log"
	"testing"

	"github.com/FirebaseExtended/fcm-external-prober/Probe/src/utils"
)

func TestGetPossibleZones(t *testing.T) {
	testStrings := []string{"REGION-a\nREGION-b\nREGION2-a\nREGION2-B",
		"INFORMATION\nMIN_CPU\nOTHER_INFORMATION", "INFORMATION"}
	maker = utils.NewFakeCommandMaker(testStrings, []bool{false, false, false}, false)
	cfg, err := ioutil.ReadFile("testConfig.txt")
	if err != nil {
		t.Log("TestGetPossibleZones: unable to parse configuration file")
		t.FailNow()
	}
	ctrl := NewController(string(cfg), maker)

	ctrl.getPossibleZones()

	if len(ctrl.vms) != 1 {
		t.Logf("TestGetPossibleZones: incorrect number of resulting zones: actual: %d, expeted: %d", len(ctrl.vms), 1)
		t.FailNow()
	}
	for _, v := range ctrl.config.Probes.Probe {
		if ctrl.vms[v.GetRegion()] != nil {
			t.Logf("TestGetPossibleZones: incorrect value in resulting vm object")
			t.Fail()
		}
	}
}

func TestStartVMs(t *testing.T) {
	testStrings := []string{"REGION-a\nREGION-b\nREGION2-a\nREGION2-b\nREGION3-a",
		"INFORMATION\nMIN_CPU\nOTHER_INFORMATION", "MIN_CPU", "INFORMATION", "", ""}
	maker = utils.NewFakeCommandMaker(testStrings, make([]bool, 6), false)
	cfg, err := ioutil.ReadFile("testConfig.txt")
	if err != nil {
		t.Log("TestStartVMs: unable to parse configuration file")
		t.FailNow()
	}
	ctrl := NewController(string(cfg), maker)

	ctrl.StartVMs()

	for _, p := range ctrl.config.Probes.Probe {
		if !ctrl.vms[p.GetRegion() + "-a"].active {
			t.Log("TestStartVMs: zonal VM for which a probe exists is not active")
			t.Fail()
		}
		ctrl.vms[p.GetRegion() + "-a"].active = false
	}

	for _, v := range ctrl.vms {
		if v.active {
			t.Log("TestStartVMs: zonal VM for which a probe does not exist is active")
		}
	}
}

func TestCommands(t *testing.T) {
	maker := new(utils.CmdMaker)
	str, err := maker.Command("gcloud", "compute", "ssh", "us-east4-a", "--zone", "us-east4-a", "--command", "echo hello").Output()
	if err != nil {
		log.Print("failed to ssh")
		t.Fail()
	}
	log.Print(string(str))
}