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

// Used along with fakeCommandRunner to provide a set of responses to cmdRunner execution in a given test case
type fakeCommandMaker struct {
	messages []string
	errors   []bool
	index    int
}

func newFakeCommandMaker(msg []string, err []bool) *fakeCommandMaker {
	ret := new(fakeCommandMaker)
	ret.messages = msg
	ret.errors = err
	ret.index = 0
	return ret
}

func (c *fakeCommandMaker) Command(name string, arg ...string) CommandRunner {
	ret := newFakeCommand(c.messages[c.index], c.errors[c.index])
	c.index++
	return ret
}

// Provide a response to a single cmdRunner execution
type fakeCommandRunner struct {
	message string
	isError bool
}

func newFakeCommand(msg string, isErr bool) *fakeCommandRunner {
	ret := new(fakeCommandRunner)
	ret.message = msg
	ret.isError = isErr
	return ret
}

func (e *fakeCommandRunner) Run() error {
	if e.isError {
		return errors.New(e.message)
	}
	return nil
}

func (e *fakeCommandRunner) Start() error {
	if e.isError {
		return errors.New(e.message)
	}
	return nil
}

func (e *fakeCommandRunner) Output() ([]byte, error) {
	if e.isError {
		return []byte{}, errors.New(e.message)
	}
	return []byte(e.message), nil
}

// Provides a specified time object for each call to time.Now() in a given test case
type fakeClock struct {
	times []time.Time
	index int
}

func newTestClock(times []time.Time) *fakeClock {
	ret := new(fakeClock)
	ret.times = times
	ret.index = 0
	return ret
}

func (t *fakeClock) Now() time.Time {
	ret := t.times[t.index]
	t.index++
	return ret
}

// Collects calls to logProbe for later analysis
type fakeLogger struct {
	testLogs []testLog
}

func (t *fakeLogger) logProbe(tim string, st string, lat int) {
	t.testLogs = append(t.testLogs, testLog{tim, st, lat})
}

type testLog struct {
	time    string
	state   string
	latency int
}
