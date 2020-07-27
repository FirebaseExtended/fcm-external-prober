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

	"github.com/FirebaseExtended/fcm-external-prober/Probe/src/utils"
	"github.com/golang/protobuf/proto"
)

var maker utils.CommandMaker

type Controller struct {
	config *ControllerConfig
	vms    map[string]*regionalVM
}

// Create a new controller with a provided configuration
func NewController(cfg string, mk utils.CommandMaker) *Controller {
	maker = mk
	ret := &Controller{new(ControllerConfig), make(map[string]*regionalVM)}

	err := proto.UnmarshalText(cfg, ret.config)
	if err != nil {
		log.Fatalf("Controller: invalid configuration: %s", err.Error())
	}
	return ret
}

func (ctrl *Controller) getPossibleZones() {
	zones, err := getCompatZones([]string{ctrl.config.GetMinCpu()})
	if err != nil {
		//TODO(langenbahn): Log this when logging is implemented
		log.Fatalf("Controller: unable to generate list of VM zones")
	}
	for _, n := range zones {
		// For now, the probe name is the same as the zone name. This will change if multiple VMs are required in a zone
		ctrl.vms[n] = newRegionalVM(n, n, ctrl.config.GetMinCpu(), ctrl.config.GetDiskImageName())
	}
}

// Start all VMs in regions in which the required hardware is available, and for which there are probes specified
func (ctrl *Controller) StartVMs() {
	ctrl.getPossibleZones()
	for _, p := range ctrl.config.Probes.Probe {
		vm, ok := ctrl.vms[p.GetRegion()+"-a"]
		if !ok {
			log.Printf("Controller: zone %s in region %s does not meet minimum requirements or does not exist", p.GetRegion()+"-a", p.GetRegion())
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
		if vm.active {
			vm.startProbes(ctrl.config.GetAccount())
		}
	}
}
