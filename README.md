# PortHunter

PortHunter is a simple concurrent TCP port scanner written in Go. It scans a range of TCP ports on a target host using goroutines and channels, with optional banner grabbing for basic service identification.

> This project was built as a learning exercise to better understand Go concurrency, networking, and how basic port scanners such as Nmap work internally. It is now considered feature complete and is no longer under active development.

## Features

- Scan any TCP port range
- Concurrent scanning using goroutines and channels
- Automatic worker scaling
- Connection timeout to avoid hanging on filtered hosts
- Hostname resolution (up to 3 IPv4 addresses)
- Basic banner grabbing (`-sV`)
- Verbose mode (`-v`)
- Save scan output to `output.txt` (`-o`)

## Requirements

- Go 1.20 or newer

## Building

```bash
go build -o porthunter
```

## Usage

```bash
porthunter <target> <start-port> <end-port> [options]
```

### Options

| Flag | Description |
|------|-------------|
| `-v` | Enables verbose output |
| `-sV` | Attempts to identify services through banner grabbing |
| `-o` | Saves all program output to `output.txt` |

Up to three options can be combined.

## Examples

```bash
porthunter scanme.nmap.org 1 1000

porthunter 192.168.1.10 20 1024 -v

porthunter scanme.nmap.org 1 1000 -v -sV -o
```

## Example Output

### Standard scan

```text
Scanning for ports on 45.33.32.156
=======================================
Port 22 is open
Port 80 is open
Port 9929 is open
```

### Banner grabbing (`-sV`)

```text
Scanning for ports on 45.33.32.156
=======================================
Port: 22
Service: SSH
Banner: SSH-2.0-OpenSSH_6.6.1p1 Ubuntu-2ubuntu2.13
=======================================
```

### Verbose mode

```text
Program starting at 2026-07-21T19:30:00.000Z
Verbose mode activated...
Resolved the hostname scanme.nmap.org to 45.33.32.156
Scanning for ports on 45.33.32.156
=======================================
Port 22 is open
```

## How It Works

1. Parses and validates command-line arguments.
2. Resolves hostnames to IPv4 addresses when necessary.
3. Calculates an appropriate number of worker goroutines.
4. Sends every port in the requested range through a channel.
5. Workers perform TCP connect scans using `net.DialTimeout`.
6. If `-sV` is enabled, attempts to read and identify service banners.
7. If `-o` is enabled, output is also written to `output.txt`.
8. Waits for all workers to complete before exiting.

## Worker Scaling

Workers are automatically chosen based on the scan range.

- Minimum: **1 worker**
- Maximum: **1000 workers**

## Limitations

- TCP connect scan only
- IPv4 only
- Banner detection limited to SSH, FTP and SMTP
- Banner grabbing depends on the service sending data immediately after connection
- No UDP scanning
- No OS fingerprinting
- No CIDR/network scanning
- No structured output formats (JSON/XML)

## Technologies Used

- Go
- Goroutines
- Channels
- WaitGroups
- `net.DialTimeout`
- `io.MultiWriter`

## Learning Goals

This project was primarily created to practice:

- Go concurrency
- Goroutines and channels
- WaitGroups
- TCP networking
- CLI application development
- Concurrent worker pool patterns
- Basic service detection through banner grabbing

## Disclaimer

Use this software only on systems that you own or have explicit permission to scan. Unauthorized port scanning may violate network policies or local laws.
