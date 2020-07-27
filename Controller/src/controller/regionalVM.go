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
	"log"
)

type regionalVM struct {
	name      string
	zone      string
	cpuMin    string
	imageName string
	active    bool
	probes    []*ProbeConfig
}

func newRegionalVM(name string, zone string, cpu string, img string) *regionalVM {
	return &regionalVM{name: name, zone: zone, cpuMin: cpu, imageName: img}
}

func (vm *regionalVM) startVM(cpu string, img string, sa string) error {
	err := maker.Command("gcloud", "compute", "instances", "create", vm.name, "--zone", vm.zone,
		"--min-cpu-platform", cpu, "--image", img, "--service-account", sa).Run()
	if err != nil {
		return err
	}
	vm.active = true
	return nil
}

func (vm *regionalVM) startProbes(acc *AccountInfo) {
	s := &ProbeConfigs{Probe: vm.probes}
	p := proto.MarshalTextString(s)
	a := proto.MarshalTextString(acc)

	// Send protobuf string with probe behavior as argument to probe function on VM
	err := maker.Command("gcloud", "compute", "ssh", vm.name, "--zone", vm.zone, "--force-key-file-overwrite",
		"--command", "go", "run", "Probe/main.go", "-probes="+p, "-account="+a).Run()
	if err != nil {
		log.Printf("startProbes: unable to start probes on VM %s in zone %s", vm.name, vm.zone)
	}
	log.Printf(string(str))
}

func (vm *regionalVM) stopVM() error {
	err := maker.Command("gcloud", "compute", "instances", "delete", vm.name).Run()
	if err != nil {
		log.Printf("stopVM: unable to stop VM %s in zone %s", vm.name, vm.zone)
		return err
	}
	return nil
}
