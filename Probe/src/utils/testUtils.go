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
	repeat   bool
}

func NewFakeCommandMaker(msg []string, err []bool, repeat bool) *FakeCommandMaker {
	return &FakeCommandMaker{msg, err, 0, repeat}
}

func (c *FakeCommandMaker) Command(name string, arg ...string) CommandRunner {
	if c.repeat {
		c.index++
		return NewFakeCommand(c.messages[0], c.errors[0])
	}
	ret := NewFakeCommand(c.messages[c.index], c.errors[c.index])
	c.index++
	return ret
}

func (c *FakeCommandMaker) TimesCalled() int {
	return c.index
}

// Provide a response to a single CmdRunner execution
type FakeCommandRunner struct {
	message string
	isError bool
}

func NewFakeCommand(msg string, isErr bool) *FakeCommandRunner {
	return &FakeCommandRunner{msg, isErr}
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
	times  []time.Time
	index  int
	repeat bool
}

func NewFakeClock(times []time.Time, repeat bool) *FakeClock {
	return &FakeClock{times, 0, repeat}
}

func (t *FakeClock) Now() time.Time {
	if t.repeat {
		t.index++
		return t.times[0]
	}
	ret := t.times[t.index]
	t.index++
	return ret
}

func (t *FakeClock) TimesCalled() int {
	return t.index
}

type FakeBoolClock struct {
	times []time.Time
	index int
	toChange *bool
}

func NewFakeBoolClock(times []time.Time, toChange *bool) *FakeBoolClock {
	return &FakeBoolClock{times, 0, toChange}
}

func (f *FakeBoolClock) Now() time.Time {
	if f.index == len(f.times) - 1 {
		*(f.toChange) = false
	}
	ret := f.times[f.index]
	f.index++
	return ret
}

func (f *FakeBoolClock) TimesCalled() int {
	return f.index
}
