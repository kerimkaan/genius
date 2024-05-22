/*
Copyright © 2024 Kerim Kaan Dönmez <kaan@kerimkaan.com>
*/
package cmd

import (
	"fmt"
	"log"
	"net"

	"github.com/miekg/dns"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get a brief system information",
	Long: `
Get a brief system information such as hostname, OS, architecture, kernel version, platform, platform family,
platform version, virtualization system, virtualization role, hostID, uptime, boot time, procs, load average, CPU model,
CPU cores, total memory, memory usage, total disk space, disk space used, disk space free, disk space used percentage,
disk filesystem, interface name, interface hardware address, interface MTU, IPv4 of en0, DNS servers.`,
	Run: func(cmd *cobra.Command, args []string) {
		hInfo, err := host.Info()
		if err != nil {
			log.Println(err)
			panic(err)
		}
		// Get the system cpu load average
		loadAvg, err := load.Avg()
		if err != nil {
			log.Println(err)
			panic(err)
		}

		cpuInfo, err := cpu.Info()
		if err != nil {
			log.Println(err)
			panic(err)
		}

		// Memory information
		memInfo, err := mem.VirtualMemory()
		if err != nil {
			log.Println(err)
			panic(err)
		}
		// Get disk usage information
		diskInfo, err := disk.Usage("/")
		if err != nil {
			log.Println(err)
			panic(err)
		}

		virtSystem := !(hInfo.VirtualizationSystem == "")
		virtName := hInfo.VirtualizationSystem
		virtRole := hInfo.VirtualizationRole
		if !virtSystem {
			virtName = "Not Available"
			virtRole = "Not Available"
		}
		fmt.Println("Hostname: ", hInfo.Hostname)
		fmt.Println("OS: ", hInfo.OS)
		fmt.Println("Architecture: ", hInfo.KernelArch)
		fmt.Println("Kernel Version: ", hInfo.KernelVersion)
		fmt.Println("Platform: ", hInfo.Platform)
		fmt.Println("Platform Family: ", hInfo.PlatformFamily)
		fmt.Println("Platform Version: ", hInfo.PlatformVersion)
		fmt.Println("Virtualization System: ", virtName)
		fmt.Println("Virtualization Role: ", virtRole)
		fmt.Println("HostID: ", hInfo.HostID)
		fmt.Println("Uptime: ", hInfo.Uptime)
		fmt.Println("Boot Time: ", hInfo.BootTime)
		fmt.Println("Procs: ", hInfo.Procs)

		fmt.Println("Load Average (1/5/15): ", loadAvg.Load1, loadAvg.Load5, loadAvg.Load15)

		fmt.Println("CPU Model: ", cpuInfo[0].ModelName)
		fmt.Println("CPU Cores: ", cpuInfo[0].Cores)

		fmt.Println("Total Memory: ", memInfo.Total/1024/1024, "MB")
		usedPercentWith2Digit := fmt.Sprintf("%.2f", memInfo.UsedPercent)
		fmt.Println("Memory usage (%): ", usedPercentWith2Digit)

		fmt.Println("Total Disk Space: ", diskInfo.Total/1024/1024/1024, "GB")
		fmt.Println("Disk Space Used: ", diskInfo.Used/1024/1024/1024, "GB")
		fmt.Println("Disk Space Free: ", diskInfo.Free/1024/1024/1024, "GB")
		fmt.Println("Disk Space Used (%): ", diskInfo.UsedPercent)
		fmt.Println("Disk filesystem: ", diskInfo.Fstype)

		ifaces, err := net.Interfaces()
		if err != nil {
			log.Println(err)
			panic(err)
		}
		for _, i := range ifaces {
			if i.Name == "en0" {
				addrs, err := i.Addrs()
				if err != nil {
					log.Println(err)
					panic(err)
				}
				for _, addr := range addrs {
					ipv4 := addr.(*net.IPNet).IP.To4()
					if ipv4 != nil {
						fmt.Println("Interface Name: ", i.Name)
						fmt.Println("Interface Hardware Address: ", i.HardwareAddr)
						fmt.Println("Interface MTU: ", i.MTU)
						fmt.Println("IPv4 of en0: ", ipv4)
					}
				}
			}
		}

		resolvConf, err := dns.ClientConfigFromFile("/etc/resolv.conf")
		if err != nil {
			log.Println(err)
			panic(err)
		}
		fmt.Println("DNS Servers: ", resolvConf.Servers)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
