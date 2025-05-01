/*
Copyright © 2024 Kerim Kaan Dönmez <kaan@kerimkaan.com>
*/
package cmd

import (
	"fmt"
	"genius/helpers"

	"github.com/spf13/cobra"
)

var largestFoldersCmd = &cobra.Command{
	Use:   "largest-folders",
	Short: "Lists the largest folders in the user's home directory",
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, err := helpers.GetUserHomeDir()
		if err != nil {
			fmt.Println("Could not get user's home directory:", err)
			return
		}
		fmt.Printf("Top 5 largest folders under %s:\n", homeDir)
		largestFolders, err := helpers.GetLargestFoldersConcurrent(homeDir, 5)
		if err != nil {
			fmt.Println("Error while getting folder sizes:", err)
		} else {
			for i, folder := range largestFolders {
				sizeGB := float64(folder.Size) / (1024 * 1024 * 1024)
				fmt.Printf("%d. %s - %.2f GB\n", i+1, folder.Path, sizeGB)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(largestFoldersCmd)
}
