package probe

import (
	"os/exec"
	"strings"
)

func findDevice() (string, error) {
	d := exec.Command("emulator", "-list-avds")
	err := d.Run()
	if err != nil {
		return "", err
	}
	o := make([]byte, 0)
	_, err = d.Stdout.Write(o)
	if err != nil {
		return "", err
	}
	n := strings.Split(string(o), "\n")
	return n[0], nil
}

func startEmulator() error {
	dev, err := findDevice()
	if err != nil {
		return err
	}
	err = exec.Command("emulator", "-avd", dev).Run()
	if err != nil {
		return err
	}
	err = exec.Command("adb", "install",
		"../../FCMExternalProberTarget/app/build/outputs/apk/debug/app-debug.apk").Run()
	if err != nil {
		return err
	}
	err = exec.Command("adb", "shell", "am", "start", "-n",
		"com.google.firebase.messaging.testing.fcmexternalprobertarget/" +
		"com.google.firebase.messaging.testing.fcmexternalprobertarget.MainActivity").Run()
	if err != nil {
		return err
	}
	return nil
}

func getToken() (string, error) {
	com := exec.Command("bash", "receive", "token.txt")
	err := com.Run()
	if err != nil {
		return "", err
	}
	tok := make([]byte, 0)
	_, err = com.Stdout.Write(tok)
	if err != nil {
		return "", err
	}
	return string(tok), nil
}

func getMessage(t string) (string, error) {
	com := exec.Command("bash", "receive", "logs/" + t + ".txt")
	err := com.Run()
	if err != nil {
		return "", err
	}
	msg := make([]byte, 0)
	_, err = com.Stdout.Write(msg)
	if err != nil {
		return "", err
	}
	return string(msg), nil
}

func killEmulator() error {
	err := exec.Command("adb", "uninstall",
		"com.google.firebase.messaging.testing.fcmexternalprobertarget").Run()
	if err != nil {
		return err
	}
	err = exec.Command("adb", "emu", "kill").Run()
	if err != nil {
		return err
	}
	return nil
}