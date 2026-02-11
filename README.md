# Genius

Genius is a simple and easy-to-use CLI tool for gathering system information.
It supports **Linux** and **macOS** systems.

## Features

- **System Information** — hostname, OS, architecture, kernel version, platform, virtualization, uptime
- **CPU** — model, core count, load average, per-core usage and times
- **Memory** — total RAM, usage percentage, swap
- **Disk** — total/used/free space, filesystem type
- **Network** — all active non-loopback interfaces with IPv4, MAC, MTU
- **DNS** — nameservers from `/etc/resolv.conf`
- **NTP** — local NTP configuration, time comparison against a configurable NTP server
- **Software** — Python version (and Homebrew version on macOS)
- **Largest Folders** — top N largest folders in home directory (concurrent scanning)

## Download

You can download the latest version from the [releases page](https://github.com/kerimkaan/genius/releases).

## Installation

### From source

```bash
git clone https://github.com/kerimkaan/genius.git
cd genius
go mod download
CGO_ENABLED=0 go build -ldflags="-s -w" -o genius main.go
sudo mv genius /usr/local/bin/
```

### With custom version

```bash
CGO_ENABLED=0 go build -ldflags="-s -w -X genius/constants.Version=1.0.0" -o genius main.go
```

## Usage

### System Information

```bash
genius info                          # Full system report
genius info --ntp-server pool.ntp.org  # Use custom NTP server for time comparison
genius i                             # Short alias
```

### CPU Details

```bash
genius cpu    # Detailed CPU information
genius c      # Short alias
```

### Largest Folders

```bash
genius largest-folders          # Top 5 largest folders in home directory
genius largest-folders -n 10    # Top 10 largest folders
genius lf                       # Short alias
```

### Help

```bash
genius help
genius info --help
```

## Supported Platforms

| Platform      | Architecture |
|---------------|-------------|
| Linux         | amd64       |
| Linux         | arm64       |
| macOS (Darwin)| amd64       |
| macOS (Darwin)| arm64       |

## Development

```bash
# Build
go build -o genius main.go

# Run tests
go test -v -race ./...

# Vet
go vet ./...
```

## License

See [LICENSE](LICENSE) for details.
