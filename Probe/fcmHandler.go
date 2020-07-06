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

package main

import (
	"encoding/json"
	"os/exec"
	"time"
)

// Represents an authentication response, into which JSON can be parsed
type Auth struct {
	Token string `json:"access_token"`
	Ttl time.Duration `json:"expires_in,int"`
	TokenType string `json:"token_type"`
}

var fcmAuth Auth
var deadline = time.Now()

func getAuth() (string, error) {
	if time.Now().After(deadline) {
		err := prepareAuth()
		if err != nil {
			return "", err
		}
	}
	return fcmAuth.Token, nil
}

func prepareAuth() error {
	// GET request for authentication credentials for interacting with FCM and Cloud Logger
	get, err := exec.Command("curl",
		"http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/" + serviceAccount + "/token",
		"-H", "Metadata-Flavor: Google").Output()
	if err != nil {
		return err
	}
	err = json.Unmarshal(get, &fcmAuth)
	if err != nil {
		return err
	}
	updateDeadline()
	return nil
}

func updateDeadline() {
	deadline = time.Now().Add(fcmAuth.Ttl * time.Second)
}

func sendMessage(time string) error {
	auth, err := getAuth()
	if err != nil {
		return err
	}
	err = exec.Command("bash", "send", "-d", deviceToken, "-a", auth, "-t", time, "-p", projectID).Run()
	if err != nil {
		return err
	}
	return nil
}
