package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"

	"github.com/shirou/gopsutil/disk"
)

type LuctusStat struct {
	Serverip        string  `json:"serverip"`
	CpuIdle         int     `json:"cpuidle"`
	CpuSteal        float64 `json:"cpusteal"`
	CpuIowait       float64 `json:"cpuiowait"`
	RamTotal        int     `json:"ramtotal"`
	RamUsed         int     `json:"ramused"`
	RamFree         int     `json:"ramfree"`
	DiskTotal       int     `json:"disktotal"`
	DiskUsed        int     `json:"diskused"`
	DiskFree        int     `json:"diskfree"`
	DiskPercentUsed int     `json:"diskpercentused"`
}

// Because fuck IPv6
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()

}

func main() {
	for {
		//RAM
		memory, err := memory.Get()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return
		}
		//CPU
		before, err := cpu.Get()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return
		}
		time.Sleep(10 * time.Minute)
		after, err := cpu.Get()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return
		}
		//Disk
		s, err := disk.Usage("/")
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return
		}

		totalCpu := float64(after.Total - before.Total)

		toSend := LuctusStat{
			Serverip:        GetOutboundIP(),
			CpuIdle:         int(math.Round(float64(after.Idle-before.Idle) / totalCpu * 100)),
			CpuSteal:        float64(after.Iowait-before.Iowait) / totalCpu * 100,
			CpuIowait:       float64(after.Steal-before.Steal) / totalCpu * 100,
			RamTotal:        int(math.Round(float64(memory.Total) / (1024 * 1024 * 1024))),
			RamUsed:         int(math.Round(float64(memory.Used) / (1024 * 1024 * 1024))),
			RamFree:         int(math.Round(float64(memory.Free) / (1024 * 1024 * 1024))),
			DiskPercentUsed: int(s.UsedPercent),
			DiskFree:        int(math.Round(float64(s.Free) / 1024 / 1024 / 1024)),
			DiskTotal:       int(math.Round(float64(s.Total) / 1024 / 1024 / 1024)),
			DiskUsed:        int(math.Round(float64(s.Used) / 1024 / 1024 / 1024)),
		}

		fmt.Println(toSend)
		toSendJson, err := json.Marshal(toSend)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return
		}
		//fmt.Println(toSendJson)
		resp, err := http.Post("http://localhost:7077/linuxstat", "application/json", bytes.NewBuffer(toSendJson))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return
		}
		fmt.Println(string(body))
	}
}
