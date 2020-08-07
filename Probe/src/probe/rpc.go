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
	"context"
	"errors"
	"google.golang.org/grpc/status"
	"time"

	"github.com/FirebaseExtended/fcm-external-prober/Controller/src/controller"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
)

var client controller.ProbeCommunicatorClient
var pingConfig *controller.PingConfig
var hostname string

func getMetadata() {
	//TODO(langenbahn): Implement this
	// This function will query the VM metadata for:
	// DNS address of the controller
	// Port on which the controller server is listening
	// TLS certificate for authenticating connection to controller
	// Configuration for
}


func initClient() error {
	tls, err := credentials.NewClientTLSFromFile("cert.pem", "")
	if err != nil {
		return err
	}
	//TODO(langenbahn): make the internal DNS address of the controller, port, and tls certificate available to the probe
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, ":", grpc.WithTransportCredentials(tls), grpc.WithBlock())
	if err != nil {
		return err
	}
	client = controller.NewProbeCommunicatorClient(conn)

	hostname, err = getHostname()
	if err != nil {
		return err
	}

	cfg, err := register()
	if err != nil {
		return err
	}
	probeConfigs = cfg.GetProbes()
	account = cfg.GetAccount()
	pingConfig = cfg.GetPingConfig()
	return nil
}

func register() (*controller.RegisterResponse, error) {
	req := &controller.RegisterRequest{Source: &hostname}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pingConfig.GetTimeout()))
	defer cancel()

	for i := 0; i < int(pingConfig.GetRetries()); i++ {
		cfg, err := client.Register(ctx, req)
		st := status.Convert(err)
		switch st.Code() {
		case codes.DeadlineExceeded:
			return nil, err
		case codes.OK:
			return cfg, nil
		default:
			time.Sleep(time.Duration(pingConfig.GetRetryInterval()))
		}
	}
	return nil, errors.New("register: maximum register retries exceeded")
}

func communicate() error {
	stop := false
	for !stop {
		hb, err := pingServer(stop)
		if err != nil {
			return err
		}
		stop = hb.GetStop()
		time.Sleep(time.Duration(pingConfig.GetInterval()) * time.Minute)
	}
	return nil
}

func confirmStop() error {
	// Probe is ceasing to run, so server response doesn't matter
	_, err := pingServer(false)
	if err != nil {
		return errors.New("ConfirmStop: failed to communicate stopping to server")
	}
	return nil
}

func pingServer(stop bool) (*controller.Heartbeat, error) {
	hb := &controller.Heartbeat{Stop: &stop, Source: &hostname}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pingConfig.GetTimeout())*time.Second)
	defer cancel()

	for i := 0; i < int(pingConfig.GetRetries()); i++ {
		hb, err := client.Ping(ctx, hb)
		st := status.Convert(err)
		switch st.Code() {
		case codes.DeadlineExceeded:
			return nil, err
		case codes.OK:
			return hb, nil
		default:
			time.Sleep(time.Duration(pingConfig.GetRetryInterval()))
		}
	}
	return nil, errors.New("pingServer: maximum register retries exceeded")

}

func getHostname() (string, error) {
	n, err := maker.Command("hostname").Output()
	if err != nil {
		return "", err
	}
	return string(n), nil
}
