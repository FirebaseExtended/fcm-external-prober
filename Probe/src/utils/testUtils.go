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

package utils

import (
	"errors"
	"time"
)

// Used along with FakeCommandRunner to provide a set of responses to CmdRunner execution in a given test case
type FakeCommandMaker struct {
	messages []string
	errors   []bool
	index    int
}

func NewFakeCommandMaker(msg []string, err []bool) *FakeCommandMaker {
	ret := new(FakeCommandMaker)
	ret.messages = msg
	ret.errors = err
	ret.index = 0
	return ret
}

func (c *FakeCommandMaker) Command(name string, arg ...string) CommandRunner {
	ret := NewFakeCommand(c.messages[c.index], c.errors[c.index])
	c.index++
	return ret
}

// Provide a response to a single CmdRunner execution
type FakeCommandRunner struct {
	message string
	isError bool
}

func NewFakeCommand(msg string, isErr bool) *FakeCommandRunner {
	ret := new(FakeCommandRunner)
	ret.message = msg
	ret.isError = isErr
	return ret
}

func (e *FakeCommandRunner) Run() error {
	if e.isError {
		return errors.New(e.message)
	}
	return nil
}

func (e *FakeCommandRunner) Start() error {
	if e.isError {
		return errors.New(e.message)
	}
	return nil
}

func (e *FakeCommandRunner) Output() ([]byte, error) {
	if e.isError {
		return []byte{}, errors.New(e.message)
	}
	return []byte(e.message), nil
}

// Provides a specified time object for each call to time.Now() in a given test case
type FakeClock struct {
	times []time.Time
	index int
}

func NewTestClock(times []time.Time) *FakeClock {
	ret := new(FakeClock)
	ret.times = times
	ret.index = 0
	return ret
}

func (t *FakeClock) Now() time.Time {
	ret := t.times[t.index]
	t.index++
	return ret
}

// Collects calls to LogProbe for later analysis
type FakeLogger struct {
	testLogs []TestLog
}

func (t *FakeLogger) LogProbe(tim string, st string, lat int) {
	t.testLogs = append(t.testLogs, TestLog{tim, st, lat})
}

type TestLog struct {
	time    string
	state   string
	latency int
}
