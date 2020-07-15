package controller

import (
	"github.com/golang/protobuf/proto"
	"log"
	"os/exec"
)

type regionalVM struct {
	name string
	zone string
	cpuMin string
	imageName string
	active bool
	probes []*ProbeConfig
}

func newRegionalVM(name string, zone string, cpu string, img string) *regionalVM {
	ret := new(regionalVM)
	ret.name = name
	ret.zone = zone
	ret.cpuMin = cpu
	ret.imageName = img
	return ret
}

func (vm regionalVM) startVM(cpu string, img string, sa string) error {
	err := exec.Command("gcloud", "compute", "instances", "create", vm.name, "--zone", vm.zone,
		"--min-cpu-platform", cpu, "--image", img, "--service-account", sa).Run()
	if err != nil {
		return err
	}
	vm.active = true
	return nil
}

func (vm regionalVM) startProbes(acc *AccountInfo) {
	s := new(ProbeConfigs)
	s.Probe = vm.probes
	p := proto.MarshalTextString(s)
	a := proto.MarshalTextString(acc)
	// Send protobuf string with probe behavior as argument to probe function on VM
	err := exec.Command("gcloud","compute", "ssh", vm.name, ";", "go", "probe", "-probes=" + p, "-account=" + a).Start()
	if err != nil {
		log.Printf("startProbes: unable to start probes on VM %s in zone %s", vm.name, vm.zone)
	}
}





