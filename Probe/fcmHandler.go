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
	"time"
)

// Represents an authentication response, into which JSON can be parsed
type Auth struct {
	Token     string        `json:"access_token"`
	Ttl       time.Duration `json:"expires_in,int"`
	TokenType string        `json:"token_type"`
	deadline  time.Time
}

func (a *Auth) getToken() (string, error) {
	if clock.Now().After(a.deadline) {
		err := a.prepareAuth()
		if err != nil {
			return "", err
		}
	}
	return a.Token, nil
}

func (a *Auth) prepareAuth() error {
	// GET request for authentication credentials for interacting with FCM and Cloud Logger
	get, err := exe.Command("curl",
		"http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/"+serviceAccount+"/token",
		"-H", "Metadata-Flavor: Google").Output()
	if err != nil {
		return err
	}
	err = json.Unmarshal(get, a)
	if err != nil {
		return err
	}
	a.updateDeadline()
	return nil
}

func (a *Auth) updateDeadline() {
	a.deadline = clock.Now().Add(a.Ttl * time.Second)
}

func (a *Auth) sendMessage(time string) error {
	auth, err := a.getToken()
	if err != nil {
		return err
	}
	err = exe.Command("bash", "send", "-d", deviceToken, "-a", auth, "-t", time, "-p", projectID).Run()
	if err != nil {
		return err
	}
	return nil
}
