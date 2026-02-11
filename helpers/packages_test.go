package helpers

import (
	"testing"
)

func TestGetPythonVersion(t *testing.T) {
	version, err := GetPythonVersion()
	if err != nil {
		// python3 may not be installed; that's ok
		t.Skipf("Skipping: python3 not available: %v", err)
	}
	// Should either return "Not found" or a version string
	if version == "" {
		t.Error("GetPythonVersion() returned empty string")
	}
	if version != "Not found" {
		// Version should look like a semver (e.g. "3.12.1")
		if len(version) < 3 {
			t.Errorf("GetPythonVersion() returned unexpected value: %q", version)
		}
	}
}

func TestGetHomeBrewVersion(t *testing.T) {
	version, err := GetHomeBrewVersion()
	if err != nil {
		t.Skipf("Skipping: brew not available: %v", err)
	}
	// Should either return "Not found" or a version string
	if version == "" {
		t.Error("GetHomeBrewVersion() returned empty string")
	}
}

func TestRunCommandNotFound(t *testing.T) {
	result, err := runCommand("nonexistent-command-xyz-12345")
	if err != nil {
		// Some systems might return an error
		return
	}
	if result != "Not found" {
		t.Errorf("runCommand for nonexistent binary returned %q, want 'Not found'", result)
	}
}
