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

package probe

import (
	"log"
	"sync"

	"github.com/FirebaseExtended/fcm-external-prober/Controller/src/controller"
	"github.com/FirebaseExtended/fcm-external-prober/Probe/src/utils"
)

var (
	probeConfigs *controller.ProbeConfigs
	account      *controller.AccountInfo
	maker        utils.CommandMaker
	clock        utils.Timer
	logger       Logger
	fcmAuth      Auth
	deviceToken  string
	probing      = true
	probeLock    sync.Mutex
)

// Handles startup/teardown of emulator/app, also starts and stops probing
func Control(mk utils.CommandMaker, clk utils.Timer, lg Logger) {
	maker = mk
	clock = clk
	logger = lg

	err := initClient()
	if err != nil {
		log.Fatalf("Control: unable to initialize gRPC client")
	}

	defer destroyEnvironment()

	initEnvironment()
	tok, err := getToken()
	if err != nil {
		log.Fatalf("Control: could not acquire device token, %s", err.Error())
	}
	deviceToken = tok
	ps := makeProbes()
	pwg := startProbes(ps)
	rwg := startResolver()

	communicate()

	stopProbes(pwg)
	stopResolver(rwg)

	err = confirmStop()

	// Connection was lost to the controller, so assume it has terminated and delete VM instance
	if err != nil {
		deleteVM()
	}
}

func initEnvironment() {
	err := startEmulator()
	if err != nil {
		log.Fatalf("probe: could not start emulator: %s", err.Error())
	}
	err = startApp()
	if err != nil {
		log.Fatalf("probe: could not install app: %s", err.Error())
	}
}

func destroyEnvironment() {
	err := uninstallApp()
	if err != nil {
		log.Printf("probe: unable to uninstall app: %s", err.Error())
	}
	err = killEmulator()
	if err != nil {
		log.Fatalf("probe: could not kill emulator: %s", err.Error())
	}
}

func makeProbes() []*probe {
	var ret []*probe
	for _, p := range probeConfigs.GetProbe() {
		ret = append(ret, newProbe(p))
	}
	return ret
}

func startProbes(ps []*probe) *sync.WaitGroup {
	pwg := new(sync.WaitGroup)
	for _, p := range ps {
		pwg.Add(1)
		go p.probe(pwg)
	}
	return pwg
}

func stopProbes(pwg *sync.WaitGroup) {
	probeLock.Lock()
	probing = false
	probeLock.Unlock()
	pwg.Wait()
}

func startResolver() *sync.WaitGroup {
	rwg := new(sync.WaitGroup)
	rwg.Add(1)
	initResolver()
	go resolveProbes(rwg)
	return rwg
}

func stopResolver(rwg *sync.WaitGroup) {
	closeUnresolved()
	rwg.Wait()
}

func deleteVM() {
	log.Printf("delete VM")
	//maker.Command("gcloud", "compute", "instances", "delete", hostname, "--zone", hostname, "--quiet")
}
