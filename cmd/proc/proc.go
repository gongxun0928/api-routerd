// SPDX-License-Identifier: Apache-2.0

package proc

import (
	"api-routerd/cmd/share"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	log "github.com/sirupsen/logrus"
)

const (
	ProcMiscPath    = "/proc/misc"
	ProcNetArpPath  = "/proc/net/arp"
	ProcModulesPath = "/proc/modules"
)

type NetARP struct {
	IPAddress string `json:"ip_address"`
	HWType    string `json:"hw_type"`
	Flags     string `json:"flags"`
	HWAddress string `json:"hw_address"`
	Mask      string `json:"mask"`
	Device    string `json:"device"`
}

type Modules struct {
	Module     string `json:"module"`
	MemorySize string `json:"memory_size"`
	Instances  string `json:"instances"`
	Dependent  string `json:"dependent"`
	State      string `json:"state"`
}

func GetVersion(rw http.ResponseWriter) error {
	infostat, err := host.Info()
	if err != nil {
		return err
	}

	return share.JsonResponse(infostat, rw)
}

func GetUserStat(rw http.ResponseWriter) error {
	userstat, err := host.Users()
	if err != nil {
		return err
	}

	return share.JsonResponse(userstat, rw)
}

func GetTemperatureStat(rw http.ResponseWriter) error {
	tempstat, err := host.SensorsTemperatures()
	if err != nil {
		return err
	}

	return share.JsonResponse(tempstat, rw)
}

func GetNetStat(rw http.ResponseWriter, protocol string) error {
	conn, err := net.Connections(protocol)
	if err != nil {
		return err
	}

	return share.JsonResponse(conn, rw)
}

func GetNetStatPid(rw http.ResponseWriter, protocol string, process string) error {
	pid, err := strconv.ParseInt(process, 10, 32)
	if err != nil || protocol == "" || pid == 0 {
		return errors.New("Can't parse request")
	}

	conn, err := net.ConnectionsPid(protocol, int32(pid))
	if err != nil {
		return err
	}

	return share.JsonResponse(conn, rw)
}

func GetProtoCountersStat(rw http.ResponseWriter) error {
	protocols := []string{"ip", "icmp", "icmpmsg", "tcp", "udp", "udplite"}

	proto, err := net.ProtoCounters(protocols)
	if err != nil {
		return err
	}

	return share.JsonResponse(proto, rw)
}

func GetNetDev(rw http.ResponseWriter) error {
	netdev, err := net.IOCounters(true)
	if err != nil {
		return err
	}

	return share.JsonResponse(netdev, rw)
}

func GetInterfaceStat(rw http.ResponseWriter) error {
	interfaces, err := net.Interfaces()
	if err != nil {
		return err
	}

	return share.JsonResponse(interfaces, rw)
}

func GetSwapMemoryStat(rw http.ResponseWriter) error {
	swap, err := mem.SwapMemory()
	if err != nil {
		return err
	}

	return share.JsonResponse(swap, rw)
}

func GetVirtualMemoryStat(rw http.ResponseWriter) error {
	virt, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	return share.JsonResponse(virt, rw)
}

func GetCPUInfo(rw http.ResponseWriter) error {
	cpus, err := cpu.Info()
	if err != nil {
		return err
	}

	return share.JsonResponse(cpus, rw)
}

func GetCPUTimeStat(rw http.ResponseWriter) error {
	cpus, err := cpu.Times(true)
	if err != nil {
		return err
	}

	return share.JsonResponse(cpus, rw)
}

func GetAvgStat(rw http.ResponseWriter) error {
	avgstat, r := load.Avg()
	if r != nil {
		return r
	}

	return share.JsonResponse(avgstat, rw)
}

func GetMisc(rw http.ResponseWriter) error {
	lines, err := share.ReadFullFile(ProcMiscPath)
	if err != nil {
		log.Fatalf("Failed to read: %s", ProcMiscPath)
		return errors.New("Failed to read misc")
	}

	miscMap := make(map[int]string)
	for _, line := range lines {
		fields := strings.Fields(line)

		deviceNum, err := strconv.Atoi(fields[0])
		if err != nil {
			continue
		}
		miscMap[deviceNum] = fields[1]
	}

	return share.JsonResponse(miscMap, rw)
}

func GetNetArp(rw http.ResponseWriter) error {
	lines, err := share.ReadFullFile(ProcNetArpPath)
	if err != nil {
		log.Fatalf("Failed to read: %s", ProcNetArpPath)
		return errors.New("Failed to read /proc/net/arp")
	}

	netarp := make([]NetARP, len(lines)-1)
	for i, line := range lines {
		if i == 0 {
			continue
		}

		fields := strings.Fields(line)

		arp := NetARP{}
		arp.IPAddress = fields[0]
		arp.HWType = fields[1]
		arp.Flags = fields[2]
		arp.HWAddress = fields[3]
		arp.Mask = fields[4]
		arp.Device = fields[5]
		netarp[i-1] = arp
	}

	return share.JsonResponse(netarp, rw)
}

// GetModules Get all installed modules
func GetModules(rw http.ResponseWriter) error {
	lines, err := share.ReadFullFile(ProcModulesPath)
	if err != nil {
		log.Fatalf("Failed to read: %s", ProcModulesPath)
		return errors.New("Failed to read /proc/modules")
	}

	modules := make([]Modules, len(lines))
	for i, line := range lines {
		fields := strings.Fields(line)

		module := Modules{}

		for j, field := range fields {
			switch j {
			case 0:
				module.Module = field
				break
			case 1:
				module.MemorySize = field
				break
			case 2:
				module.Instances = field
				break
			case 3:
				module.Dependent = field
				break
			case 4:
				module.State = field
				break
			}
		}

		modules[i] = module
	}

	return share.JsonResponse(modules, rw)
}

func GetProcessInfo(rw http.ResponseWriter, proc string, property string) error {
	pid, err := strconv.ParseInt(proc, 10, 32)
	if err != nil {
		return err
	}

	p, err := process.NewProcess(int32(pid))
	if err != nil {
		return err
	}

	switch property {
	case "pid-connections":
		conn, err := p.Connections()
		if err != nil {
			return err
		}

		return share.JsonResponse(conn, rw)

	case "pid-rlimit":
		rlimit, err := p.Rlimit()
		if err != nil {
			return err
		}

		return share.JsonResponse(rlimit, rw)

	case "pid-rlimit-usage":
		rlimit, err := p.RlimitUsage(true)
		if err != nil {
			return err
		}

		return share.JsonResponse(rlimit, rw)

	case "pid-status":
		s, err := p.Status()
		if err != nil {
			return err
		}

		return share.JsonResponse(s, rw)

	case "pid-username":
		u, err := p.Username()
		if err != nil {
			return err
		}

		return share.JsonResponse(u, rw)

	case "pid-open-files":
		f, err := p.OpenFiles()
		if err != nil {
			return err
		}

		return share.JsonResponse(f, rw)

	case "pid-fds":
		f, err := p.NumFDs()
		if err != nil {
			return err
		}

		return share.JsonResponse(f, rw)

	case "pid-name":
		n, err := p.Name()
		if err != nil {
			return err
		}

		return share.JsonResponse(n, rw)

	case "pid-memory-percent":
		m, err := p.MemoryPercent()
		if err != nil {
			return err
		}

		return share.JsonResponse(m, rw)

	case "pid-memory-maps":
		m, err := p.MemoryMaps(true)
		if err != nil {
			return err
		}

		return share.JsonResponse(m, rw)

	case "pid-memory-info":
		m, err := p.MemoryInfo()
		if err != nil {
			return err
		}

		return share.JsonResponse(m, rw)

	case "pid-io-counters":
		m, err := p.IOCounters()
		if err != nil {
			return err
		}

		return share.JsonResponse(m, rw)
	}

	return nil
}
