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

/*
Package probe implements an FCM probe that will:
	Initialize a new emulator and target app
	Send a specified number of messages to the app
	Attempt to verify that the app received those messages
	Log findings to Cloud Logger
	Terminate emulator and target app when finished
*/
package probe

import (
	"log"
	"sync"
	"time"

	"github.com/FirebaseExtended/fcm-external-prober/Controller/src/controller"
)

const timeFileFormat = "2006-01-02-T150405.000-Z0700"
const timeLogFormat = time.UnixDate

type probe struct {
	config *controller.ProbeConfig
}

func newProbe(cfg *controller.ProbeConfig) *probe {
	ret := new(probe)
	ret.config = cfg
	return ret
}

func (p *probe) probe(pwg *sync.WaitGroup) {
	if p.config.GetType() == controller.ProbeType_UNSPECIFIED {
		for probing {
			tim := clock.Now()
			err := fcmAuth.sendMessage(tim.Format(timeFileFormat), int(p.config.GetType()))
			if err != nil {
				log.Printf("probe: unable to send message: %s", err.Error())
				continue
			}
			sp := newSentProbe(tim, p)
			addProbe(sp)
			// Time interval between probes
			time.Sleep(time.Duration(p.config.GetSendInterval()) * time.Second)
		}
	}
	pwg.Done()
}

func (p *probe) getTypeString() string {
	// Will need to update this function if additional probe types are added
	switch p.config.GetType() {
	case controller.ProbeType_UNSPECIFIED:
		return "default"
	case controller.ProbeType_TOPIC:
		return "topic"
	default:
		return ""
	}
}
