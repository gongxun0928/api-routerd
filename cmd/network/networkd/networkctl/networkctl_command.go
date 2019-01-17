package networkctl

import (
	"os/exec"
	"strings"

	"github.com/RestGW/api-routerd/cmd/share"
)

// StatusResponse JSON response for networkctl
type StatusResponse struct {
	Index       string `json:"index"`
	Link        string `json:"link"`
	Type        string `json:"type"`
	Operational string `json:"operational"`
	Setup       string `json:"setup"`
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
func ExecuteNetworkctl() ([]StatusResponse, error) {
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

	r := make([]StatusResponse, len(lines)-4)
	for i, line := range lines {
		if i == 0 || i >= len(lines)-3 {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		n := StatusResponse{
			Index:       fields[0],
			Type:        fields[1],
			Operational: fields[2],
			Setup:       fields[3],
		}

		r[i-1] = n
	}

	return r, nil
}
