/*
Copyright © 2024 Kerim Kaan Dönmez <kaan@kerimkaan.com>
*/
package cmd

import (
	"fmt"
	"genius/helpers"
	"log"
	"net"
	"os"
	"time"

	"github.com/beevik/ntp"
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
	Use:     "info",
	Aliases: []string{"i"},
	Short:   "Get a brief system information",
	Long: `
Get a brief system information such as hostname, OS, architecture, kernel version, platform, platform family,
platform version, virtualization system, virtualization role, hostID, uptime, boot time, procs, load average, CPU model,
CPU cores, total memory, memory usage, total disk space, disk space used, disk space free, disk space used percentage,
disk filesystem, interface name, interface hardware address, interface MTU, IPv4 of en0, DNS servers.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		// Is the system Windows?
		// If it is, the program will exit.
		// Because the program is not compatible with Windows.
		if helpers.IsWindows() {
			fmt.Println("The program is not compatible with Windows.")
			os.Exit(0)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		t0 := time.Now()
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
		// Swap memory information
		swapInfo, err := mem.SwapMemory()
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

		// Get NTP configurations from /etc/ntp.conf
		ntpConfig, err := helpers.ReadNTPConfFile()
		if err != nil {
			log.Println(err)
			// panic(err)
		}

		virtSystem := !(hInfo.VirtualizationSystem == "")
		virtName := hInfo.VirtualizationSystem
		virtRole := hInfo.VirtualizationRole
		if !virtSystem {
			virtName = "Not Available"
			virtRole = "Not Available"
		}
		fmt.Println("============================================")
		fmt.Println("System Information")
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
		fmt.Println("Uptime: ", hInfo.Uptime/60/60/24, "days", hInfo.Uptime/60/60%24, "hours", hInfo.Uptime/60%60, "minutes")
		// BootTime is in Unix timestamp
		// Convert it to human readable format
		fmt.Println("Last Boot Time in Local Time: ", time.Unix(int64(hInfo.BootTime), 0).Format("2006-01-02 15:04:05"))
		fmt.Println("Procs: ", hInfo.Procs)
		fmt.Println("============================================")
		fmt.Println("CPU Model: ", cpuInfo[0].ModelName)
		fmt.Println("CPU Cores: ", cpuInfo[0].Cores)
		fmt.Println("Load Average (1/5/15): ", loadAvg.Load1, loadAvg.Load5, loadAvg.Load15)
		fmt.Println("============================================")
		fmt.Println("Total Memory: ", memInfo.Total/1024/1024, "MB")
		usedPercentWith2Digit := fmt.Sprintf("%.2f", memInfo.UsedPercent)
		fmt.Println("Memory usage (%): ", usedPercentWith2Digit)
		fmt.Println("Swap Total: ", swapInfo.Total/1024/1024, "MB")
		fmt.Println("Swap Used (%): ", fmt.Sprintf("%.2f", swapInfo.UsedPercent))
		fmt.Println("============================================")

		fmt.Println("Total Disk Space: ", diskInfo.Total/1024/1024/1024, "GB")
		fmt.Println("Disk Space Used: ", diskInfo.Used/1024/1024/1024, "GB")
		fmt.Println("Disk Space Free: ", diskInfo.Free/1024/1024/1024, "GB")
		fmt.Println("Disk Space Used (%): ", fmt.Sprintf("%.2f", diskInfo.UsedPercent))
		fmt.Println("Disk filesystem: ", diskInfo.Fstype)
		fmt.Println("============================================")

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
		fmt.Println("============================================")

		resolvConf, err := dns.ClientConfigFromFile("/etc/resolv.conf")
		if err != nil {
			log.Println(err)
			panic(err)
		}
		fmt.Println("DNS Servers: ", resolvConf.Servers)

		if *ntpConfig == nil {
			fmt.Println("NTP Configuration is not available")
		} else {
			for _, ntp := range *ntpConfig {
				fmt.Println("NTP Server: ", ntp.Server)
				fmt.Println("NTP Burst Option: ", ntp.IBurst)
			}
		}

		// Get the ntp pool time
		ntpTime, err := ntp.Time("0.tr.pool.ntp.org")
		if err != nil {
			log.Println("Error getting NTP time")
			log.Println(err)
		}
		// Get the current time of the system
		currentTime := time.Now()
		fmt.Println("Current Time of the System: ", currentTime.Format("2006-01-02 15:04:05"))
		fmt.Println("NTP (0.tr.pool.ntp.org) Time: ", ntpTime.Format("2006-01-02 15:04:05"))
		// Calculate the time difference between the system and NTP time
		timeDiff := ntpTime.Sub(currentTime)
		fmt.Println("Time Difference: ", timeDiff)
		fmt.Println("============================================")
		if helpers.IsMacOS() {
			brewVersion, err := helpers.GetHomeBrewVersion()
			if err != nil {
				log.Println(err)
			}
			fmt.Println("Homebrew Version: ", brewVersion)
		}

		pythonVersion, err := helpers.GetPythonVersion()
		if err != nil {
			log.Println(err)
		}
		fmt.Println("Python Version: ", pythonVersion)
		t1 := time.Now()
		fmt.Println("Time taken to get the system information: ", t1.Sub(t0))
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
