package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/netip"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const MaxIpToScanPerHostname = 3

var (
	scanService = false
	verbose = false
	saveOutput = false
	outputBuffer bytes.Buffer
	writer = io.MultiWriter(os.Stdout, &outputBuffer)
)

func main() {
	target, start, end := parseArgs()
	workersAmount := (end - start) / 50

	if workersAmount < 1 {
		workersAmount = 1
	}

	if workersAmount > 1000 {
		workersAmount = 1000
	}

	var ips [MaxIpToScanPerHostname]net.IP

	if net.ParseIP(target) != nil {
		if verbose {
			fmt.Fprintln(writer, "One target mode actived...")	
		}
		scanSingleIp(target, start, end, workersAmount)
		return
	}

	ips = resolveHostname(target)

	counter := 0
	for _, ip := range ips {
		if ip != nil {
			counter++
		}
	}

	if counter > 1 && verbose {
		fmt.Fprintln(writer,target + " resolved in multiple ips. " + "scanning " + strconv.Itoa(counter) + " ips...")
	}

	for _, ip := range ips {
		if ip != nil {
			scanSingleIp(ip.String(), start, end, workersAmount)
		}
	}

	if saveOutput {
		saveOutputToFile()	
	}
}


func parseArgs() (string, int, int) {
	if len(os.Args) < 4 {
		printUsageAndExample()
		os.Exit(1)
	}

	hasOptions := true

	if len(os.Args) == 4 {
		hasOptions = false
	}

	target := os.Args[1]

	start, err := strconv.Atoi(os.Args[2])
	if err != nil {
		printUsageAndExample()
		os.Exit(1)
	}

	end, err := strconv.Atoi(os.Args[3])
	if err != nil || end > 65535 {
		printUsageAndExample()
		os.Exit(1)
	}

	if start <= 0 || end <= 0 {
		fmt.Fprintln(writer,"Ports can't be 0 or negative")
		os.Exit(1)
	}

	if start > end {
		fmt.Fprintln(writer,"Start range must be less than or equal to end")
		os.Exit(1)
	}

	if hasOptions {
		flags := os.Args[4:]
		counter := 0

		for _, flag := range flags {

			switch {
			case flag == "-v":
				verbose = true
				fmt.Fprintln(writer,"Program starting at " + time.Now().Format("2006-01-02T15:04:05.000Z07:00"))
				fmt.Fprintln(writer,"Verbose mode activated...")
			case flag == "-sV":
				scanService = true
			case flag == "-o":
				saveOutput = true
			default:
				fmt.Fprintln(writer,"Unrecognized option '" + flag + "'")
				os.Exit(4)
			}

			counter++

			if counter > 2 {
				fmt.Fprintln(writer,"Max of 3 options are supported. Skipping some extra flags...")
				break
			}
		}
	}

	return target, start, end
}


func printUsageAndExample() {
	fmt.Fprintln(writer,"Usage: porthunter <ip-address> <range-start> <range-end>\nExample: porthunter 192.168.0.1 1 100 - ports across 1 and 100 on 192.168.0.1")
}

func scanSingleIp(ip string, start int, end int, workersAmount int) {
	var wg sync.WaitGroup
	ports := make(chan int, workersAmount)

	for i := 0; i < workersAmount; i++ {
		wg.Add(1)
		go scanPorts(ip, ports, &wg)
	}

	fmt.Fprint(writer, "Scanning for ports on " + ip + "\n")
	fmt.Fprintln(writer,"=======================================")

	for i := start; i < end+1; i++ {
		ports <- i
	}

	close(ports)

	wg.Wait()
}

func resolveHostname(hostname string) [MaxIpToScanPerHostname]net.IP {
	ips, err := net.LookupIP(hostname)
	if err != nil {
		fmt.Fprintln(writer,"Couldn't resolve hostname " + hostname + ":\n" + err.Error())
		os.Exit(2)
	}

	var ipv4ips [MaxIpToScanPerHostname]net.IP
	counter := 0

	for _, ip := range ips {
		addr, err := netip.ParseAddr(ip.String())
		if err != nil {
			fmt.Fprintln(writer,"Internal error resolving hostname: " + err.Error())
			os.Exit(3)
		}
		if addr.Is4() {
			
			if verbose {
				fmt.Fprintln(writer,"Resolved the hostname " + hostname + " to " + addr.String())
			}			

			if counter >= MaxIpToScanPerHostname {
				if verbose {
					fmt.Fprintln(writer,"Max ip is set to " + strconv.Itoa(MaxIpToScanPerHostname) + ". skiping some ips")
				}
				break
			}

			ipv4ips[counter] = ip
			counter++
		}
	}

	return ipv4ips
}

func scanPorts(target string, ports <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for port := range ports {
		portToString := strconv.Itoa(port)

		conn, err := net.DialTimeout("tcp", target+":"+portToString, 3*time.Second)

		if err == nil {
			if scanService {
				tryIdentifyService(conn, portToString)
			} else {
				printOpenPort(portToString)
			}
		}
	}
}

func tryIdentifyService(conn net.Conn, port string) {
	buffer := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, err := conn.Read(buffer)
	if err == nil {
		banner := string(buffer)
		service := identifyBanner(banner)
		if service == "" {
			printPortAndService(port, service, banner, true)
		} else {
			printPortAndService(port, service, banner, false)
		}
		conn.Close()
		return
	}

	printOpenPort(port)
	conn.Close()
}

func identifyBanner(banner string) string {
	upper := strings.ToUpper(banner)
	switch { 
	case strings.Contains(upper, "SSH"):
		return "SSH"
	case strings.Contains(upper, "FTP"):
		return "FTP"
	case strings.Contains(upper, "SMTP"):
		return "SMTP"
	default:
		return ""
	}
}

func printOpenPort(port string) {
	fmt.Fprintln(writer,"Port " + port + " is open")
}

func printPortAndService(port string, service string, banner string, failure bool) {
	if failure {
		fmt.Fprintln(writer,"Port: " + port)
		fmt.Fprintln(writer,"Was not possible to identify banner - " + banner)
		fmt.Fprintln(writer,"=======================================")
		return
	}

	fmt.Fprintln(writer,"Port: " + port)
	fmt.Fprintln(writer,"Service: " + service)
	fmt.Fprint(writer, "Banner: " + banner)
	fmt.Fprintln(writer,"=======================================")
}

func saveOutputToFile() {
	output := outputBuffer.Bytes()
	outputFile, err := os.Create("output.txt")
	
	if err != nil {
		fmt.Fprintln(writer,"Occurred an error trying to save output")
		return
	}
	outputFile.Write(output)
}