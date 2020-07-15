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
)

func TestProbe(t *testing.T) {
	probeType = "default"
	probeNumber = 1
	probeInterval = 0
	testTimes := []time.Time{time.Unix(1, 0), time.Unix(2, 0), time.Unix(2, 0)}
	clock = newTestClock(testTimes)
	at := "1111"
	ei := 1
	tt := "Bearer"
	j := fmt.Sprintf("{\"access_token\":\"%s\",\"expires_in\":%d,\"token_type\":\"%s\"}", at, ei, tt)
	c := []string{"TEST_DEVICE", "", "", "", "", "TEST_TOKEN", j, "", "2000", "", ""}
	var b []bool
	for i := 0; i < len(c); i++ {
		b = append(b, false)
	}
	maker = newFakeCommandMaker(c, b)
	fakeLogger := new(fakeLogger)
	pLog = fakeLogger
	probe()

	res := fakeLogger.testLogs
	if len(res) != 1 {
		t.Logf("TestProbe: incorrect number of probe resolution logs: actual: %d, expected: 1", len(res))
		t.Fail()
	}
	if res[0].time != testTimes[0].Format(timeLogFormat) || res[0].state != "resolved" || res[0].latency != 1000 {
		t.Logf("TestProbe: test probe not resolved correctly")
		t.Fail()
	}
}
