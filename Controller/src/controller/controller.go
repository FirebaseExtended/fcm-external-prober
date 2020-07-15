package controller

import (
	"github.com/golang/protobuf/proto"
	"log"
)

type Controller struct {
	config *ControllerConfig
	vms map[string]*regionalVM
}

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
		log.Fatalf("Controller: unable to generate list of VM zones")
	}
	for _, c := range z {
		ctrl.vms[c] = newRegionalVM(c, c, ctrl.config.GetMinCpu(), ctrl.config.GetDiskImageName())
	}
}

func (ctrl *Controller) StartVMs() {
	ctrl.getPossibleZones()
	for _, p := range ctrl.config.Probes.Probe {
		vm, ok := ctrl.vms[p.GetRegion() + "a"]
		if !ok {
			log.Printf("Controller: zone %s in region %s does not meet minimum requirements or does not exist", p.GetRegion() + "a", p.GetRegion())
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

func (ctrl *Controller) StartProbes() {
	for _, vm := range ctrl.vms {
		vm.startProbes(ctrl.config.GetAccount())
	}
}

