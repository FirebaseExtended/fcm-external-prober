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

func TestProbe(t *testing.T) {
	interval := int32(0)
	typ := controller.ProbeType_UNSPECIFIED
	cfg := &controller.ProbeConfig{SendInterval: &interval, Type: &typ}
	testClock := utils.NewFakeBoolClock(make([]time.Time, 6), &probing)
	clock = testClock
	testMaker := utils.NewFakeCommandMaker([]string{"0.0", ""}, []bool{false}, true)
	maker = testMaker

	p := newProbe(cfg)
	pwg := new(sync.WaitGroup)
	pwg.Add(1)

	initResolver()
	go p.probe(pwg)
	pwg.Wait()
	addProbe(nil)

	// Subtract command call and one clock call for initResolver
	if testClock.TimesCalled() - 2 != 2*(testMaker.TimesCalled() - 1) {
		t.Log("TestProbe: clock and maker not accessed same number of times")
		t.Fail()
	}

	i := 0
	current := <-unresolved
	for current != nil {
		i++
		current = <-unresolved
	}

	// Subtract command call for initResolver
	if i != testMaker.TimesCalled() - 1 {
		t.Log("TestProbe: number of probes sent not equal to number of times sendMessage was called")
		t.Fail()
	}
}
