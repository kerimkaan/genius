package helpers

import (
	"fmt"
	"genius/types"
	"os"
	"strings"
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
