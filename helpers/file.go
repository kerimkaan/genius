package helpers

import (
	"fmt"
	"genius/types"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// CheckFileExists checks if the file exists in the given path
// and returns a boolean value.
func CheckFileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// ReadNTPConfFile reads the NTP configuration file and returns the content.
// The NTP configuration file is located at /etc/ntp.conf
// and it contains the NTP server addresses.
// The NTP server addresses are used to get the time from the NTP server.
func ReadNTPConfFile() (*[]types.NTPConfiguration, error) {
	ntpConfFile := "/etc/ntp.conf"
	// Or if it has chrony installed
	// ntpConfFile := "/etc/chrony/chrony.conf"

	if !CheckFileExists(ntpConfFile) {
		return nil, fmt.Errorf("NTP configuration file %s does not exist.", ntpConfFile)
	}
	// Read the NTP configuration file
	// and return the content
	ntpFile, err := os.ReadFile(ntpConfFile)
	if err != nil {
		return nil, err
	}
	stringNTPFile := string(ntpFile)

	// We have NTP config file something like this:
	// server          0.us.pool.ntp.org               iburst
	// server          1.us.pool.ntp.org               iburst
	// server          2.us.pool.ntp.org               iburst
	// server          3.us.pool.ntp.org               iburst

	// Find the server addresses in the NTP configuration file
	// and return them
	if strings.Index(stringNTPFile, "server") == -1 {
		return nil, fmt.Errorf("NTP server addresses not found in %s", ntpConfFile)
	}
	// Remove the comments and get the server addresses
	if strings.Index(stringNTPFile, "#") != -1 {
		stringNTPFile = stringNTPFile[:strings.Index(stringNTPFile, "#")] // Remove the comments
	}
	stringNTPFile = stringNTPFile[strings.Index(stringNTPFile, "server")+7:] // Get the server addresses
	// Replace the "server" keyword with an empty string
	stringNTPFile = strings.ReplaceAll(stringNTPFile, "server", "")

	stringNTPFile = strings.ReplaceAll(stringNTPFile, "\n", " ")
	stringArrayNTPFile := strings.Fields(stringNTPFile)

	var ntpConfig []types.NTPConfiguration
	for i, server := range stringArrayNTPFile {
		if server == "iburst" {
			continue
		} else if server == "" {
			continue
		}
		var iburst bool
		if i+1 >= len(stringArrayNTPFile) {
			iburst = false
		} else {
			iburst = stringArrayNTPFile[i+1] == "iburst"
		}
		ntpConfig = append(ntpConfig, types.NTPConfiguration{
			Server: server,
			IBurst: iburst,
		})
	}
	return &ntpConfig, nil
}

// FolderSize holds folder path and its size
type FolderSize struct {
	Path string
	Size int64
}

// GetLargestFolders returns the largest n folders under the given root directory
func GetLargestFolders(root string, n int) ([]FolderSize, error) {
	folders := make(map[string]int64)

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && path != root {
			var size int64
			filepath.Walk(path, func(fp string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					size += info.Size()
				}
				return nil
			})
			folders[path] = size
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Convert map to slice
	var folderSizes []FolderSize
	for k, v := range folders {
		folderSizes = append(folderSizes, FolderSize{Path: k, Size: v})
	}
	// Sort by size descending
	sort.Slice(folderSizes, func(i, j int) bool {
		return folderSizes[i].Size > folderSizes[j].Size
	})
	if len(folderSizes) > n {
		folderSizes = folderSizes[:n]
	}
	return folderSizes, nil
}

// GetLargestFoldersConcurrent returns the largest n folders under the given root directory using concurrency
func GetLargestFoldersConcurrent(root string, n int) ([]FolderSize, error) {
	dirs, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	folderSizesCh := make(chan FolderSize, len(dirs))

	for _, dir := range dirs {
		if dir.IsDir() {
			wg.Add(1)
			go func(d os.DirEntry) {
				defer wg.Done()
				path := filepath.Join(root, d.Name())
				size := int64(0)
				filepath.Walk(path, func(fp string, info os.FileInfo, err error) error {
					if err != nil {
						return nil // Permission denied gibi hataları atla
					}
					if !info.IsDir() {
						size += info.Size()
					}
					return nil
				})
				folderSizesCh <- FolderSize{Path: path, Size: size}
			}(dir)
		}
	}

	wg.Wait()
	close(folderSizesCh)

	var folderSizes []FolderSize
	for fs := range folderSizesCh {
		folderSizes = append(folderSizes, fs)
	}

	sort.Slice(folderSizes, func(i, j int) bool {
		return folderSizes[i].Size > folderSizes[j].Size
	})
	if len(folderSizes) > n {
		folderSizes = folderSizes[:n]
	}
	return folderSizes, nil
}

// GetUserHomeDir returns the current user's home directory
func GetUserHomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return home, nil
}
