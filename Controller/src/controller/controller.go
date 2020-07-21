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

/* Package controller is responsible for creating VMs in various regions, according to
 * the provided configuration string
 * The controller is also responsible for initiating probing on the VMs that are created
 */
package controller

import (
	"log"

	"github.com/golang/protobuf/proto"
)

type Controller struct {
	config *ControllerConfig
	vms    map[string]*regionalVM
}

// Create a new controller with a provided configuration
func NewController(cfg string) *Controller {
	ret := new(Controller)
	ret.config = new(ControllerConfig)
	err := proto.UnmarshalText(cfg, ret.config)
	if err != nil {
		log.Fatalf("Controller: invalid configuration: %s", err.Error())
	}
	return ret
}

func (ctrl *Controller) getPossibleZones() {
	z, err := getCompatZones([]string{ctrl.config.GetMinCpu()})
	if err != nil {
		//TODO(langenbahn): Log this when logging is implemented
		log.Fatalf("Controller: unable to generate list of VM zones")
	}
	for _, c := range z {
		ctrl.vms[c] = newRegionalVM(c, c, ctrl.config.GetMinCpu(), ctrl.config.GetDiskImageName())
	}
}

// Start all VMs in regions in which the required hardware is available, and for which there are probes specified
func (ctrl *Controller) StartVMs() {
	ctrl.getPossibleZones()
	for _, p := range ctrl.config.Probes.Probe {
		vm, ok := ctrl.vms[p.GetRegion()+"a"]
		if !ok {
			log.Printf("Controller: zone %s in region %s does not meet minimum requirements or does not exist", p.GetRegion()+"a", p.GetRegion())
			continue
		}
		if !vm.active {
			err := vm.startVM(ctrl.config.GetMinCpu(), ctrl.config.GetDiskImageName(), ctrl.config.GetAccount().GetServiceAccount())
			if err != nil {
				log.Printf("Controller: regional VM could not be started in zone %s", vm.zone)
			}
		}
		vm.probes = append(vm.probes, p)
	}
}

// Start probing logic in all active VMs
func (ctrl *Controller) StartProbes() {
	for _, vm := range ctrl.vms {
		vm.startProbes(ctrl.config.GetAccount())
	}
}
