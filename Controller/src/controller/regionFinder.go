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
	"errors"
	"regexp"
	"strings"
)

func getCompatZones(reqs []string) ([]string, error) {
	r, err := findZones()
	if err != nil {
		return nil, err
	}
	var ret []string
	for _, s := range r {
		str, err := listZoneInfo(s)
		if err != nil {
			return nil, err
		}
		if meetsRequirements(str, reqs) {
			ret = append(ret, s)
		}
	}
	if len(ret) == 0 {
		return nil, errors.New("getCompatZones: no zones available with given requirements")
	}
	return ret, nil
}

func findZones() ([]string, error) {
	str, err := listZones()
	if err != nil {
		return nil, err
	}
	// Match all strings that represent zone names that end in -a, i.e. us-east1-a
	// For now, assume that only one zone is needed per region
	reg := regexp.MustCompile(".*-a\\b")
	return reg.FindAllString(str, -1), nil
}

func listZones() (string, error) {
	out, err := maker.Command("gcloud", "compute", "zones", "list").Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func listZoneInfo(zone string) (string, error) {
	out, err := maker.Command("gcloud", "compute", "zones", "describe", zone).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func meetsRequirements(inf string, req []string) bool {
	for i := 0; i < len(req); i++ {
		if !strings.Contains(inf, req[i]) {
			return false
		}
	}
	return true
}
