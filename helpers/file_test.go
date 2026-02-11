package helpers

import (
	"genius/types"
	"os"
	"path/filepath"
	"testing"
)

func TestCheckFileExists(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "genius-test-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	t.Run("existing file", func(t *testing.T) {
		if !CheckFileExists(tmpFile.Name()) {
			t.Errorf("CheckFileExists(%q) = false, want true", tmpFile.Name())
		}
	})

	t.Run("non-existing file", func(t *testing.T) {
		if CheckFileExists("/nonexistent/path/to/file.txt") {
			t.Error("CheckFileExists for non-existent path should return false")
		}
	})
}

func TestReadNTPConfFile_NoFile(t *testing.T) {
	// Override ntpConfPaths for testing
	origPaths := ntpConfPaths
	ntpConfPaths = []string{"/nonexistent/ntp.conf"}
	defer func() { ntpConfPaths = origPaths }()

	_, err := ReadNTPConfFile()
	if err == nil {
		t.Error("ReadNTPConfFile() should return error when no config file exists")
	}
}

func TestReadNTPConfFile_ValidConfig(t *testing.T) {
	// Create a temp NTP config file
	tmpFile, err := os.CreateTemp("", "ntp-conf-test-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := `# NTP configuration
# This is a comment
driftfile /var/lib/ntp/ntp.drift

server 0.us.pool.ntp.org iburst
server 1.us.pool.ntp.org iburst
server 2.us.pool.ntp.org
# server 3.us.pool.ntp.org iburst
pool ntp.ubuntu.com iburst
`
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	// Override ntpConfPaths for testing
	origPaths := ntpConfPaths
	ntpConfPaths = []string{tmpFile.Name()}
	defer func() { ntpConfPaths = origPaths }()

	configs, err := ReadNTPConfFile()
	if err != nil {
		t.Fatalf("ReadNTPConfFile() returned error: %v", err)
	}

	expected := []types.NTPConfiguration{
		{Server: "0.us.pool.ntp.org", IBurst: true},
		{Server: "1.us.pool.ntp.org", IBurst: true},
		{Server: "2.us.pool.ntp.org", IBurst: false},
		{Server: "ntp.ubuntu.com", IBurst: true},
	}

	if len(configs) != len(expected) {
		t.Fatalf("got %d configs, want %d", len(configs), len(expected))
	}

	for i, cfg := range configs {
		if cfg.Server != expected[i].Server {
			t.Errorf("config[%d].Server = %q, want %q", i, cfg.Server, expected[i].Server)
		}
		if cfg.IBurst != expected[i].IBurst {
			t.Errorf("config[%d].IBurst = %v, want %v", i, cfg.IBurst, expected[i].IBurst)
		}
	}
}

func TestReadNTPConfFile_NoServers(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "ntp-conf-empty-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := `# Only comments here
# No server entries
driftfile /var/lib/ntp/ntp.drift
`
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	origPaths := ntpConfPaths
	ntpConfPaths = []string{tmpFile.Name()}
	defer func() { ntpConfPaths = origPaths }()

	_, err = ReadNTPConfFile()
	if err == nil {
		t.Error("ReadNTPConfFile() should return error when no servers found")
	}
}

func TestReadNTPConfFile_InlineComments(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "ntp-conf-inline-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := `server time.google.com iburst # Google NTP
server time.cloudflare.com # Cloudflare NTP
`
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	origPaths := ntpConfPaths
	ntpConfPaths = []string{tmpFile.Name()}
	defer func() { ntpConfPaths = origPaths }()

	configs, err := ReadNTPConfFile()
	if err != nil {
		t.Fatalf("ReadNTPConfFile() returned error: %v", err)
	}

	if len(configs) != 2 {
		t.Fatalf("got %d configs, want 2", len(configs))
	}
	if configs[0].Server != "time.google.com" || !configs[0].IBurst {
		t.Errorf("config[0] = %+v, want {Server: time.google.com, IBurst: true}", configs[0])
	}
	if configs[1].Server != "time.cloudflare.com" || configs[1].IBurst {
		t.Errorf("config[1] = %+v, want {Server: time.cloudflare.com, IBurst: false}", configs[1])
	}
}

func TestGetUserHomeDir(t *testing.T) {
	homeDir, err := GetUserHomeDir()
	if err != nil {
		t.Fatalf("GetUserHomeDir() returned error: %v", err)
	}
	if homeDir == "" {
		t.Error("GetUserHomeDir() returned empty string")
	}
	// Should be an absolute path
	if !filepath.IsAbs(homeDir) {
		t.Errorf("GetUserHomeDir() returned non-absolute path: %s", homeDir)
	}
}

func TestGetLargestFoldersConcurrent(t *testing.T) {
	// Create a temp directory with subdirectories
	tmpDir, err := os.MkdirTemp("", "genius-folders-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create subdirectories with files of known sizes
	dirs := []struct {
		name string
		size int
	}{
		{"large", 1000},
		{"medium", 500},
		{"small", 100},
	}

	for _, d := range dirs {
		dirPath := filepath.Join(tmpDir, d.name)
		if err := os.Mkdir(dirPath, 0755); err != nil {
			t.Fatalf("failed to create dir %s: %v", d.name, err)
		}
		filePath := filepath.Join(dirPath, "data.bin")
		data := make([]byte, d.size)
		if err := os.WriteFile(filePath, data, 0644); err != nil {
			t.Fatalf("failed to write file in %s: %v", d.name, err)
		}
	}

	t.Run("returns correct order", func(t *testing.T) {
		results, err := GetLargestFoldersConcurrent(tmpDir, 3)
		if err != nil {
			t.Fatalf("GetLargestFoldersConcurrent() returned error: %v", err)
		}

		if len(results) != 3 {
			t.Fatalf("got %d results, want 3", len(results))
		}

		// Should be sorted largest first
		if results[0].Size < results[1].Size || results[1].Size < results[2].Size {
			t.Errorf("results not sorted by size descending: %+v", results)
		}

		if filepath.Base(results[0].Path) != "large" {
			t.Errorf("first result should be 'large', got %s", filepath.Base(results[0].Path))
		}
	})

	t.Run("respects limit", func(t *testing.T) {
		results, err := GetLargestFoldersConcurrent(tmpDir, 2)
		if err != nil {
			t.Fatalf("GetLargestFoldersConcurrent() returned error: %v", err)
		}

		if len(results) != 2 {
			t.Errorf("got %d results, want 2", len(results))
		}
	})

	t.Run("handles empty directory", func(t *testing.T) {
		emptyDir, err := os.MkdirTemp("", "genius-empty-test-*")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(emptyDir)

		results, err := GetLargestFoldersConcurrent(emptyDir, 5)
		if err != nil {
			t.Fatalf("GetLargestFoldersConcurrent() returned error: %v", err)
		}

		if len(results) != 0 {
			t.Errorf("got %d results for empty dir, want 0", len(results))
		}
	})
}
