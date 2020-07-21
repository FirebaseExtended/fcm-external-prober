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

import "log"

type Logger interface {
	LogProbe(sp *sentProbe, st string, lat int)
}

type CloudLogger struct {
	// TODO(langenbahn): Add fields when logging is implemented
}

func NewCloudLogger() *CloudLogger {
	return new(CloudLogger)
}

func (c *CloudLogger) LogProbe(sp *sentProbe, st string, lat int) {
	//TODO(langenbahn): Implement logging
	log.Printf("time: %s status: %s latency: %d", sp.sendTime.Format(timeLogFormat), st, lat)
}
