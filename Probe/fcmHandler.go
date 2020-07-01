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
	"errors"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var fcmAuthToken string
var deadline = time.Now()

func getAuth() (string, error) {
	if time.Now().After(deadline) {
		err := prepareAuth()
		if err != nil {
			return "", err
		}
	}
	return fcmAuthToken, nil
}

func prepareAuth() error {
	get, err := exec.Command("curl",
		"http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/send-service-account@gifted-cooler-279818.iam.gserviceaccount.com/token",
		"-H", "Metadata-Flavor: Google").Output()
	if err != nil {
		return err
	}
	ret := regexp.MustCompile("[{\":,}]+").Split(string(get), -1)
	err = updateDeadline(&ret)
	if err != nil {
		return err
	}
	fcmAuthToken, err = findValue("access_token", &ret)
	if err != nil {
		return err
	}
	return nil
}

func updateDeadline(kv *[]string) error {
	ttl, err := findValue("expires_in", kv)
	if err != nil {
		return err
	}
	dur, err := time.ParseDuration(ttl + "s")
	if err != nil {
		return err
	}
	deadline = time.Now().Add(dur)
	return nil
}

func findValue(key string, list *[]string) (string, error) {
	// Assumes list is stored in a {key, :, value} configuration
	for i := 0; i < len(*list)-1; i++ {
		if strings.Compare(key, (*list)[i]) == 0 {
			return (*list)[i+1], nil
		}
	}
	return "", errors.New("probe/fcmHandler: key not found")
}

func sendMessage(time string) error {
	auth, err := getAuth()
	if err != nil {
		return err
	}
	err = exec.Command("bash", "send", "-d", deviceToken, "-a", auth, "-t", time).Run()
	if err != nil {
		return err
	}
	return nil
}
