/*
Copyright © 2024 Kerim Kaan Dönmez <kaan@kerimkaan.com>
*/
package cmd

import (
	"fmt"
	"genius/helpers"

	"github.com/spf13/cobra"
)

var folderCount int

var largestFoldersCmd = &cobra.Command{
	Use:     "largest-folders",
	Aliases: []string{"lf"},
	Short:   "Lists the largest folders in the user's home directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		homeDir, err := helpers.GetUserHomeDir()
		if err != nil {
			return fmt.Errorf("could not get user's home directory: %w", err)
		}

		fmt.Printf("Top %d largest folders under %s:\n", folderCount, homeDir)
		largestFolders, err := helpers.GetLargestFoldersConcurrent(homeDir, folderCount)
		if err != nil {
			return fmt.Errorf("error while getting folder sizes: %w", err)
		}

		for i, folder := range largestFolders {
			sizeGB := float64(folder.Size) / (1024 * 1024 * 1024)
			fmt.Printf("  %d. %s - %.2f GB\n", i+1, folder.Path, sizeGB)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(largestFoldersCmd)
	largestFoldersCmd.Flags().IntVarP(&folderCount, "count", "n", 5, "Number of largest folders to display")
}
