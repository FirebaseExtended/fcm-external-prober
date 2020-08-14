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
	SetRegion(reg string)
	SetError(dest string)
	SetLog(dest string)
	LogProbe(sp *sentProbe, st string, lat int, tok string)
	LogError(desc string)
	LogErrorf(desc string, args ...interface{})
	LogFatal(desc string)
	LogFatalf(desc string, args ...interface{})
}

// Logger that sends logs to cloud logger via gcloud
type CloudLogger struct {
	Region    string // Region in which VM is located
	LogDest   string // Location in Cloud Logger where probes should be logged
	ErrorDest string // Location in Cloud Logger where errors should be logged
}

func NewCloudLogger() *CloudLogger {
	return &CloudLogger{"unspecified", "defaultLog", "defaultError"}
}

type probeLog struct {
	SendTime  string `json:"sendTime"`  // Send time of probe
	ProbeType string `json:"probeType"` // Type of probe
	Latency   int    `json:"latency"`   // Latency of probe
	State     string `json:"state"`     // State of probe
	Region    string `json:"region"`    // Region in which VM is located
	Token     string `json:"token"`     // Device token of device on which the app is located
}

type errorLog struct {
	Desc   string `json:"description"` // Error description
	Region string `json:"region"`      // Region in which VM is located
}

// Set the region in which VM is located
func (c *CloudLogger) SetRegion(reg string) {
	c.Region = reg
}

// Set error log location
func (c *CloudLogger) SetError(dest string) {
	c.ErrorDest = dest
}

// Set probe log location
func (c *CloudLogger) SetLog(dest string) {
	c.LogDest = dest
}

// Log probe information to specified log
func (c *CloudLogger) LogProbe(sp *sentProbe, st string, lat int, tok string) {
	cl := &probeLog{sp.sendTime.Format(timeLogFormat), sp.probe.config.Type.String(), lat, st,
		c.Region, tok}
	l, err := json.Marshal(cl)
	if err != nil {
		c.LogError(fmt.Sprintf("Unable to log probe: unable to marshal JSON: %v", err))
		return
	}

	err = maker.Command("gcloud", "logging", "write",
		"--payload-type=json", c.LogDest, string(l)).Run()
	if err != nil {
		c.LogError(fmt.Sprintf("Unable to log probe: unable to send to server: %v", err))
	}
}

// Log errors to specified log
func (c *CloudLogger) LogError(desc string) {
	//TODO(langenbahn) add any other useful information for errors
	el := &errorLog{desc, c.Region}
	l, err := json.Marshal(el)
	//TODO(langenbahn): in the case that an error message cannot be sent to the server,
	//send through gRPC connection if available
	if err != nil {
		log.Printf("Unable to log error: unable to marshal JSON: %v", err)
		return
	}

	err = maker.Command("gcloud", "logging", "write",
		"--payload-type=json", c.ErrorDest, string(l)).Run()
	if err != nil {
		log.Printf("Unable to log error: unable to send to server: %v", err)
	}
}

func (c *CloudLogger) LogErrorf(desc string, args ...interface{}) {
	c.LogError(fmt.Sprintf(desc, args...))
}

// Log error and terminate
func (c *CloudLogger) LogFatal(desc string) {
	c.LogError(fmt.Sprintf("fatal: %s", desc))
	os.Exit(1)
}

// Log error with format and terminate
func (c *CloudLogger) LogFatalf(desc string, args ...interface{}) {
	c.LogError(fmt.Sprintf("fatal: "+desc, args...))
	os.Exit(1)
}
