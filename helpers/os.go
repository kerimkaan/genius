package helpers

import "runtime"

func IsWindows() bool {
	// GOOS is the running program's operating system target:
	// one of darwin, freebsd, linux, and so on.
	if runtime.GOOS == "windows" {
		return true
	}
	return false
}

func IsMacOS() bool {
	if runtime.GOOS == "darwin" {
		return true
	}
	return false
}
