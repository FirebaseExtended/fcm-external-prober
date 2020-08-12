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
	"testing"

	"github.com/FirebaseExtended/fcm-external-prober/Controller/src/controller"
	"github.com/FirebaseExtended/fcm-external-prober/Probe/src/utils"
	"github.com/golang/protobuf/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TestClient struct {
}

func (tc *TestClient) Register(ctx context.Context, in *controller.RegisterRequest, opts ...grpc.CallOption) (*controller.RegisterResponse, error) {
	// Determine return values based on register request source string
	switch in.GetSource() {
	case "Exceeded":
		return nil, status.Error(codes.DeadlineExceeded, "exceeded")
	case "Unavailable":
		return nil, status.Error(codes.Unavailable, "unavailable")
	default:
		return nil, nil
	}
}

func (tc *TestClient) Ping(ctx context.Context, in *controller.Heartbeat, opts ...grpc.CallOption) (*controller.Heartbeat, error) {
	switch in.GetSource() {
	case "Exceeded":
		return nil, status.Error(codes.DeadlineExceeded, "exceeded")
	case "Unavailable":
		return nil, status.Error(codes.Unavailable, "unavailable")
	case "Stop":
		*in.Source = "Controller"
		*in.Stop = true
		return in, nil
	default:
		*in.Source = "Controller"
		*in.Stop = false
		return in, nil
	}
}

func initVars(host string, retries int32) {
	hostname = host
	client = new(TestClient)
	pingConfig = &controller.PingConfig{Retries: &retries}
	metadata = &controller.MetadataConfig{RegisterRetries: &retries}
}

func TestGetMetadata(t *testing.T) {
	ip := "TEST_IP"
	testMeta := &controller.MetadataConfig{HostIp: &ip}
	encodedMeta := proto.MarshalTextString(testMeta)
	testString := "testItem.key:   probeData\n" +
		"testItem.value: " + encodedMeta
	maker = utils.NewFakeCommandMaker([]string{testString}, []bool{false}, false)

	err := getMetadata()
	if err != nil {
		t.Logf("TestGetMetadata: error returned on valid input %v", err)
		t.FailNow()
	}
	if metadata.GetHostIp() != "TEST_IP" {
		t.Logf("TestGetMetadata: metadata unmarshalled incorrectly")
		t.Fail()
	}
}

func TestRegisterExpected(t *testing.T) {
	initVars("testHost", 1)

	_, err := register()

	if err != nil {
		t.Log("TestRegisterExpected: received non-nil error output on valid input")
		t.Fail()
	}
}

func TestRegisterDeadlineExceeded(t *testing.T) {
	initVars("Exceeded", 1)

	_, err := register()
	if err == nil {
		t.Log("TestRegisterDeadlineExceeded: received nil error on error case")
		t.FailNow()
	}

	st, ok := status.FromError(err)
	if !ok || st.Code() != codes.DeadlineExceeded {
		t.Logf("TestRegisterDeadlineExceeded: incorrect error returned %v", err)
		t.Fail()
	}
}

func TestRegisterUnavailable(t *testing.T) {
	initVars("Unavailable", 5)

	_, err := register()
	if err == nil {
		t.Log("TestRegisterUnavailable: received nil error on error case")
		t.FailNow()
	}

	if err.Error() != "register: maximum register retries exceeded" {
		t.Logf("TestRegisterUnavailable: incorrect error returned %v", err)
		t.Fail()
	}
}

func TestCommunicate(t *testing.T) {
	initVars("Stop", 1)

	err := communicate()
	if err != nil {
		t.Log("TestCommunicate: non-nil error returned on valid input")
		t.Fail()
	}
}

func TestPingServerExpected(t *testing.T) {
	initVars("testHost", 1)

	_, err := pingServer(false)
	if err != nil {
		t.Log("TestPingServerExpected: received non-nil error output on valid input")
		t.Fail()
	}
}

func TestPingServerExceeded(t *testing.T) {
	initVars("Exceeded", 1)

	_, err := pingServer(false)
	if err == nil {
		t.Log("TestPingServerDeadlineExceeded: received nil error on error case")
		t.FailNow()
	}

	st, ok := status.FromError(err)
	if !ok || st.Code() != codes.DeadlineExceeded {
		t.Logf("TestPingServerDeadlineExceeded: incorrect error returned %v", err)
		t.Fail()
	}
}

func TestPingServerUnavailable(t *testing.T) {
	initVars("Unavailable", 5)

	_, err := pingServer(false)
	if err == nil {
		t.Log("TestPingServerUnavailable: received nil error on error case")
		t.FailNow()
	}

	if err.Error() != "pingServer: maximum register retries exceeded" {
		t.Logf("TestRegisterUnavailable: incorrect error returned %v", err)
		t.Fail()
	}
}
