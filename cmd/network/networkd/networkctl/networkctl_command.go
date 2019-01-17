package networkctl

import (
	"os/exec"
	"strings"

	"github.com/RestGW/api-routerd/cmd/share"
)

// Response JSON response for networkctl
type Response struct {
	Index       string `json:"index"`
	Link        string `json:"link"`
	Type        string `json:"type"`
	Operational string `json:"operational"`
	Setup       string `json:"setup"`
}

// ResponseLLDP LLDP response
type ResponseLLDP struct {
	Link       string `json:"link"`
	ChassisID  string `json:"chassis_id"`
	System     string `json:"system"`
	Capability string `json:"capability"`
	Port       string `json:"port"`
	PortDdesc  string `json:"port_description"`
}

func getPathNetworkctl() (string, error) {
	err := share.CheckBinaryExists("networkctl")
	if err != nil {
		return "", err
	}

	path, err := exec.LookPath("networkctl")
	if err != nil {
		return "", err
	}

	return path, nil
}

// ExecuteNetworkctl execute networkctl same as status
func ExecuteNetworkctl() ([]Response, error) {
	path, err := getPathNetworkctl()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(path)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(stdout), "\n")

	r := make([]Response, len(lines)-4)
	for i, line := range lines {
		if i == 0 || i >= len(lines)-3 {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		n := Response{
			Index:       fields[0],
			Link:        fields[1],
			Type:        fields[2],
			Operational: fields[3],
			Setup:       fields[4],
		}

		r[i-1] = n
	}

	return r, nil
}

// ExecuteNetworkctlLLDP execute networkctl same as status
func ExecuteNetworkctlLLDP() ([]ResponseLLDP, error) {
	path, err := getPathNetworkctl()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(path, "lldp")
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(stdout), "\n")

	r := make([]ResponseLLDP, len(lines)-6)
	for i, line := range lines {
		if i == 0 || i >= len(lines)-4 {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}

		n := ResponseLLDP{
			Link:       fields[0],
			ChassisID:  fields[1],
			System:     fields[2],
			Capability: fields[3],
			Port:       fields[4],
			PortDdesc:  fields[5],
		}

		r[i-1] = n
	}

	return r, nil
}
