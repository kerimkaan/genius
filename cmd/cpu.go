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
	Run: func(cmd *cobra.Command, args []string) {
		physicalCores, _ := cpu.Counts(false)

		logicalCores, _ := cpu.Counts(true)

		cpuInfo, err := cpu.Info()
		if err != nil {
			fmt.Println(err)
		}

		percents, err := cpu.Percent(0, true)
		if err != nil {
			if err.Error() == "not implemented yet" {
				percents = nil
			}
		}

		times, err := cpu.Times(true)
		if err != nil {
			if err.Error() == "not implemented yet" {
				times = []cpu.TimesStat{}
			}
		}

		fmt.Println("Physical cores:", physicalCores)
		fmt.Println("Logical cores:", logicalCores)
		fmt.Println("CPU Info:", cpuInfo)
		if percents == nil {
			fmt.Println("CPU Percent: Not implemented yet")
		} else {
			fmt.Println("CPU Percent:", percents)
		}
		if len(times) == 0 {
			fmt.Println("CPU Times: Not implemented yet")
		} else {
			fmt.Println("CPU Times:", times)
		}
	},
}

func init() {
	rootCmd.AddCommand(cpuCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cpuCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cpuCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
