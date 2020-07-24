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

// Collects calls to LogProbe for later analysis
type fakeLogger struct {
	testLogs []testLog
}

func (t *fakeLogger) LogProbe(sp *sentProbe, st string, lat int) {
	t.testLogs = append(t.testLogs, testLog{sp.sendTime.Format(timeLogFormat), st, lat})
}

type testLog struct {
	time    string
	state   string
	latency int
}
