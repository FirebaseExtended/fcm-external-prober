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

package main

import (
	"flag"
	"github.com/golang/protobuf/proto"
	"log"
	"probe"
	"utils"
)

var probeConfigs *probe.ProbeConfigs
var account *probe.AccountInfo
var maker utils.CommandMaker
var clock utils.Timer
var logger utils.Logger

func initFlags() {
	p := *flag.String("probes", "", "ProbeConfigs protobuf with probe behaviors")
	a := *flag.String("account", "", "AccountInfo protobuf with GCP account info")
	flag.Parse()
	err := proto.UnmarshalText(p, probeConfigs)
	if err != nil {
		log.Fatalf("initFlags: invalid probe configuration: %s", err.Error())
	}
	err = proto.UnmarshalText(a, account)
	if err != nil {
		log.Fatalf("initFlags: invalid account information: %s", err.Error())
	}
}

func main() {
	initFlags()
	maker = new(utils.CmdMaker)
	logger = utils.NewCloudLogger()
	handler := probe.NewAppHandler()
	handler.StartEmulator()

	p := makeProbes()
	startProbes(p)
	r := probe.NewResolver(logger)
	go r.ResolveProbes(p)
}

func makeProbes() []*probe.Probe {
	var ret []*probe.Probe
	for _, p := range probeConfigs.Probes {
		ret = append(ret, probe.NewProbe(p, clock))
	}
	return ret
}

func startProbes(ps []*probe.Probe) {
	for _, p := range ps {
		p.Probe()
	}
}

