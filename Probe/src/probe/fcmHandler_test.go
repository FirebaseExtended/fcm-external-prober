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
	"fmt"
	"testing"
	"time"

	"github.com/FirebaseExtended/fcm-external-prober/Probe/src/utils"
)

func TestGetTokenAfterDeadline(t *testing.T) {
	tTime := time.Unix(0, 1)
	clock = utils.NewFakeClock([]time.Time{tTime, tTime}, false)
	at := "1111"
	ei := 1
	tt := "Bearer"
	s := fmt.Sprintf("{\"access_token\":\"%s\",\"expires_in\":%d,\"token_type\":\"%s\"}", at, ei, tt)
	maker = utils.NewFakeCommandMaker([]string{s}, []bool{false}, false)
	tAuth := new(Auth)

	str, err := tAuth.getToken()

	if err != nil {
		t.Logf("testGetTokenAfterDeadline: error on valid input: + %s", err.Error())
		t.Fail()
	}
	if str != at {
		t.Logf("testGetTokenAfterDeadline: token not updated after deadline")
		t.Fail()
	}
}

func TestGetTokenBeforeDeadline(t *testing.T) {
	tTime := time.Unix(0, 1)
	clock = utils.NewFakeClock([]time.Time{tTime}, false)
	tAuth := new(Auth)
	tAuth.Token = "TEST_TOKEN"
	tAuth.deadline = time.Unix(0, 2)

	str, err := tAuth.getToken()

	if err != nil {
		t.Logf("testGetTokenBeforeDeadline: error on valid input: + %s", err.Error())
		t.Fail()
	}
	if str != "TEST_TOKEN" {
		t.Logf("testGetTokenBeforeDeadline: token updated before deadline")
		t.Fail()
	}
}

func TestPrepareAuth(t *testing.T) {
	tTime := time.Unix(100, 0)
	cTime := tTime.Add(1 * time.Second)
	clock = utils.NewFakeClock([]time.Time{tTime}, false)
	at := "1111"
	ei := 1
	tt := "Bearer"
	s := fmt.Sprintf("{\"access_token\":\"%s\",\"expires_in\":%d,\"token_type\":\"%s\"}", at, ei, tt)
	maker = utils.NewFakeCommandMaker([]string{s}, []bool{false}, false)
	tAuth := new(Auth)

	err := tAuth.prepareAuth()

	if err != nil {
		t.Logf("testPrepareAuth: error on valid input: + %s", err.Error())
		t.Fail()
	}
	if tAuth.Token != at || tAuth.Ttl != time.Duration(ei)*time.Nanosecond || tAuth.TokenType != tt {
		t.Logf("testPrepareAuth: token response JSON incorrectly parsed or modified")
		t.Fail()
	}
	if !cTime.Equal(tAuth.deadline) {
		t.Logf("testPrepareAuth: deadline not correctly updated: actual: %s, expected: %s", tAuth.deadline, cTime)
		t.Fail()
	}
}

func TestPrepareAuthCommandFail(t *testing.T) {
	s := "COMMAND_ERROR"
	maker = utils.NewFakeCommandMaker([]string{s}, []bool{true}, false)
	tAuth := new(Auth)

	err := tAuth.prepareAuth()

	if err == nil {
		t.Log("testPrepareAuthCommandFail: no error on failed cmdRunner")
		t.FailNow()
	}
	if err.Error() != s {
		t.Logf("testPrepareAuthCommandFail: unexpected error: actual: %s, expected: %s", err.Error(), s)
	}
}

func TestPrepareAuthInvalidJSON(t *testing.T) {
	s := "INVALID_JSON"
	maker = utils.NewFakeCommandMaker([]string{s}, []bool{true}, false)
	tAuth := new(Auth)

	err := tAuth.prepareAuth()
	if err == nil {
		t.Log("testPrepareAuthCommandFail: no error on invalid JSON")
		t.Fail()
	}
}
