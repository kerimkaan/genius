package helpers

import (
	"log"
	"os/exec"
	"strings"
)

func GetHomeBrewVersion() (string, error) {
	// Get the Homebrew version
	// brew --version
	cmd := exec.Command("brew", "--version")
	stdout, err := cmd.Output()
	if err != nil {
		if strings.Contains(err.Error(), "executable file not found") {
			return "Not found", nil
		}
		log.Println(err)
		// panic(err)
		return "", err
	}
	stringified := string(stdout)
	// Remove the new line character
	stringified = strings.TrimSuffix(stringified, "\n")
	// Remove the Homebrew keyword
	stringified = strings.ReplaceAll(stringified, "Homebrew", "")
	// Remove the whitespace
	stringified = strings.TrimSpace(stringified)
	return stringified, nil
}

func GetPythonVersion() (string, error) {
	// Get the Python version
	// python --version
	cmd := exec.Command("python3", "--version")
	stdout, err := cmd.Output()
	if err != nil {
		if strings.Contains(err.Error(), "executable file not found") {
			return "Not found", nil
		}
		log.Println(err)
		// panic(err)
		return "", err
	}
	stringified := string(stdout)
	// Remove the new line character
	stringified = strings.TrimSuffix(stringified, "\n")
	// Remove the Python keyword
	stringified = strings.ReplaceAll(stringified, "Python", "")
	// Remove the whitespace
	stringified = strings.TrimSpace(stringified)
	return stringified, nil
}
