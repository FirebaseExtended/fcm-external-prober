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
	"testing"

	"github.com/FirebaseExtended/fcm-external-prober/Probe/src/utils"
)

func TestGetCompatZones(t *testing.T) {
	testStrings := []string{"us-east1-a\nus-east1-b\nus-east2-abc\nus-east3-a\tus-east3\tUP",
		"Item1\nItem2\nItem3\nItem4", "Item1\nItem2\nItem3"}
	reqs := []string{"Item1", "Item4"}
	expectedOut := []string{"us-east1-a"}
	maker = utils.NewFakeCommandMaker(testStrings, []bool{false, false, false}, false)

	res, err := getCompatZones(reqs)

	if err != nil {
		t.Logf("TestGetCompatZones: returned error on valid input")
		t.Fail()
	}
	if len(res) != len(expectedOut) {
		t.Logf("TestGetCompatZones: incorrect output length: actual: %d, expected: %d", len(res), len(expectedOut))
		t.FailNow()
	}
	for i, s := range expectedOut {
		if res[i] != s {
			t.Logf("TestGetCompatZones: incorrect output value: actual: %s, expected: %s", res[i], s)
			t.Fail()
		}
	}
}

func TestFindZones(t *testing.T) {
	testString := "us-east1-a\nus-east1-b\nus-east2-abc\nus-east3-a\tus-east3\tUP"
	expectedOut := []string{"us-east1-a", "us-east3-a"}
	maker = utils.NewFakeCommandMaker([]string{testString}, []bool{false}, false)
	res, err := findZones()
	if err != nil {
		t.Log("TestFindZones: returned error on valid input")
		t.Fail()
	}
	if len(res) != len(expectedOut) {
		t.Logf("TestFindZones: incorrect output length: actual: %d, expected: %d", len(res), len(expectedOut))
		t.FailNow()
	}
	for i, s := range expectedOut {
		if s != res[i] {
			t.Logf("TestFindZones: incorrect output value: actual: %s, expected: %s", res[i], s)
			t.Fail()
		}
	}
}

func TestMeetsRequirements(t *testing.T) {
	testString := "Item1\nItem2\nItem3\nItem4"
	testReqs := []string{"Item1", "Item2", "Item4"}
	res := meetsRequirements(testString, testReqs)
	if !res {
		t.Log("TestMeetsRequirements: returned false on true input")
		t.Fail()
	}
}

func TestMeetsRequirementsNegative(t *testing.T) {
	testString := "Item1\nItem2\nItem3\nItem4"
	testReqs := []string{"Item1", "Item4", "Item5"}
	res := meetsRequirements(testString, testReqs)
	if res {
		t.Log("TestMeetsRequirementsNegative: returned true on false input")
		t.Fail()
	}
}
