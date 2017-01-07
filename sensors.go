package main

import (
	"bytes"
	"os/exec"
	"regexp"
	"strings"
)

func ReadSensors() (map[string]string, error) {
	cmd := exec.Command("sensors")
	var output bytes.Buffer
	cmd.Stdout = &output
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return parseOutput(output.String()), nil
}

func StripUnits(s string) string {
	if strings.HasSuffix(s, "°C") {
		return s[:len(s)-len("°C")]
	} else if strings.HasSuffix(s, " RPM") {
		return s[:len(s)-len(" RPM")]
	}
	return s
}

func parseOutput(out string) map[string]string {
	res := map[string]string{}
	lines := strings.Split(out, "\n")
	fieldExpr := regexp.MustCompile(`^(.*?):\s*(.*?)(\(|$)`)
	for _, line := range lines {
		match := fieldExpr.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		if match[1] == "Adapter" {
			continue
		}
		res[match[1]] = strings.TrimSpace(match[2])
	}
	return res
}
