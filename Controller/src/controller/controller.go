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
	"sync"
	"time"

	"github.com/FirebaseExtended/fcm-external-prober/Probe/src/utils"
)

var (
	maker          utils.CommandMaker
	clock          utils.Timer
	logger		   Logger
	vms            map[string]*regionalVM
	stoppedVMs     int
	stoppedVMsLock sync.Mutex
	stopping       bool
	config         *ControllerConfig
)

type Controller struct{}

// Create a new controller with a provided configuration
func NewController(cfg *ControllerConfig, cmd utils.CommandMaker, clk utils.Timer, log Logger) *Controller {
	config = cfg
	maker = cmd
	clock = clk
	logger = log
	vms = make(map[string]*regionalVM)
	return &Controller{}
}

// Start gRPC server for regional VMs to connect to
func (ctrl *Controller) InitServer() {
	err := initServer()
	if err != nil {
		log.Fatalf("Controller: unable to start rpc server, %v", err)
	}
	err = addMetadata()
	if err != nil {
		log.Fatalf("Controller: unable to add project metadata %v", err)
	}
}

// Start all VMs in regions in which the required hardware is available, and for which there are probes specified
func (ctrl *Controller) InitProbes() {
	getPossibleZones()

	// Assign all probe configurations to their designated regions
	for _, p := range config.Probes.Probe {
		vm, ok := vms[p.GetRegion()+"-a"]
		if !ok {
			logger.LogErrorf("Controller: zone %s in region %s does not meet minimum requirements or does not exist", p.GetRegion()+"-a", p.GetRegion())
			continue
		}
		vm.probes = append(vm.probes, p)
	}

	// Start VMs in regions to which probes are designated, remove those that contain no probes
	for z, v := range vms {
		if len(v.probes) == 0 {
			delete(vms, z)
			continue
		}
		err := v.startVM()
		if err != nil {
			logger.LogErrorf("Controller: regional VM could not be started in zone %s", v.zone)
		}
	}
}

func getPossibleZones() {
	z, err := getCompatZones([]string{config.GetMinCpu()})
	if err != nil {
		log.Fatalf("Controller: unable to generate list of VM zones: %v", err)
	}
	for _, n := range z {
		// For now, the probe name is the same as the zone name. This will change if multiple VMs are required in a zone
		vms[n] = newRegionalVM(n, n)
		vms[n].setState(inactive)
	}
}

func (ctrl *Controller) MonitorProbes() {
	checkVMs(time.Duration(config.PingConfig.GetTimeout()) * time.Minute)
}
