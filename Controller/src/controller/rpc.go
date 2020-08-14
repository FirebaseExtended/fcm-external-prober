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

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const certFile = "cert.pem"

var cert []byte

func initServer() error {
	err := makeCert()
	if err != nil {
		return err
	}
	tls, err := credentials.NewServerTLSFromFile("cert.pem", "key.pem")
	if err != nil {
		return err
	}
	srv := grpc.NewServer(grpc.Creds(tls))
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GetMetadata().GetPort()))
	if err != nil {
		return err
	}
	RegisterProbeCommunicatorServer(srv, new(CommunicatorServer))
	go srv.Serve(lis)
	return nil
}

func makeCert() error {
	err := maker.Command("openssl", "req",
		"-x509",
		"-newkey", "rsa:4096",
		"-keyout", "key.pem",
		"-out", certFile,
		"-days", "365",
		"-nodes",
		"-subj", "/CN="+config.Metadata.GetHostIp()).Run()
	if err != nil {
		logger.LogErrorf("%v", err)
		return err
	}
	// Read public certificate so that it can be included in regional vm metadata
	cert, err = ioutil.ReadFile(certFile)
	if err != nil {
		return err
	}
	config.GetMetadata().Cert = proto.String(string(cert))
	return nil
}

func addMetadata() error {
	md := proto.MarshalTextString(config.GetMetadata())
	err := maker.Command("gcloud", "compute", "project-info", "add-metadata", fmt.Sprintf("--metadata=probeData=%s", md)).Run()
	if err != nil {
		return err
	}
	return nil
}

// Handles communication between controller and regional VMs
type CommunicatorServer struct {
	UnimplementedProbeCommunicatorServer
}

// Provides regional VMs with information about which probes to run
func (cs *CommunicatorServer) Register(ctx context.Context, in *RegisterRequest) (*RegisterResponse, error) {
	vm, ok := vms[in.GetSource()]
	if !ok {
		logger.LogErrorf("Register: given source %s does not correspond to existing VM", in.GetSource())
		return &RegisterResponse{}, errors.New("invalid source")
	}
	vm.setState(idle)
	vm.updatePingTime()
	return &RegisterResponse{
		Probes:     &ProbeConfigs{Probe: vm.probes},
		Account:    config.GetMetadata().GetAccount(),
		PingConfig: config.GetPingConfig()}, nil
}

// Processes incoming information from probes
func (cs *CommunicatorServer) Ping(ctx context.Context, in *Heartbeat) (*Heartbeat, error) {
	if in.GetStop() {
		vms[in.GetSource()].restartVM()
	} else {
		vms[in.GetSource()].setState(probing)
	}
	vms[in.GetSource()].updatePingTime()
	src := "Controller"
	in.Source = &src
	in.Stop = &stopping
	return in, nil
}

func checkVMs(max time.Duration) {
	for stoppedVMs < len(vms) {
		for _, vm := range vms {

			if isTimedOut(vm, max) {
				vm.restartVM()
			}
		}
		time.Sleep(time.Duration(config.GetPingConfig().GetInterval()) * time.Minute)
	}
}

func isTimedOut(vm *regionalVM, max time.Duration) bool {
	vm.stateLock.Lock()
	defer vm.stateLock.Unlock()
	return (vm.state == starting || vm.state == idle || vm.state == probing) && clock.Now().After(vm.lastPing.Add(max))
}
