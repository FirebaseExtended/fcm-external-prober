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

package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/FirebaseExtended/fcm-external-prober/Controller/src/controller"
	"github.com/FirebaseExtended/fcm-external-prober/Probe/src/utils"
	"github.com/golang/protobuf/proto"
)

func main() {
	cf := flag.String("config", "config.txt", "text file in which a protobuf with config information is located")
	flag.Parse()
	c, err := ioutil.ReadFile(*cf)
	if err != nil {
		log.Fatalf("Main: could not read from specified config file: %s", err.Error())
	}
	cfg := new(controller.ControllerConfig)
	err = proto.UnmarshalText(string(c), cfg)
	if err != nil {
		log.Fatalf("Main: invalid configuration: %s", err.Error())
	}
	ctrl := controller.NewController(cfg, new(utils.CmdMaker), new(utils.ProbeClock), new(controller.ControllerLogger))
	ctrl.InitServer()
	ctrl.InitProbes()
	ctrl.MonitorProbes()
}
