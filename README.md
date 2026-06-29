# PortHunter
PortHunter is a simple multithreaded TCP port scanner written in Go. It scans a range of ports on a target host using goroutines and channels, with optional service detection via banner grabbing.

> This project was built as a learning exercise to better understand Go concurrency, networking, and how basic port scanners such as Nmap work internally.

## Features
- Scan any TCP port range
- Concurrent scanning using goroutines
- Automatic worker scaling
- Connection timeout to avoid hanging
- Hostname resolution with support for up to 3 IPv4 addresses
- Optional service/banner identification (`-sV`)
- Verbose mode (`-v`)
- Save output to file (`-o`)

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
| `-v`  | Verbose mode — shows timestamps, resolved IPs, and diagnostic info |
| `-sV` | Service scan — attempts to grab banners and identify services |
| `-o`  | Save output to `output.txt` |

Up to 3 options can be combined in any order.

### Examples
```bash
porthunter scanme.nmap.org 1 1000
porthunter 192.168.1.1 20 1024 -v -sV -o
```

## Example Output

**Standard scan:**
```text
Scanning for ports on 45.33.32.156
=======================================
Port 22 is open
Port 80 is open
Port 9929 is open
```

**With `-sV` (service scan):**
```text
Scanning for ports on 45.33.32.156
=======================================
Port: 22
Service: SSH
Banner: SSH-2.0-OpenSSH_6.6.1p1 Ubuntu-2ubuntu2.13
=======================================
Port 80 is open. Could not identify service.
```

**With `-v` (verbose):**
```text
Program starting at 2024-11-03T21:04:05.000-03:00
Verbose mode activated...
Resolved the hostname scanme.nmap.org to 45.33.32.156
Scanning for ports on 45.33.32.156
=======================================
Port 22 is open
```

## How It Works
1. Validates command-line arguments and parses any option flags.
2. If a hostname is given, resolves it to up to 3 IPv4 addresses.
3. Calculates an appropriate number of worker goroutines based on the scan range.
4. Creates a channel containing every port in the specified range.
5. Worker goroutines receive ports from the channel and attempt TCP connections using `net.DialTimeout`.
6. If `-sV` is set, attempts to read a banner and identify the service (SSH, FTP, SMTP).
7. If `-o` is set, all output is also written to `output.txt`.
8. The program waits for every worker to finish before exiting.

## Worker Scaling
The number of workers is automatically calculated based on the scan range size.
- Minimum: **1 worker**
- Maximum: **1000 workers**

This keeps small scans lightweight while allowing large scans to complete much faster.

## Limitations
- TCP connect scan only
- Service identification limited to SSH, FTP, and SMTP
- Does not perform OS fingerprinting
- Does not support UDP scanning
- No structured output formats (JSON/XML)

## Learning Goals
This project explores several Go concepts, including:
- Goroutines and channels
- WaitGroups
- TCP networking with timeouts
- CLI argument parsing
- Concurrent programming patterns
- I/O multiplexing with `io.MultiWriter`

## Disclaimer
Use this software only on systems that you own or have explicit permission to scan. Unauthorized port scanning may violate network policies or local laws.
