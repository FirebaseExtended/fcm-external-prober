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

package probe

import (
	"errors"
	"log"
	"strings"
	"time"
)

func findDevice() (string, error) {
	// Assumes the first AVD on the list is the one to be used
	out, err := maker.Command("emulator", "-list-avds").Output()
	if err != nil {
		return "", err
	}
	dev := strings.Split(string(out), "\n")
	return dev[0], nil
}

func startEmulator() error {
	dev, err := findDevice()
	if err != nil {
		return err
	}
	err = maker.Command("emulator", "-avd", dev, "-no-snapshot", "-delay-adb").Start()
	if err != nil {
		return err
	}
	err = maker.Command("adb", "wait-for-device").Run()
	if err != nil {
		return err
	}
	return nil
}

func startApp() error {
	err := maker.Command("adb", "install",
		"../../FCMExternalProberTarget/app/build/outputs/apk/debug/app-debug.apk").Run()
	if err != nil {
		return err
	}
	err = maker.Command("adb", "shell", "am", "start", "-n",
		"com.google.firebase.messaging.testing.fcmexternalprobertarget/"+
			"com.google.firebase.messaging.testing.fcmexternalprobertarget.MainActivity").Run()
	if err != nil {
		return err
	}
	return nil
}

func getToken() (string, error) {
	for i := 0; i < 10; i++ {
		tok, err := maker.Command("bash", "receive", "token.txt").Output()
		if err != nil {
			return "", err
		}
		if string(tok) != "nf" {
			return string(tok), nil
		}
		time.Sleep(5 * time.Second)
	}
	return "", errors.New("timed out on token generation")
}

func getMessage(fn string) (string, error) {
	log.Print(fn)
	msg, err := maker.Command("bash", "receive", fn+".txt", "-p", "logs/").Output()
	if err != nil {
		return "", err
	}
	return string(msg), nil
}

func uninstallApp() error {
	err := maker.Command("adb", "uninstall",
		"com.google.firebase.messaging.testing.fcmexternalprobertarget").Run()
	if err != nil {
		return err
	}
	return nil
}

func killEmulator() error {
	err := maker.Command("adb", "emu", "kill").Run()
	if err != nil {
		return err
	}
	return nil
}
