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

func TestRestartVM(t *testing.T) {
	maker = utils.NewFakeCommandMaker(make([]string, 3), make([]bool, 3), false)
	clock = utils.NewFakeClock([]time.Time{time.Unix(0, 0), time.Unix(0,0)}, false)
	stopping = false
	vm := newRegionalVM("", "")

	if vm.state != inactive {
		t.Log("TestRestartVM: VM initialized with incorrect state")
		t.Fail()
	}

	vm.restartVM()

	if vm.state != starting {
		t.Log("TestRestartVM: VM state not set to 'starting' after restarting when controller is not stopping")
		t.Fail()
	}

	stopping = true
	vm.restartVM()

	if vm.state != stopped {
		t.Log("TestRestartVM: VM state not set to 'stopping' after restarting when controller is stopping")
		t.Fail()
	}
}

func TestSetState(t *testing.T) {
	maker = utils.NewFakeCommandMaker(make([]string, 5), make([]bool, 5), false)
	clock = utils.NewFakeClock([]time.Time{time.Unix(0, 0)}, false)
	stopping = false
	stoppedVMs = 0
	vm := newRegionalVM("", "")
	testStates := []vmState{inactive, starting, idle, probing, stopped}

	for i := 0; i < 4; i++ {
		vm.setState(testStates[i])

		if stoppedVMs > 0 {
			t.Log("TestSetState: stoppedVMs incremented when vm state was not set to stopping")
			t.Fail()
		}
	}

	vm.setState(testStates[4])

	if stoppedVMs != 1 {
		t.Log("TestSetState: stoppedVMs not incremented when vm state was set to stopping")
		t.Fail()
	}

	// Once a VM is in stopping state, it should not be able to leave that state
	vm.setState(inactive)

	if vm.state != stopped {
		t.Log("TestSetState: vm escaped stopped state")
		t.Fail()
	}
}
