/*
 * Copyright 2020 Google LLC
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

package controller

import (
	"fmt"
	"log"
	"os"
)

// Wrapper for controller error logging to Cloud Logger
type Logger interface {
	LogFatal(desc string)
	LogFatalf(desc string, args ...interface{})
	LogError(desc string)
	LogErrorf(desc string, args ...interface{})
}

// Logs controller-based errors
type ControllerLogger struct {
	Destination string // Location in Cloud Logger to which errors are logged
}

// Log error and terminate program
func (c *ControllerLogger) LogFatal(desc string) {
	c.LogError(fmt.Sprintf("fatal: %s", desc))
	os.Exit(1)
}

// Log error with format and terminate program
func (c *ControllerLogger) LogFatalf(desc string, args ...interface{}) {
	c.LogError(fmt.Sprintf("fatal: "+desc, args...))
	os.Exit(1)
}

// Log error to Cloud Logger
func (c *ControllerLogger) LogError(desc string) {
	err := maker.Command("gcloud", "logging", "write", config.GetControllerLogDestination(), desc).Run()
	if err != nil {
		log.Printf("Unable to log error: unable to send to server: %v", err)
	}
}

// Log error with format
func (c *ControllerLogger) LogErrorf(desc string, args ...interface{}) {
	c.LogError(fmt.Sprintf(desc, args...))
}
