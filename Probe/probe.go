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
	"flag"
	"log"
	"sync"
	"time"
)

const timeFileFormat = "2006-01-02-T150405.000-Z0700"
const timeLogFormat = time.UnixDate

var region string
var probeType string
var probeNumber int
var serviceAccount string
var deviceToken string
var projectID string
var probeInterval int
var probeTimeout int

var clock Timer
var maker CommandMaker
var pLog logger

func main() {
	initFlags()
	clock = new(probeClock)
	maker = new(cmdMaker)
	pLog = newCloudLogger()
	probe()
}

func initFlags() {
	flag.StringVar(&region, "region", "default", "regional server in which VM is located")
	flag.StringVar(&probeType, "type", "default", "type of probe behavior")
	flag.IntVar(&probeNumber, "number", 4, "number of total probes (messages) sent")
	flag.StringVar(&serviceAccount, "account",
		"send-service-account@gifted-cooler-279818.iam.gserviceaccount.com", "service account with FCM privileges")
	flag.StringVar(&projectID, "project", "gifted-cooler-279818", "GCP project in which this VM exists")
	flag.IntVar(&probeInterval, "interval", 4, "number of seconds between successive probes")
	flag.IntVar(&probeTimeout, "timeout", 10, "number of seconds before a probe will be timed out if unreconciled")
	flag.Parse()
}

func probe() {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go resolveProbes(wg)

	fcmAuth := new(Auth)

	err := startEmulator()
	if err != nil {
		log.Fatalf("probe: could not start emulator: %s", err.Error())
	}

	if probeType == "default" {
		err := startApp()
		if err != nil {
			log.Fatalf("probe: could not install app: %s", err.Error())
		}

		deviceToken, err = getToken()
		if err != nil {
			log.Fatalf("probe: unable to retrieve FCM token: %s", err.Error())
		}

		for i := 0; i < probeNumber; i++ {
			tim := clock.Now()
			err = fcmAuth.sendMessage(tim.Format(timeFileFormat))
			if err != nil {
				log.Printf("probe: unable to send message: %s", err.Error())
			}
			addProbe(tim)
			// Time interval between probes
			time.Sleep(time.Duration(probeInterval) * time.Second)
		}
	}

	stopResolving()
	wg.Wait()
	err = uninstallApp()
	if err != nil {
		log.Printf("probe: unable to uninstall app: %s", err.Error())
	}
	err = killEmulator()
	if err != nil {
		log.Fatalf("probe: could not kill emulator: %s", err.Error())
	}
}
