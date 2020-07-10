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
	"time"
)

// Used along with testCommand to provide a set of responses to command execution in a given test case
type testExecute struct {
	messages []string
	errors []bool
	index int
}

func newTestExecute(msg []string, err []bool) *testExecute {
	ret := new(testExecute)
	ret.messages = msg
	ret.errors = err
	ret.index = 0
	return ret
}

func (c *testExecute) Command(name string, arg ...string) Commander {
	ret := newTestCommand(c.messages[c.index], c.errors[c.index])
	c.index++
	return ret
}

// Provide a response to a single command execution
type testCommand struct {
	message string
	isError bool
}

func newTestCommand(msg string, isErr bool) *testCommand {
	ret := new(testCommand)
	ret.message = msg
	ret.isError = isErr
	return ret
}

func (e *testCommand) Run() error {
	if e.isError {
		return errors.New(e.message)
	}
	return nil
}

func (e *testCommand) Start() error {
	if e.isError {
		return errors.New(e.message)
	}
	return nil
}

func (e *testCommand) Output() ([]byte, error) {
	if e.isError {
		return []byte{}, errors.New(e.message)
	}
	return []byte(e.message), nil
}

// Provides a specified time object for each call to time.Now() in a given test case
type testClock struct {
	times []time.Time
	index int
}

func newTestClock(times []time.Time) *testClock {
	ret := new(testClock)
	ret.times = times
	ret.index = 0
	return ret
}

func (t *testClock) Now() time.Time {
	ret := t.times[t.index]
	t.index++
	return ret
}

// Collects calls to logProbe for later analysis
type testLogger struct {
	testLogs []testLog
}

func (t *testLogger) logProbe(tim string, st string, lat int) {
	t.testLogs = append(t.testLogs, testLog{tim, st, lat})
}

type testLog struct {
	time string
	state string
	latency int
}
