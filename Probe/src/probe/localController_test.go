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
	"io/ioutil"
	"sync"
	"testing"
	"time"

	"github.com/FirebaseExtended/fcm-external-prober/Probe/src/utils"
	"github.com/golang/protobuf/proto"
)

const numProbes = 10

func makeTestProbeConfig() *ProbeConfig {
	t := ProbeType_UNSPECIFIED
	si := int32(0)
	return &ProbeConfig{
		Type:         &t,
		SendInterval: &si,
	}
}

func makeTestProbeConfigs() *ProbeConfigs {
	testProbe := makeTestProbeConfig()
	var testProbes []*ProbeConfig
	for i := 0; i < numProbes; i++ {
		testProbes = append(testProbes, testProbe)
	}
	return &ProbeConfigs{Probes: testProbes}
}

func equals(o1 *ProbeConfig, o2 *ProbeConfig) bool {
	return o1.GetType() == o2.GetType() && o1.GetSendInterval() == o2.GetSendInterval()
}

func TestTest(t *testing.T) {
	probeConfig := makeTestProbeConfigs()

	ioutil.WriteFile("testProto.txt", []byte(proto.MarshalTextString(probeConfig)), 0644)
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
	rwg := startResolver()
	addProbe(nil)
	rwg.Wait()
}

func TestStopResolver(t *testing.T) {
	rwg := new(sync.WaitGroup)
	rwg.Add(1)
	go resolveProbes(rwg)
	stopResolver(rwg)
}
