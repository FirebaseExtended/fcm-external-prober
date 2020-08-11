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
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Wrapper for logging
type Logger interface {
	LogProbe(sp *sentProbe, st string, lat int)
	LogError(desc string)
	LogFatal(desc string)
	LogFatalf(desc string, args ...interface{})
}

// Logger that sends logs to cloud logger via gcloud
type CloudLogger struct {
}

//TODO(langenbahn): add log/error destinations in config proto
type probeLog struct {
	sendTime  string
	probeType string
	latency   int
	state     string
	region    string
	token     string
}

type errorLog struct {
	desc   string
	region string
	token  string
}

// Log probe information to specified log
func (c *CloudLogger) LogProbe(sp *sentProbe, st string, lat int) {
	cl := &probeLog{sp.sendTime.Format(timeLogFormat), sp.probe.config.Type.String(), lat, st,
		sp.probe.config.GetRegion(), deviceToken}
	l, err := json.Marshal(cl)
	if err != nil {
		c.LogError(fmt.Sprintf("Unable to log probe: unable to marshal JSON: %v", err))
		return
	}

	err = maker.Command("gcloud", "logging", "write",
		"--payload-type=json", "PROBELOG", string(l)).Run()
	if err != nil {
		c.LogError(fmt.Sprintf("Unable to log probe: unable to send to server: %v", err))
	}
}

func (c *CloudLogger) LogFatal(desc string) {
	c.LogError(fmt.Sprintf("fatal: %s", desc))
	os.Exit(1)
}

func (c *CloudLogger) LogFatalf(desc string, args ...interface{}) {
	c.LogError(fmt.Sprintf("fatal: "+desc, args...))
	os.Exit(1)
}

// Log errors to specified log
func (c *CloudLogger) LogError(desc string) {
	//TODO(langenbahn) add any other useful information for errors
	el := &errorLog{desc, probeConfigs.Probe[0].GetRegion(), deviceToken}
	l, err := json.Marshal(el)
	//TODO(langenbahn): in the case that an error message cannot be sent to the server,
	//send through gRPC connection if available
	if err != nil {
		log.Printf("Unable to log error: unable to marshal JSON: %v", err)
		return
	}

	err = maker.Command("gcloud", "logging", "write",
		"--payload-type=json", "ERRORLOG", string(l)).Run()
	if err != nil {
		log.Printf("Unable to log probe: unable to send to server: %v", err)
	}
}
