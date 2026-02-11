/*
Copyright © 2024 Kerim Kaan Dönmez <kaan@kerimkaan.com>
*/
package cmd

import (
	"bufio"
	"fmt"
	"genius/helpers"
	"net"
	"os"
	"strings"
	"time"

	"github.com/beevik/ntp"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/spf13/cobra"
)

const separator = "============================================"

// Default NTP pool server; can be overridden via --ntp-server flag.
var ntpServer string

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"i"},
	Short:   "Get a brief system information",
	Long: `Get a brief system information such as hostname, OS, architecture, kernel version,
platform, platform family, platform version, virtualization system, virtualization role,
hostID, uptime, boot time, procs, load average, CPU model, CPU cores, total memory,
memory usage, total disk space, disk space used, disk space free, disk space used percentage,
disk filesystem, network interfaces, DNS servers.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		t0 := time.Now()

		if err := printHostInfo(); err != nil {
			return err
		}
		if err := printCPUInfo(); err != nil {
			return err
		}
		if err := printMemoryInfo(); err != nil {
			return err
		}
		if err := printDiskInfo(); err != nil {
			return err
		}
		if err := printNetworkInfo(); err != nil {
			return err
		}
		if err := printDNSInfo(); err != nil {
			return err
		}
		printNTPInfo()
		printNTPTimeComparison(ntpServer)
		printPackageVersions()

		fmt.Println(separator)
		fmt.Printf("Time taken: %s\n", time.Since(t0))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().StringVar(&ntpServer, "ntp-server", "pool.ntp.org", "NTP server to compare time against")
}

// printHostInfo prints system and host information.
func printHostInfo() error {
	hInfo, err := host.Info()
	if err != nil {
		return fmt.Errorf("failed to get host info: %w", err)
	}

	virtName := hInfo.VirtualizationSystem
	virtRole := hInfo.VirtualizationRole
	if virtName == "" {
		virtName = "Not Available"
		virtRole = "Not Available"
	}

	fmt.Println(separator)
	fmt.Println("System Information")
	fmt.Printf("  Hostname:              %s\n", hInfo.Hostname)
	fmt.Printf("  OS:                    %s\n", hInfo.OS)
	fmt.Printf("  Architecture:          %s\n", hInfo.KernelArch)
	fmt.Printf("  Kernel Version:        %s\n", hInfo.KernelVersion)
	fmt.Printf("  Platform:              %s\n", hInfo.Platform)
	fmt.Printf("  Platform Family:       %s\n", hInfo.PlatformFamily)
	fmt.Printf("  Platform Version:      %s\n", hInfo.PlatformVersion)
	fmt.Printf("  Virtualization System: %s\n", virtName)
	fmt.Printf("  Virtualization Role:   %s\n", virtRole)
	fmt.Printf("  HostID:                %s\n", hInfo.HostID)
	fmt.Printf("  Uptime:                %d days %d hours %d minutes\n",
		hInfo.Uptime/60/60/24, hInfo.Uptime/60/60%24, hInfo.Uptime/60%60)
	fmt.Printf("  Last Boot Time:        %s\n",
		time.Unix(int64(hInfo.BootTime), 0).Format("2006-01-02 15:04:05"))
	fmt.Printf("  Procs:                 %d\n", hInfo.Procs)
	return nil
}

// printCPUInfo prints CPU model, core count and load average.
func printCPUInfo() error {
	cpuInfo, err := cpu.Info()
	if err != nil {
		return fmt.Errorf("failed to get CPU info: %w", err)
	}

	loadAvg, err := load.Avg()
	if err != nil {
		return fmt.Errorf("failed to get load average: %w", err)
	}

	fmt.Println(separator)
	fmt.Println("CPU")
	if len(cpuInfo) > 0 {
		fmt.Printf("  Model:                 %s\n", cpuInfo[0].ModelName)
		fmt.Printf("  Cores:                 %d\n", cpuInfo[0].Cores)
	}
	fmt.Printf("  Load Average (1/5/15): %.2f %.2f %.2f\n", loadAvg.Load1, loadAvg.Load5, loadAvg.Load15)
	return nil
}

// printMemoryInfo prints RAM and swap memory statistics.
func printMemoryInfo() error {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("failed to get memory info: %w", err)
	}

	swapInfo, err := mem.SwapMemory()
	if err != nil {
		return fmt.Errorf("failed to get swap info: %w", err)
	}

	fmt.Println(separator)
	fmt.Println("Memory")
	fmt.Printf("  Total Memory:  %d MB\n", memInfo.Total/1024/1024)
	fmt.Printf("  Memory Usage:  %.2f%%\n", memInfo.UsedPercent)
	fmt.Printf("  Swap Total:    %d MB\n", swapInfo.Total/1024/1024)
	fmt.Printf("  Swap Usage:    %.2f%%\n", swapInfo.UsedPercent)
	return nil
}

// printDiskInfo prints disk usage information for the root filesystem.
func printDiskInfo() error {
	diskInfo, err := disk.Usage("/")
	if err != nil {
		return fmt.Errorf("failed to get disk info: %w", err)
	}

	fmt.Println(separator)
	fmt.Println("Disk")
	fmt.Printf("  Total:      %d GB\n", diskInfo.Total/1024/1024/1024)
	fmt.Printf("  Used:       %d GB\n", diskInfo.Used/1024/1024/1024)
	fmt.Printf("  Free:       %d GB\n", diskInfo.Free/1024/1024/1024)
	fmt.Printf("  Usage:      %.2f%%\n", diskInfo.UsedPercent)
	fmt.Printf("  Filesystem: %s\n", diskInfo.Fstype)
	return nil
}

// printNetworkInfo prints information for all active, non-loopback network interfaces.
func printNetworkInfo() error {
	ifaces, err := net.Interfaces()
	if err != nil {
		return fmt.Errorf("failed to get network interfaces: %w", err)
	}

	fmt.Println(separator)
	fmt.Println("Network Interfaces")

	found := false
	for _, iface := range ifaces {
		// Skip loopback and down interfaces
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			ipv4 := ipNet.IP.To4()
			if ipv4 == nil {
				continue
			}
			found = true
			fmt.Printf("  Interface: %s\n", iface.Name)
			fmt.Printf("    MAC:  %s\n", iface.HardwareAddr)
			fmt.Printf("    MTU:  %d\n", iface.MTU)
			fmt.Printf("    IPv4: %s\n", ipv4)
		}
	}
	if !found {
		fmt.Println("  No active network interfaces found")
	}
	return nil
}

// printDNSInfo reads /etc/resolv.conf and prints DNS server addresses.
func printDNSInfo() error {
	fmt.Println(separator)
	fmt.Println("DNS")

	servers, err := parseDNSServers("/etc/resolv.conf")
	if err != nil {
		fmt.Printf("  DNS Servers: not available (%s)\n", err)
		return nil
	}
	fmt.Printf("  DNS Servers: %s\n", strings.Join(servers, ", "))
	return nil
}

// parseDNSServers parses nameserver entries from a resolv.conf file.
func parseDNSServers(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open %s: %w", path, err)
	}
	defer file.Close()

	var servers []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "nameserver") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				servers = append(servers, fields[1])
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(servers) == 0 {
		return nil, fmt.Errorf("no nameserver entries found")
	}
	return servers, nil
}

// printNTPInfo reads and prints local NTP configuration.
func printNTPInfo() {
	fmt.Println(separator)
	fmt.Println("NTP Configuration")

	ntpConfig, err := helpers.ReadNTPConfFile()
	if err != nil {
		fmt.Printf("  %s\n", err)
		return
	}
	for _, entry := range ntpConfig {
		fmt.Printf("  Server: %s (iburst: %t)\n", entry.Server, entry.IBurst)
	}
}

// printNTPTimeComparison queries an NTP server and shows the time difference.
func printNTPTimeComparison(server string) {
	ntpTime, err := ntp.Time(server)
	if err != nil {
		fmt.Printf("  NTP time query failed: %s\n", err)
		return
	}

	currentTime := time.Now()
	timeDiff := ntpTime.Sub(currentTime)

	fmt.Printf("  System Time:       %s\n", currentTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("  NTP Time (%s): %s\n", server, ntpTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("  Time Difference:   %s\n", timeDiff)
}

// printPackageVersions prints the versions of common software packages.
func printPackageVersions() {
	fmt.Println(separator)
	fmt.Println("Software Versions")

	if helpers.IsMacOS() {
		brewVersion, err := helpers.GetHomeBrewVersion()
		if err != nil {
			fmt.Printf("  Homebrew: error (%s)\n", err)
		} else {
			fmt.Printf("  Homebrew: %s\n", brewVersion)
		}
	}

	pythonVersion, err := helpers.GetPythonVersion()
	if err != nil {
		fmt.Printf("  Python:   error (%s)\n", err)
	} else {
		fmt.Printf("  Python:   %s\n", pythonVersion)
	}
}
