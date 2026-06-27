# PortHunter

PortHunter is a simple multithreaded TCP port scanner written in Go. It scans a range of ports on a target host using goroutines and channels to perform concurrent connections.

> This project was built as a learning exercise to better understand Go concurrency, networking, and how basic port scanners such as Nmap work internally.

## Features

- Scan any TCP port range
- Concurrent scanning using goroutines
- Automatic worker scaling
- Connection timeout to avoid hanging
- Simple command-line interface

## Requirements

- Go 1.20 or newer

## Building

```bash
go build -o porthunter
```

## Usage

```bash
porthunter <target> <start-port> <end-port>
```

### Example

```bash
porthunter scanme.nmap.org 1 1000
```

or

```bash
porthunter 192.168.1.1 20 1024
```

## Example Output

```text
Scanning for ports on scanme.nmap.org

Porta 22 aberta!
Porta 80 aberta!
Porta 9929 aberta!
```

## How It Works

1. Validates command-line arguments.
2. Calculates an appropriate number of worker goroutines.
3. Creates a channel containing every port in the specified range.
4. Worker goroutines receive ports from the channel.
5. Each worker attempts a TCP connection using `net.DialTimeout`.
6. If the connection succeeds, the port is reported as open.
7. The program waits until every worker finishes.

## Worker Scaling

The number of workers is automatically calculated based on the scan size.

- Minimum: **1 worker**
- Maximum: **1000 workers**

This keeps small scans lightweight while allowing large scans to complete much faster.

## Limitations

- TCP connect scan only
- Does not detect service versions
- Does not perform OS fingerprinting
- Does not support UDP scanning
- No output formats (JSON/XML)
- No banner grabbing

## Learning Goals

This project explores several Go concepts, including:

- Goroutines
- Channels
- WaitGroups
- TCP networking
- Timeouts
- Concurrent programming

## Disclaimer

Use this software only on systems that you own or have explicit permission to scan. Unauthorized port scanning may violate network policies or local laws.