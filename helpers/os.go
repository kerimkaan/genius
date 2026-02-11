package helpers

import "runtime"

// IsWindows reports whether the current operating system is Windows.
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// IsMacOS reports whether the current operating system is macOS (darwin).
func IsMacOS() bool {
	return runtime.GOOS == "darwin"
}

// IsLinux reports whether the current operating system is Linux.
func IsLinux() bool {
	return runtime.GOOS == "linux"
}
