package main

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/gookit/color"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

func getHeader() string {
	host, _ := host.Info()
	username := os.Getenv("USERNAME")

	return fmt.Sprintf("%s@%s", username, host.Hostname)
}

func printHeader() {
	host, _ := host.Info()
	username := os.Getenv("USERNAME")

	color.FgLightRed.Print(username)
	color.FgWhite.Print("@")
	color.FgLightRed.Println(host.Hostname)
}

func getOS() string {
	host, _ := host.Info()
	return host.Platform
}

func getKernel() string {
	host, _ := host.Info()

	if host.OS == "windows" {
		version := strings.Split(strings.Split(host.KernelVersion, " ")[0], ".")
		return fmt.Sprintf("%s.%s.%s", version[0], version[1], version[2])
	}

	return host.KernelVersion
}

func getUptime() string {
	host, _ := host.Info()
	days := host.Uptime / 60 / 60 / 24
	hours := host.Uptime / 60 / 60 % 24
	minutes := host.Uptime / 60 % 60
	seconds := host.Uptime % 60
	uptime := ""

	if days != 0 {
		uptime = fmt.Sprintf("%s%dd ", uptime, days)
	}

	if hours != 0 {
		uptime = fmt.Sprintf("%s%dh ", uptime, hours)
	}

	if minutes != 0 {
		uptime = fmt.Sprintf("%s%dm ", uptime, minutes)
	}

	if seconds != 0 {
		uptime = fmt.Sprintf("%s%ds ", uptime, seconds)
	}

	return uptime
}

func getMemory() string {
	mem, _ := mem.VirtualMemory()
	used := float64(mem.Used/1024/1024) / 1000
	free := float64(mem.Free/1024/1024) / 1000
	full := used + free
	percent := math.Round((used * 100) / full)

	return fmt.Sprintf("%.2fGiB / %.2fGiB (%.0f%%)", used, full, percent)
}

func getDisk(path string) string {
	disk, _ := disk.Usage(path)
	used := disk.Used / 1024 / 1024 / 1024
	full := disk.Total / 1024 / 1024 / 1024
	return fmt.Sprintf("%dGiB / %dGiB (%.0f%%)", used, full, math.Round(disk.UsedPercent))
}

// TODO: Add colors. Use ANSI Escape Color Codes
func alignText(title string, data string) {
	space := ""
	space_length := len(getHeader()) - len(title) - 1
	for i := 0; i < space_length; i++ {
		space += " "
	}
	fmt.Printf("%s%s%s\n", title, space, data)
}

func main() {
	printHeader()
	alignText("OS", getOS())
	alignText("Kernel", getKernel())
	alignText("Uptime", getUptime())
	alignText("RAM", getMemory())

	if len(os.Args) == 2 && os.Args[1] == "--disk" {
		space := ""
		space_length := len(getHeader()) - len("Disks") - 1
		for i := 0; i < space_length; i++ {
			space += " "
		}

		fmt.Printf("\nDisks%s", space)

		host, _ := host.Info()
		os := host.OS

		disks, _ := disk.Partitions(true)
		if os != "windows" {
			disks, _ = disk.Partitions(false)
		}

		for i, disk := range disks {
			if os != "windows" && disk.Fstype != "ext4" {
				return
			}

			if i > 0 {
				fmt.Printf("%s%s\n", space+space+" ", getDisk(disk.Mountpoint))
			} else {
				fmt.Printf("%s\n", getDisk(disk.Mountpoint))
			}
		}
	}
}
