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

	"github.com/FirebaseExtended/fcm-external-prober/Probe/src/probe"
	"github.com/FirebaseExtended/fcm-external-prober/Probe/src/utils"
)

func initFlags() (*string, *string) {
	cfgs := flag.String("probes", "", "ProbeConfigs protobuf with probe behaviors")
	acct := flag.String("account", "", "AccountInfo protobuf with GCP account info")
	flag.Parse()

	return cfgs, acct
}

func main() {
	cfgs, acct := initFlags()
	m := new(utils.CmdMaker)
	c := new(utils.ProbeClock)
	l := probe.NewCloudLogger()
	probe.Control(cfgs, acct, m, c, l)
}
