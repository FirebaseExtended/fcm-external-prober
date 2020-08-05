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
	"fmt"
	"log"
	"time"

	"github.com/FirebaseExtended/fcm-external-prober/Controller/src/controller"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var client controller.ProbeCommunicatorClient
var pingConfig *controller.PingConfig
var hostname string
var retries int32

func initClient() error {
	tls, err := credentials.NewClientTLSFromFile("cert.pem", "")
	if err != nil {
		return err
	}
	//TODO: determine through which mode the internal DNS address of the controller will be available
	//Might be best to just throw it in metadata with the cert
	conn, err := grpc.Dial("crdhost.c.gifted-cooler-279818.internal:50001", grpc.WithTransportCredentials(tls), grpc.WithBlock(), grpc.WithTimeout(15 * time.Second))
	if err != nil {
		return err
	}
	client = controller.NewProbeCommunicatorClient(conn)

	hostname, err = getHostname()
	if err != nil {
		return err
	}

	c := make(chan *controller.RegisterResponse)
	fmt.Println("Register")
	go register(c)

	cfg := <- c
	fmt.Println("Done Register")
	probeConfigs = cfg.GetProbes()
	account = cfg.GetAccount()
	pingConfig = cfg.GetPingConfig()
	return nil
}

func register(out chan *controller.RegisterResponse) {
	TESTHOSTNAME := "us-central1-a"
	req := &controller.RegisterRequest{Source: &TESTHOSTNAME}
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	v, err := client.Register(ctx, req)
	select {
	case <-ctx.Done():
		log.Fatalf("Register: unable to communicate with server: %v", ctx.Err())
	case out <-v:
		if err != nil {
			log.Fatalf("Register: error in registering with server: %v", err)
		}
	}
}

func communicate() error {
	c := make(chan *controller.Heartbeat)
	hb := new(controller.Heartbeat)
	stop := false;
	for !hb.GetStop() {
		fmt.Println("Ping")
		go pingServer(c, &stop)
		hb = <- c
		fmt.Println("Done Ping")
		// If the source returned from pingServer is this VM, the server did not respond after maximum retries
		if hb.GetSource() == hostname {
			return errors.New("Communicate: connection with server lost")
		}
		time.Sleep(time.Duration(pingConfig.GetInterval()) * time.Minute) 
	}
	return nil
}

func confirmStop() error {
	fmt.Println("Confirming Stop")
	c := make(chan *controller.Heartbeat)
	hb := new(controller.Heartbeat)
	stop := true
	go pingServer(c, &stop)
	hb = <- c
	if hb.GetSource() == hostname {
		return errors.New("ConfirmStop: failed to communicate stopping to server")
	}
	return nil
}

func pingServer(out chan *controller.Heartbeat, stop *bool) {
	TESTHOSTNAME := "us-central1-a"
	hb := &controller.Heartbeat{Stop: stop, Source: &TESTHOSTNAME}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pingConfig.GetTimeout()) * time.Second)
	defer cancel()

	v, err := client.Ping(ctx, hb)
	select {
	case <-ctx.Done():
		log.Printf("PingServer: connection error: %v, retries: %d", ctx.Err(), retries)
		retries++
		if retries > pingConfig.GetRetries() {
			out<-hb
		}
	case out <-v:
		if err != nil {
			log.Printf("PingServer: error in pinging server: %v", err)
		}
	}
}

func getHostname() (string, error) {
	n, err := maker.Command("hostname").Output()
	if err != nil {
		return "", err
	}
	return string(n), nil
}