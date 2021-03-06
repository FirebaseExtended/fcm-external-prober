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
	"sync"
	"time"
)

type vmState int

const (
	inactive vmState = iota
	starting
	idle
	probing
	stopped
)

type regionalVM struct {
	name      string
	zone      string
	state     vmState
	stateLock sync.Mutex
	probes    []*ProbeConfig
	lastPing  time.Time
}

func newRegionalVM(name string, zone string) *regionalVM {
	return &regionalVM{name: name, zone: zone, lastPing: clock.Now()}
}

func (vm *regionalVM) startVM() error {
	err := maker.Command("gcloud", "compute", "instances", "create", vm.name, "--zone", vm.zone,
		"--quiet", "--min-cpu-platform", config.GetMinCpu(),
		"--service-account", config.GetMetadata().GetAccount().GetServiceAccount(),
		"--image", config.GetImageName(), "--machine-type", "n1-standard-4", "--scopes", "cloud-platform",
		"--metadata-from-file=startup-script="+config.GetStartupScriptPath()).Run()
	if err != nil {
		return err
	}
	vm.updatePingTime()
	vm.setState(starting)
	return nil
}

func (vm *regionalVM) stopVM() {
	err := maker.Command("gcloud", "compute", "instances", "delete", vm.name, "--zone", vm.zone, "--quiet").Run()
	if err != nil {
		logger.LogErrorf("stopVM: unable to stop VM %s in zone %s: %v", vm.name, vm.zone, err)
	}
}

func (vm *regionalVM) restartVM() {
	vm.stopVM()
	if stopping || vm.state == stopped {
		// Controller is shutting down, so do not start VM again
		vm.setState(stopped)
		return
	}
	vm.setState(starting)
	err := vm.startVM()
	if err != nil {
		logger.LogErrorf("stopVM: unable to restart VM %s in zone %s: %v", vm.name, vm.zone, err)
		vm.setState(stopped)
	}
}

func (vm *regionalVM) updatePingTime() {
	vm.lastPing = clock.Now()
}

func (vm *regionalVM) setState(s vmState) {
	vm.stateLock.Lock()
	defer vm.stateLock.Unlock()
	if vm.state == stopped {
		return
	} else if s == stopped {
		stoppedVMsLock.Lock()
		stoppedVMs++
		stoppedVMsLock.Unlock()
	}
	vm.state = s
}
