package helpers

import (
	"bufio"
	"fmt"
	"genius/types"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// CheckFileExists checks if the file exists in the given path
// and returns a boolean value.
func CheckFileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// ntpConfPaths lists the possible NTP configuration file locations.
var ntpConfPaths = []string{
	"/etc/ntp.conf",
	"/etc/chrony/chrony.conf",
	"/etc/chrony.conf",
}

// ReadNTPConfFile reads the NTP configuration file and returns the parsed
// NTP server entries. It checks multiple known config file paths
// (/etc/ntp.conf, /etc/chrony/chrony.conf, /etc/chrony.conf).
func ReadNTPConfFile() ([]types.NTPConfiguration, error) {
	var confFile string
	for _, path := range ntpConfPaths {
		if CheckFileExists(path) {
			confFile = path
			break
		}
	}
	if confFile == "" {
		return nil, fmt.Errorf("no NTP configuration file found (checked: %s)", strings.Join(ntpConfPaths, ", "))
	}

	file, err := os.Open(confFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", confFile, err)
	}
	defer file.Close()

	var configs []types.NTPConfiguration
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Remove inline comments
		if idx := strings.Index(line, "#"); idx >= 0 {
			line = strings.TrimSpace(line[:idx])
		}

		// Skip empty lines and non-server lines
		if line == "" || (!strings.HasPrefix(line, "server") && !strings.HasPrefix(line, "pool")) {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		server := fields[1]
		iburst := false
		for _, f := range fields[2:] {
			if f == "iburst" {
				iburst = true
				break
			}
		}
		configs = append(configs, types.NTPConfiguration{
			Server: server,
			IBurst: iburst,
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading %s: %w", confFile, err)
	}

	if len(configs) == 0 {
		return nil, fmt.Errorf("no NTP server entries found in %s", confFile)
	}

	return configs, nil
}

// maxConcurrentWalkers limits the number of concurrent goroutines for folder scanning.
const maxConcurrentWalkers = 10

// GetLargestFoldersConcurrent returns the largest n folders under the given
// root directory using bounded concurrency.
func GetLargestFoldersConcurrent(root string, n int) ([]types.FolderSize, error) {
	dirs, err := os.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", root, err)
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrentWalkers)
	folderSizesCh := make(chan types.FolderSize, len(dirs))

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		wg.Add(1)
		go func(d os.DirEntry) {
			defer wg.Done()
			sem <- struct{}{}        // acquire
			defer func() { <-sem }() // release

			path := filepath.Join(root, d.Name())
			size := int64(0)
			_ = filepath.Walk(path, func(fp string, info os.FileInfo, err error) error {
				if err != nil {
					return nil // skip permission errors
				}
				if !info.IsDir() {
					size += info.Size()
				}
				return nil
			})
			folderSizesCh <- types.FolderSize{Path: path, Size: size}
		}(dir)
	}

	wg.Wait()
	close(folderSizesCh)

	var folderSizes []types.FolderSize
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

// GetUserHomeDir returns the current user's home directory.
func GetUserHomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return home, nil
}
