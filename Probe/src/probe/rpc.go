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
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/FirebaseExtended/fcm-external-prober/Controller/src/controller"
	"github.com/golang/protobuf/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

const certFile = "cert.pem"

var (
	client     controller.ProbeCommunicatorClient
	pingConfig *controller.PingConfig
	hostname   string
	metadata   *controller.MetadataConfig
)

// Retrieve metadata from string manually from flattened format instead of using JSON unmarshalling
// because data is deeply nested, and unmarshalling JSON would require several nested structs or type assertions
func getProbeData(raw string) (*controller.MetadataConfig, error) {
	items := strings.Split(raw, "commonInstanceMetadata.items")
	var probeData []string
	for i, item := range items {
		// Search for item "probeData" key, manipulate the associated value at the next index
		if strings.Contains(item, "probeData") {
			// data will come in the form of: 'value: "DATA"'
			probeData = strings.SplitN(items[i+1], ": ", 2)
			// Remove trailing newline character from cert for unquoting
			probeData[1] = strings.TrimSuffix(probeData[1], "\n")
			var err error
			probeData[1], err = strconv.Unquote(probeData[1])
			if err != nil {
				return nil, err
			}
			break
		}
	}

	if probeData == nil || len(probeData) < 2 {
		return nil, errors.New("getProbeData: unable to parse probe metadata")
	}

	meta := new(controller.MetadataConfig)
	err := proto.UnmarshalText(probeData[1], meta)
	if err != nil {
		return nil, err
	}
	return meta, nil
}

func getMetadata() error {
	out, err := maker.Command("gcloud", "compute", "project-info", "describe",
		"--format=flattened(commonInstanceMetadata.items[])").Output()
	if err != nil {
		return err
	}
	metadata, err = getProbeData(string(out))
	if err != nil {
		return err
	}
	cf, err := os.Create("cert.pem")
	if err != nil {
		return err
	}
	_, err = cf.Write([]byte(metadata.GetCert()))
	if err != nil {
		return err
	}
	return nil
}

func initClient() error {
	tls, err := credentials.NewClientTLSFromFile(certFile, "")
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(metadata.GetRegisterTimeout())*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", metadata.GetHostIp(), metadata.GetPort()),
		grpc.WithTransportCredentials(tls), grpc.WithBlock())
	if err != nil {
		return err
	}
	client = controller.NewProbeCommunicatorClient(conn)

	cfg, err := register()
	if err != nil {
		return err
	}
	probeConfigs = cfg.GetProbes()
	pingConfig = cfg.GetPingConfig()
	return nil
}

func register() (*controller.RegisterResponse, error) {
	req := &controller.RegisterRequest{Source: hostname}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(metadata.GetRegisterTimeout())*time.Second)
	defer cancel()

	for i := 0; i < int(metadata.GetRegisterRetries()); i++ {
		cfg, err := client.Register(ctx, req)
		st := status.Convert(err)
		switch st.Code() {
		case codes.DeadlineExceeded:
			return nil, err
		case codes.OK:
			return cfg, nil
		default:
			time.Sleep(time.Duration(metadata.GetRegisterRetryInterval()) * time.Second)
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
	_, err := pingServer(true)
	if err != nil {
		return errors.New("ConfirmStop: failed to communicate stopping to server")
	}
	return nil
}

func pingServer(stop bool) (*controller.Heartbeat, error) {
	hb := &controller.Heartbeat{Stop: stop, Source: hostname}
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
	n, err := maker.Command("curl", "-H", "Metadata-Flavor:Google", "http://metadata.google.internal/computeMetadata/v1/instance/name").Output()
	if err != nil {
		return "", err
	}
	// Remove trailing newline from command output
	return strings.TrimSuffix(string(n), "\n"), nil
}
