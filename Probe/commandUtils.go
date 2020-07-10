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

import "os/exec"

type CommandMaker interface {
	Command(name string, arg ...string) CommandRunner
}

type CommandRunner interface {
	Run() error
	Output() ([]byte, error)
	Start() error
}

type cmdMaker struct{}

func (e cmdMaker) Command(name string, arg ...string) CommandRunner {
	return newCMDRunner(exec.Command(name, arg...))
}

type cmdRunner struct {
	cmd exec.Cmd
}

func newCMDRunner(cmd *exec.Cmd) *cmdRunner {
	ret := new(cmdRunner)
	ret.cmd = *cmd
	return ret
}

func (c cmdRunner) Run() error {
	return c.cmd.Run()
}

func (c cmdRunner) Output() ([]byte, error) {
	return c.cmd.Output()
}

func (c cmdRunner) Start() error {
	return c.cmd.Start()
}
