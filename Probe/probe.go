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

/*
Package probe implements an FCM probe that will:
	Initialize a new emulator and target app
	Send a specified number of messages to the app
	Attempt to verify that the app received those messages
	Log findings to Cloud Logger
	Terminate emulator and target app when finished
*/
package main

import (
	"log"
	"os"
	"strconv"
	"time"
)

const timeFileFormat = "2006-01-02-T150405.000-Z0700"
const timeLogFormat = time.UnixDate

var region string
var probeType string
var probeNumber int
var deviceToken string

func main() {
	arg := os.Args[1:]
	region = arg[0]
	probeType = arg[1]
	var err error
	probeNumber, err = strconv.Atoi(arg[2])
	if err != nil {
		log.Fatal("probe: invalid argument (non-integer value for number of probes)")
	}
	unresolvedProbes = make([]time.Time, 0)
	probe()
}

func probe() {
	c := make(chan bool)
	go resolveProbes(c)
	err := startEmulator()
	if err != nil {
		log.Fatal("probe: could not start emulator")
	}
	if probeType == "default" {
		err = startApp()
		if err != nil {
			log.Print("probe: could not install app: " + err.Error())
			//continue
		}
		time.Sleep(10 * time.Second)
		for i := 0; i < probeNumber; i++ {
			deviceToken, err = getToken()
			if err != nil {
				log.Print("probe: unable to retrieve FCM token")
			} else {
				tim := time.Now()
				err = sendMessage(tim.Format(timeFileFormat))
				if err != nil {
					log.Print("probe: unable to send message: " + err.Error())
				}
				addProbe(tim)
			}
			time.Sleep(10 * time.Second)
		}
		err = uninstallApp()
		if err != nil {
			log.Print("probe: unable to uninstall app")
		}
	}
	stopResolving()
	waitForResolution(c)
	err = killEmulator()
	if err != nil {
		log.Fatal("probe: could not kill emulator")
	}
}

func waitForResolution(c chan bool) {
	_ = <-c
}
