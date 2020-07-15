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
)

var probeConfigs *probe.ProbeConfigs
var account *probe.AccountInfo
var maker CommandMaker
var pLog logger

func initFlags() {


	err := proto.UnmarshalText(*flag.String("probes", "", "ProbeConfigs protobuf with probe behaviors"), probeConfigs)
	if err != nil {
		log.Fatalf("initFlags: invalid probe configuration proto: %s", err.Error())
	}
	err = proto.UnmarshalText(*flag.String("account", "", "AccountInfo protobuf with GCP account info"), account)
	if err != nil {
		log.Fatalf("initFlags: invalid account information proto: %s", err.Error())
	}
	flag.Parse()
}

func main() {

}
