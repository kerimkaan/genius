package helpers

import (
	"runtime"
	"testing"
)

func TestIsWindows(t *testing.T) {
	result := IsWindows()
	expected := runtime.GOOS == "windows"
	if result != expected {
		t.Errorf("IsWindows() = %v, want %v", result, expected)
	}
}

func TestIsMacOS(t *testing.T) {
	result := IsMacOS()
	expected := runtime.GOOS == "darwin"
	if result != expected {
		t.Errorf("IsMacOS() = %v, want %v", result, expected)
	}
}

func TestIsLinux(t *testing.T) {
	result := IsLinux()
	expected := runtime.GOOS == "linux"
	if result != expected {
		t.Errorf("IsLinux() = %v, want %v", result, expected)
	}
}

func TestOSFunctionsAreMutuallyExclusive(t *testing.T) {
	// At least one should be true on any platform, and they shouldn't overlap
	count := 0
	if IsWindows() {
		count++
	}
	if IsMacOS() {
		count++
	}
	if IsLinux() {
		count++
	}
	// On standard platforms exactly one should be true
	// (could be 0 on exotic platforms like freebsd)
	if count > 1 {
		t.Errorf("Multiple OS checks returned true: windows=%v, macos=%v, linux=%v",
			IsWindows(), IsMacOS(), IsLinux())
	}
}
