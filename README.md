# Genius

Genius is a simple and easy to use CLI tool for gathering information about a
system. It is valid for Linux and MacOS systems.

## Installation

```bash
git clone https://github.com/kerimkaan/genius.git
cd genius
go mod download
go build -ldflags="-s -w" -o /usr/local/bin/genius main.go
```

## Usage

```bash
genius help
```
