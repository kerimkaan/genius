package helpers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"os/exec"
)

// commandTimeout is the maximum duration for external command execution.
const commandTimeout = 5 * time.Second

// runCommand executes a command with a timeout and returns its trimmed stdout.
func runCommand(name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	stdout, err := cmd.Output()
	if err != nil {
		if strings.Contains(err.Error(), "executable file not found") {
			return "Not found", nil
		}
		return "", fmt.Errorf("failed to run %s: %w", name, err)
	}
	return strings.TrimSpace(string(stdout)), nil
}

// GetHomeBrewVersion returns the installed Homebrew version string.
func GetHomeBrewVersion() (string, error) {
	output, err := runCommand("brew", "--version")
	if err != nil {
		return "", err
	}
	// Remove the "Homebrew" prefix
	output = strings.ReplaceAll(output, "Homebrew", "")
	return strings.TrimSpace(output), nil
}

// GetPythonVersion returns the installed Python 3 version string.
func GetPythonVersion() (string, error) {
	output, err := runCommand("python3", "--version")
	if err != nil {
		return "", err
	}
	// Remove the "Python" prefix
	output = strings.ReplaceAll(output, "Python", "")
	return strings.TrimSpace(output), nil
}
