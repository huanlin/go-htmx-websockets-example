package hardware

import (
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

func GetSystemSection() (string, error) {
	runtimeOS := runtime.GOOS

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return "", err
	}

	hostStat, err := host.Info()
	if err != nil {
		return "", err
	}

	output := fmt.Sprintf("Hostname: %s\nTotal Memory: %d\nUsed Memory: %d\nOS: %s\n",
		hostStat.Hostname, vmStat.Total, vmStat.Used, runtimeOS)

	return output, nil
}

func GetCpuSection() (string, error) {
	cpuStat, err := cpu.Info()
	if err != nil {
		return "", err
	}

	output := fmt.Sprintf("CPU Model: %s\nCores: %d", cpuStat[0].ModelName, cpuStat[0].Cores)

	return output, nil
}

func GetDiskSection() (string, error) {
	diskStat, err := disk.Usage("/")
	if err != nil {
		return "", err
	}

	output := fmt.Sprintf("Total Disk Space: %d\nFree Disk Space: %d", diskStat.Total, diskStat.Free)

	return output, nil
}
