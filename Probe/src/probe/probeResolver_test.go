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
	"sync"
	"testing"
	"time"

	"github.com/FirebaseExtended/fcm-external-prober/Controller/src/controller"
	"github.com/FirebaseExtended/fcm-external-prober/Probe/src/utils"
)

func TestResolveProbes(t *testing.T) {
	maker = utils.NewFakeCommandMaker([]string{"1000", "nf", "nf"}, []bool{false, false, false}, false)
	clock = utils.NewFakeClock([]time.Time{time.Unix(3, 0), time.Unix(100, 0)}, false)
	fakeLogger := new(fakeLogger)
	logger = fakeLogger

	timeout := int32(2)
	ptype := controller.ProbeType_UNSPECIFIED
	testConfig := &controller.ProbeConfig{ReceiveTimeout: &timeout, Type: &ptype}
	testProbes := []*sentProbe{newSentProbe(time.Unix(1, 0), &probe{config: testConfig}),
		newSentProbe(time.Unix(2, 0), &probe{config: testConfig})}
	wg := new(sync.WaitGroup)
	wg.Add(1)
	resolve = true

	initResolver()
	go resolveProbes(wg)
	addProbe(testProbes[0])
	addProbe(testProbes[1])
	closeUnresolved()

	wg.Wait()

	// There should be 3 logs: one resolved, one unresolved, one timeout
	logs := fakeLogger.testLogs
	if len(logs) != 3 {
		t.Logf("TestResolveProbes: not all probes resolved: %d", len(logs))
		t.FailNow()
	}
	if logs[0].time != testProbes[0].sendTime.Format(timeLogFormat) || logs[0].state != "resolved" || logs[0].latency != 0 {
		t.Logf("TestResolveProbe: probe 1 resolved incorrectly")
		t.Fail()
	}
	if logs[1].time != testProbes[1].sendTime.Format(timeLogFormat) || logs[1].state != "unresolved" || logs[1].latency != -1 {
		t.Logf("TestResolveProbe: probe 2 unresolved incorrectly")
		t.Fail()
	}
	if logs[2].time != testProbes[1].sendTime.Format(timeLogFormat) || logs[2].state != "timeout" || logs[2].latency != -1 {
		t.Logf("TestResolveProbe: probe 3 not timed out")
		t.Fail()
	}
}

func TestResolveProbe(t *testing.T) {
	maker = utils.NewFakeCommandMaker([]string{"1500"}, []bool{false}, false)
	ptype := controller.ProbeType_UNSPECIFIED
	testConfig := &controller.ProbeConfig{Type: &ptype}
	testSentProbe := newSentProbe(time.Unix(1, 0), &probe{config: testConfig})
	fakeLogger := new(fakeLogger)
	logger = fakeLogger

	res := resolveProbe(testSentProbe)

	if !res {
		t.Log("TestResolveProbe: probe not resolved on valid input")
		t.Fail()
	}
	log := fakeLogger.testLogs[0]
	if log.state != "resolved" {
		t.Logf("TestResolveProbe: incorrect probe state: actual: %s expected: resolved", log.state)
		t.Fail()
	}
	if log.latency != 500 {
		t.Logf("TestResolveProbe: incorrect latency: actual %d expected: 500", log.latency)
		t.Fail()
	}
}

func TestResolveProbeGetError(t *testing.T) {
	maker = utils.NewFakeCommandMaker([]string{"INVALID_COMMAND"}, []bool{true}, false)
	ptype := controller.ProbeType_UNSPECIFIED
	testConfig := &controller.ProbeConfig{Type: &ptype}
	testSentProbe := newSentProbe(time.Unix(1, 0), &probe{config: testConfig})
	fakeLogger := new(fakeLogger)
	logger = fakeLogger

	res := resolveProbe(testSentProbe)

	if !res {
		t.Log("TestResolveProbeGetError: probe not resolved on getMessage error")
		t.Fail()
	}
	log := fakeLogger.testLogs[0]
	if log.state != "error" {
		t.Logf("TestResolveProbeGetError: incorrect probe state: actual: %s expected: error", log.state)
		t.Fail()
	}
	if log.latency != -1 {
		t.Logf("TestResolveProbeGetError: incorrect latency: actual %d expected: -1", log.latency)
		t.Fail()
	}
}

func TestResolveProbeTimeout(t *testing.T) {
	maker = utils.NewFakeCommandMaker([]string{"nf"}, []bool{false}, false)
	timeout := int32(2)
	ptype := controller.ProbeType_UNSPECIFIED
	testConfig := &controller.ProbeConfig{ReceiveTimeout: &timeout, Type: &ptype}
	testSentProbe := newSentProbe(time.Unix(1, 0), &probe{config: testConfig})
	// Set time to after timeout time
	clock = utils.NewFakeClock([]time.Time{time.Unix(2, 0).Add(time.Duration(timeout) * time.Second)}, false)
	fakeLogger := new(fakeLogger)
	logger = fakeLogger

	res := resolveProbe(testSentProbe)

	if !res {
		t.Log("TestResolveProbeTimeout: probe not resolved on timeout")
		t.Fail()
	}
	log := fakeLogger.testLogs[0]
	if log.state != "timeout" {
		t.Logf("TestResolveProbeTimeout: incorrect probe state: actual: %s expected: timeout", log.state)
		t.Fail()
	}
	if log.latency != -1 {
		t.Logf("TestResolveProbeTimeout: incorrect latency: actual %d expected: -1", log.latency)
		t.Fail()
	}
}

func TestResolveProbeUnresolved(t *testing.T) {
	maker = utils.NewFakeCommandMaker([]string{"nf"}, []bool{false}, false)
	timeout := int32(2)
	ptype := controller.ProbeType_UNSPECIFIED
	testConfig := &controller.ProbeConfig{ReceiveTimeout: &timeout, Type: &ptype}
	testSentProbe := newSentProbe(time.Unix(1, 0), &probe{config: testConfig})
	// Set time to before timeout time
	clock = utils.NewFakeClock([]time.Time{time.Unix(1, 0).Add(time.Duration(timeout) * time.Second)}, false)
	fakeLogger := new(fakeLogger)
	logger = fakeLogger

	res := resolveProbe(testSentProbe)

	if res {
		t.Log("TestResolveProbeUnresolved: probe not unresolved before timeout")
		t.Fail()
	}
	log := fakeLogger.testLogs[0]
	if log.state != "unresolved" {
		t.Logf("TestResolveProbeUnresolved: incorrect probe state: actual: %s expected: unresolved", log.state)
		t.Fail()
	}
	if log.latency != -1 {
		t.Logf("TestResolveProbeUnresolved: incorrect latency: actual %d expected: -1", log.latency)
		t.Fail()
	}
}

func TestResolveProbeInvalidMessage(t *testing.T) {
	maker = utils.NewFakeCommandMaker([]string{"INVALID_MESSAGE"}, []bool{false}, false)
	ptype := controller.ProbeType_UNSPECIFIED
	testConfig := &controller.ProbeConfig{Type: &ptype}
	testSentProbe := newSentProbe(time.Unix(1, 0), &probe{config: testConfig})
	clock = utils.NewFakeClock([]time.Time{time.Unix(1, 0)}, false)
	fakeLogger := new(fakeLogger)
	logger = fakeLogger

	res := resolveProbe(testSentProbe)

	if !res {
		t.Log("TestResolveProbeInvalidMessage: probe not resolved with error recorded reception time")
		t.Fail()
	}
	log := fakeLogger.testLogs[0]
	if log.state != "error" {
		t.Logf("TestResolveProbeInvalidMessage: incorrect probe state: actual: %s expected: error", log.state)
		t.Fail()
	}
	if log.latency != -1 {
		t.Logf("TestResolveProbeInvalidMessage: incorrect latency: actual %d expected: -1", log.latency)
		t.Fail()
	}
}

func TestCalculateLatency(t *testing.T) {
	res, err := calculateLatency(time.Unix(10, 0), "10001")

	if err != nil {
		t.Logf("TestCalculateLatency: error on valid input: %s", err.Error())
		t.Fail()
	}
	if res != 1 {
		t.Logf("TestCalculateLatency: incorrect result: actual: %d expected: %d", res, 1)
		t.Fail()
	}
}

func TestCalculateLatencyError(t *testing.T) {
	res, err := calculateLatency(time.Unix(10, 0), "INVALID_TIME")

	if err == nil {
		t.Logf("TestCalculateLatency: no error on invalid input")
		t.Fail()
	}
	if res != -1 {
		t.Logf("TestCalculateLatency: incorrect result: actual: %d expected: %d", res, -1)
		t.Fail()
	}
}
