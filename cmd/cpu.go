/*
Copyright © 2024 Kerim Kaan Dönmez <kaan@kerimkaan.com>
*/
package cmd

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/spf13/cobra"
)

// cpuCmd represents the cpu command
var cpuCmd = &cobra.Command{
	Use:     "cpu",
	Aliases: []string{"c"},
	Short:   "Detailed CPU information",
	Long:    `Get detailed CPU information such as physical cores, logical cores, CPU model, CPU cores, CPU usage percentage and CPU times.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		physicalCores, err := cpu.Counts(false)
		if err != nil {
			return fmt.Errorf("failed to get physical core count: %w", err)
		}

		logicalCores, err := cpu.Counts(true)
		if err != nil {
			return fmt.Errorf("failed to get logical core count: %w", err)
		}

		cpuInfo, err := cpu.Info()
		if err != nil {
			return fmt.Errorf("failed to get CPU info: %w", err)
		}

		fmt.Println(separator)
		fmt.Println("CPU Information")
		fmt.Printf("  Physical Cores: %d\n", physicalCores)
		fmt.Printf("  Logical Cores:  %d\n", logicalCores)

		for i, info := range cpuInfo {
			fmt.Printf("  CPU %d: %s (%.0f MHz, %d cores)\n", i, info.ModelName, info.Mhz, info.Cores)
		}

		percents, err := cpu.Percent(0, true)
		if err != nil {
			fmt.Printf("  CPU Usage: not available (%s)\n", err)
		} else {
			fmt.Println("  CPU Usage per Core:")
			for i, p := range percents {
				fmt.Printf("    Core %d: %.2f%%\n", i, p)
			}
		}

		times, err := cpu.Times(true)
		if err != nil {
			fmt.Printf("  CPU Times: not available (%s)\n", err)
		} else {
			fmt.Println("  CPU Times per Core:")
			for i, t := range times {
				fmt.Printf("    Core %d: user=%.1f system=%.1f idle=%.1f\n",
					i, t.User, t.System, t.Idle)
			}
		}
		fmt.Println(separator)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cpuCmd)
}
