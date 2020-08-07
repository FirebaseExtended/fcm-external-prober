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
	"testing"
	"time"

	"github.com/FirebaseExtended/fcm-external-prober/Probe/src/utils"
)

func TestRegisterNotFound(t *testing.T) {
	server := &CommunicatorServer{}
	src := "DOES_NOT_EXIST"
	req := &RegisterRequest{Source: &src}

	_, err := server.Register(nil, req)

	if err == nil {
		t.Log("TestRegisterNotFound: No error returned given invalid source input")
		t.Fail()
	}
}

func TestRegisterExpected(t *testing.T) {
	server := &CommunicatorServer{}
	src := "REGION"
	req := &RegisterRequest{Source: &src}
	clock = utils.NewFakeClock([]time.Time{time.Unix(0, 0), time.Unix(1, 0)}, false)
	testVM := newRegionalVM("", "")
	vms = map[string]*regionalVM{"REGION": testVM}
	cfg, err := getTestConfig("testConfig.txt")
	if err != nil {
		t.Log("TestRegisterExpected: unable to parse test configuration file")
		t.Fail()
	}
	config = cfg

	res, err := server.Register(nil, req)
	if err != nil {
		t.Log("TestRegisterExpected: Error returned on valid source input")
		t.FailNow()
	}

	if len(res.GetProbes().GetProbe()) != len(testVM.probes) {
		t.Log("TestRegisterExpected: Incorrect number of probes returned from Register")
		t.Fail()
	}
	if res.GetAccount() != config.GetAccount() {
		t.Log("TestRegisterExpected: Incorrect account information returned from Register")
		t.Fail()
	}
	if res.GetPingConfig() != config.GetPingConfig() {
		t.Log("TestRegisterExpected: Incorrect ping configuration returned from Register")
		t.Fail()
	}
	if testVM.state != idle || !testVM.lastPing.Equal(time.Unix(1, 0)) {
		t.Log("TestRegisterExpected: VM not updated correctly on Register")
		t.Fail()
	}
}

func TestRegisterStopped(t *testing.T) {
	server := &CommunicatorServer{}
	src := "REGION"
	req := &RegisterRequest{Source: &src}
	clock = utils.NewFakeClock([]time.Time{time.Unix(0, 0), time.Unix(1, 0)}, false)
	testVM := newRegionalVM("", "")
	testVM.state = stopped
	vms = map[string]*regionalVM{"REGION": testVM}

	_, err := server.Register(nil, req)
	if err != nil {
		t.Log("TestRegisterStopped: Error returned on valid source input")
		t.FailNow()
	}

	if testVM.state != stopped {
		t.Log("TestRegisterStopped: vm state changed from 'stopped'")
		t.Fail()
	}
}

func TestPingExpected(t *testing.T) {
	server := &CommunicatorServer{}
	src := "REGION"
	stp := false
	req := &Heartbeat{Source: &src, Stop: &stp}
	clock = utils.NewFakeClock([]time.Time{time.Unix(0, 0), time.Unix(1, 0)}, false)
	testVM := newRegionalVM("", "")
	vms = map[string]*regionalVM{"REGION": testVM}

	res, err := server.Ping(nil, req)
	if err != nil {
		t.Log("TestPingExpected: Error returned on valid input")
		t.FailNow()
	}

	if testVM.state != probing || !testVM.lastPing.Equal(time.Unix(1, 0)) {
		t.Log("TestRegisterExpected: VM not updated correctly on Ping")
		t.Fail()
	}
	if res.GetSource() != "Controller" || res.GetStop() != false {
		t.Log("TestPingExpected: Incorrect response from ping server")
		t.Fail()
	}
}

func TestPingClientStop(t *testing.T) {
	server := &CommunicatorServer{}
	src := "REGION"
	stp := true
	req := &Heartbeat{Source: &src, Stop: &stp}
	clock = utils.NewFakeClock([]time.Time{time.Unix(0, 0), time.Unix(1, 0)}, false)
	testVM := newRegionalVM("", "")
	testVM.state = stopped
	vms = map[string]*regionalVM{"REGION": testVM}

	res, err := server.Ping(nil, req)
	if err != nil {
		t.Log("TestPingExpected: Error returned on valid input")
		t.FailNow()
	}

	if testVM.state != stopped || !testVM.lastPing.Equal(time.Unix(1, 0)) {
		t.Log("TestRegisterExpected: VM not updated correctly on Register")
		t.Fail()
	}
	if res.GetSource() != "Controller" || res.GetStop() != false {
		t.Log("TestPingExpected: Incorrect response from ping server")
		t.Fail()
	}
}

func TestCheckVMs(t *testing.T) {
	clock = utils.NewFakeClock([]time.Time{time.Unix(1, 0)}, false)
	testVM := &regionalVM{lastPing: time.Unix(0, 0)}
	vms = map[string]*regionalVM{"REGION": testVM}
	cfg, err := getTestConfig("testConfig.txt")
	if err != nil {
		t.Log("TestRegisterExpected: unable to parse test configuration file")
		t.Fail()
	}
	config = cfg
	stopping = true

	checkVMs(0 * time.Second)
}

func TestIsTimedOut(t *testing.T) {
	clock = utils.NewFakeClock([]time.Time{time.Unix(1, 0)}, true)
	vm := &regionalVM{lastPing: time.Unix(0, 0)}
	expected := []bool{false, false, true, false, true, false, true, false, false, false}
	states := []vmState{inactive, starting, idle, probing, stopped}

	for i := 0; i < len(states); i++ {
		vm.state = states[i]
		if isTimedOut(vm, 0*time.Second) != expected[2*i] {
			t.Logf("TestIsTimedOut: incorrect output on timeout with state: %v, expected %v", vm.state, expected[2*i])
			t.Fail()
		}
		if isTimedOut(vm, 1*time.Second) != expected[2*i+1] {
			t.Logf("TestIsTimedOut: incorrect output on no timeout with state: %v, expected %v", vm.state, expected[2*i+i])
			t.Fail()
		}
	}
}
