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
	"sync"

	"github.com/FirebaseExtended/fcm-external-prober/Controller/src/controller"
	"github.com/FirebaseExtended/fcm-external-prober/Probe/src/utils"
)

var (
	probeConfigs *controller.ProbeConfigs
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

	acquireData()

	err := initClient()
	if err != nil {
		logger.LogFatalf("Control: unable to initialize gRPC client: %v", err)
	}

	defer destroyEnvironment()

	initEnvironment()
	tok, err := getToken()
	if err != nil {
		logger.LogFatalf("Control: could not acquire device token: %v", err)
	}
	deviceToken = tok
	ps := makeProbes()
	rwg, err := startResolver()
	if err != nil {
		logger.LogFatalf("Control: unable to start resolver: %v", err)
	}
	pwg := startProbes(ps)


	err = communicate()
	if err != nil {
		logger.LogErrorf("Control: communication error, %v", err)
	}

	stopProbes(pwg)
	stopResolver(rwg)

	err = confirmStop()

	// Connection was lost to the controller, so assume it has terminated and delete VM instance
	if err != nil {
		deleteVM()
	}
}

func acquireData() {
	// Set hostname first so that errors are identifiable by VM if they occur
	hostname, err := getHostname()
	if err != nil {
		logger.LogFatalf("acquireData: unable to resolve hostname: %v", err)
	}
	// Update region in logger now that it has been acquired
	logger.SetRegion(hostname)

	// Get Metadata that will be used to connect to the controller
	err = getMetadata()
	if err != nil {
		logger.LogFatalf("acquireData: unable to acquire metadata: %v", err)
	}

	// Update logger send destinations that were acquired with metadata
	logger.SetError(metadata.GetErrorLogDestination())
	logger.SetLog(metadata.GetProbeLogDestination())
}

func initEnvironment() {
	err := startEmulator()
	if err != nil {
		logger.LogFatalf("initEnvironment: could not start emulator: %v", err)
	}
	err = startApp()
	if err != nil {
		logger.LogFatalf("initEnvironment: could not install app: %v", err)
	}
}

func destroyEnvironment() {
	err := uninstallApp()
	if err != nil {
		logger.LogErrorf("destroyEnvironment: unable to uninstall app: %v", err)
	}
	err = killEmulator()
	if err != nil {
		logger.LogFatalf("destroyEnvironment: could not kill emulator: %v", err)
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

func startResolver() (*sync.WaitGroup, error) {
	rwg := new(sync.WaitGroup)
	rwg.Add(1)
	err := initResolver()
	if err != nil {
		return nil, err
	}
	go resolveProbes(rwg)
	return rwg, nil
}

func stopResolver(rwg *sync.WaitGroup) {
	closeUnresolved()
	rwg.Wait()
}

func deleteVM() {
	//TODO(langenbahn) Remove this when not developing on a GCE VM
	//maker.Command("gcloud", "compute", "instances", "delete", hostname, "--zone", hostname, "--quiet")
}
