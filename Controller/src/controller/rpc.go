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

package controller

var cert []byte

func initServer() {
	err := makeCert()
	if err != nil {
		return err
	}
	tls, err := credentials.NewServerTLSFromFile("cert.pem", "key.pem")
	if err != nil {
		return err
	}
	srv := grpc.NewServer(grpc.Creds(tls))
	lis, err := net.Listen("tcp", "crdhost.c.gifted-cooler-279818.internal:50001")
	if err != nil {
		return err
	}
	controller.RegisterProbeCommunicatorServer(srv, new(CommunicatorServer))
}

func makeCert() error {
	err := maker.Command("openssl", "req", "-x509", "-newkey", "rsa:4096", "-keyout", "key.pem", "-out", "cert.pem", "-days", "365", "-nodes", "-subj", "'/CN=crdhost.c.gifted-cooler-279818'").Run()
	if err != nil {
		return err
	}
	// Read public certificate so that it can be included in regional vm metadata
	cert, err := ioutil.ReadFile("cert.pem")
	if err != nil {
		return err
	}
	return nil
}

// Handles communication between controller and regional VMs
type CommunicatorServer struct{}

// Provides regional VMs with information about which probes to run
func (pc *CommunicatorServer) Register(ctx context.Context, in *controller.RegisterRequest) (*controller.RegisterResponse, error) {
	vm, ok := vms[in.Source]
	if !ok {
		//TODO(langenbahn): log this error
		log.Print("Register: given source does not correspond to an existing VM")
		return &controller.RegisterResponse{}
	}
	vm.setState(idle)
	vm.updatePingTime()
	return &controller.RegisterResponse{probes: &ProbeConfigs{Probe: vm.probes}, account: config.GetAccount(), pingInterval: config.GetVmPingInterval()}, nil
}

// Processes incoming information from probes
func (pc *CommunicatorServer) Ping(ctx context.Context, in *controller.Heartbeat) (*controller.Heartbeat, error) {
	if in.GetStop() {
		vms[in.GetSource()].restartVM()
	} else {
		vms[in.GetSource()].setState(probing)
	}
	vms[in.Source].updatePingTime()
	in.Source = "Controller"
	in.Stop = running
	return in, nil
}

func checkVMs(max time.Duration, wg *sync.WaitGroup) {
	for stoppedVMs < len(vms) {
		for _, v := range vms {
			vm.stateLock.Lock()
			if (vm.state == starting || vm.state == idle || vm.state == probing) && time.Now().After(v.lastPing.Add(max)) {
				vm.stateLock.Unlock()
				v.restartVM()
			} else {
				vm.stateLock.Unlock()
			}
		}
		//Have a channel here that this can block on?
		time.Sleep(1 * time.Minute)
	}
	wg.Done()
}
