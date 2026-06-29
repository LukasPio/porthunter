package main

import (
	"fmt"
	"net"
	"net/netip"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const MaxIpToScanPerHostname = 3

func main() {

	if len(os.Args) < 4 {
		printUsageAndExample()
		os.Exit(1)
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
		fmt.Println("Ports can't be 0 or negative")
		os.Exit(1)
	}

	if start > end {
		fmt.Println("Start range must be less than or equal to end")
		os.Exit(1)
	}

	workersAmount := (end - start) / 50

	if workersAmount < 1 {
		workersAmount = 1
	}

	if workersAmount > 1000 {
		workersAmount = 1000
	}

	var ips [MaxIpToScanPerHostname]net.IP

	if net.ParseIP(target) != nil {
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

	if counter > 1 {
		fmt.Println(target + " resolved in multiple ips. " + "scanning " + strconv.Itoa(counter) + " ips...")
	}

	for _, ip := range ips {
		if ip != nil {
			scanSingleIp(ip.String(), start, end, workersAmount)
		}
	}
}

func printUsageAndExample() {
	fmt.Println("Usage: porthunter <ip-address> <range-start> <range-end>\nExample: porthunter 192.168.0.1 1 100 - ports across 1 and 100 on 192.168.0.1")
}

func scanSingleIp(ip string, start int, end int, workersAmount int) {
	var wg sync.WaitGroup
	ports := make(chan int, workersAmount)

	for i := 0; i < workersAmount; i++ {
		wg.Add(1)
		go scanPorts(ip, ports, &wg)
	}

	fmt.Print("Scanning for ports on " + ip + "\n")
	fmt.Println("=======================================")

	for i := start; i < end+1; i++ {
		ports <- i
	}

	close(ports)

	wg.Wait()
}

func resolveHostname(hostname string) [MaxIpToScanPerHostname]net.IP {
	ips, err := net.LookupIP(hostname)
	if err != nil {
		fmt.Println("Couldn't resolve hostname " + hostname + ":\n" + err.Error())
		os.Exit(2)
	}

	var ipv4ips [MaxIpToScanPerHostname]net.IP
	counter := 0

	for _, ip := range ips {
		addr, err := netip.ParseAddr(ip.String())
		if err != nil {
			fmt.Println("Internal error resolving hostname: " + err.Error())
			os.Exit(3)
		}
		if addr.Is4() {

			if counter > MaxIpToScanPerHostname {
				fmt.Println("Max ip is set to " + strconv.Itoa(MaxIpToScanPerHostname) + ". skiping some ips")
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
			tryGrabBanner(conn, portToString)
		}
	}
}

func tryGrabBanner(conn net.Conn, port string) {
	buffer := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, err := conn.Read(buffer)
	if err == nil {
		banner := string(buffer)
		service := identifyBanner(banner)
		if service == "" {
			printPortWithoutService(port)
		} else {
			printPortAndService(port, service, banner)
		}
		conn.Close()
		return
	}

	printPortWithoutService(port)
	conn.Close()
}

func identifyBanner(banner string) string {
	switch {
	case strings.Contains(banner, "SSH"):
		return "SSH"
	case strings.Contains(banner, "FTP"):
		return "FTP"
	case strings.Contains(banner, "SMTP"):
		return "SMTP"
	default:
		return ""
	}
}

func printPortWithoutService(port string) {
	fmt.Println("Porta " + port + " aberta. Não foi possível identificar o serviço")
	fmt.Println("=======================================")
}

func printPortAndService(port string, service string, banner string) {
	fmt.Println("Porta: " + port)
	fmt.Println("Serviço: " + service)
	fmt.Print("Banner: " + banner)
	fmt.Println("=======================================")
}
