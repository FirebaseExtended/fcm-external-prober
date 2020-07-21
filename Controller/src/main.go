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
	"controller"
	"flag"
	"io/ioutil"
	"log"
)

func main() {
	cf := flag.String("config", "config.txt", "text file in which a protobuf with config information is located")
	c, err := ioutil.ReadFile(*cf)
	if err != nil {
		log.Fatalf("Main: could not read from specified config file: %s", err.Error())
	}
	ctrl := controller.NewController(string(c))
	ctrl.StartVMs()
	ctrl.StartProbes()
}
