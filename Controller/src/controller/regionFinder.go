package controller

import (
	"os/exec"
	"regexp"
	"strings"
)

func getCompatZones(reqs []string) ([]string, error) {
	r, err := findZones()
	if err != nil {
		return nil, err
	}
	var ret []string
	for i := 0; i < len(r); i++ {
		str, err := listZoneInfo(r[i])
		if err != nil {
			return nil, err
		}
		if meetsRequirements(str, reqs) {
			ret = append(ret, str)
		}
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
	reg := regexp.MustCompile("[a-z]+-[a-z]+[1-9]-a")
	return reg.FindAllString(str, -1), nil
}

func listZoneInfo(zone string) (string, error) {
	out, err := exec.Command("gcloud", "compute", "regions", "describe", zone).Output()
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

func listZones() (string, error) {
	out, err := exec.Command("gcloud", "compute", "zones", "list").Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}