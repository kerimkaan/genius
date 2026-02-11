/*
Copyright © 2024 Kerim Kaan Dönmez <kaan@kerimkaan.com>
*/
package cmd

import (
	"fmt"
	"genius/constants"
	"genius/helpers"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "genius",
	Version: constants.Version,
	Short:   "Genius is a CLI tool to get a brief system information.",
	Long:    `Genius is a CLI tool to get a brief system information such as platform, host, CPU, memory, disk, and network.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if helpers.IsWindows() {
			return fmt.Errorf("genius is not compatible with Windows")
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
