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

const numProbes = 10

func makeTestProbeConfig() *controller.ProbeConfig {
	return &controller.ProbeConfig{
		Type:         controller.ProbeType_UNSPECIFIED,
		SendInterval: 0,
	}
}

func makeTestProbeConfigs() *controller.ProbeConfigs {
	testProbe := makeTestProbeConfig()
	var testProbes []*controller.ProbeConfig
	for i := 0; i < numProbes; i++ {
		testProbes = append(testProbes, testProbe)
	}
	return &controller.ProbeConfigs{Probe: testProbes}
}

func equals(o1 *controller.ProbeConfig, o2 *controller.ProbeConfig) bool {
	return o1.GetType() == o2.GetType() && o1.GetSendInterval() == o2.GetSendInterval()
}

func TestMakeProbes(t *testing.T) {
	probeConfigs = makeTestProbeConfigs()
	compareConfig := makeTestProbeConfig()
	testProbes := makeProbes()
	for i := 0; i < numProbes; i++ {
		if !equals(testProbes[i].config, compareConfig) {
			t.Log("TestMakeProbes: probe's config does not match provided config")
			t.Fail()
		}
	}
}

func TestStartProbes(t *testing.T) {
	probing = false
	testConfig := makeTestProbeConfig()
	var testProbes []*probe
	for i := 0; i < numProbes; i++ {
		testProbes = append(testProbes, newProbe(testConfig))
	}
	pwg := startProbes(testProbes)
	pwg.Wait()
}

func TestStopProbes(t *testing.T) {
	testConfig := makeTestProbeConfig()
	clock = utils.NewFakeClock([]time.Time{}, true)
	maker = utils.NewFakeCommandMaker([]string{""}, []bool{false}, true)
	pwg := new(sync.WaitGroup)
	for i := 0; i < numProbes; i++ {
		pwg.Add(1)
		go newProbe(testConfig).probe(pwg)
	}
	stopProbes(pwg)
}

func TestStartResolver(t *testing.T) {
	clock = utils.NewFakeClock([]time.Time{time.Unix(0, 0)}, true)
	maker = utils.NewFakeCommandMaker([]string{"0.0"}, []bool{false}, false)
	rwg, err := startResolver()
	if err != nil {
		t.Logf("TestStartResolver: error returned on valid input")
		t.FailNow()
	}
	addProbe(nil)
	rwg.Wait()
}

func TestStopResolver(t *testing.T) {
	rwg := new(sync.WaitGroup)
	rwg.Add(1)
	go resolveProbes(rwg)
	stopResolver(rwg)
}
