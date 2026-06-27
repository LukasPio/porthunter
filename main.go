package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

var targetIp string

func main() {

	if len(os.Args) < 4 {
		println("Missing parameters\nUsage: porthunter <ip-address> <range-start> <range-end>\nExample: porthunter 192.168.0.1 1 100 - ports across 1 and 100 on 192.168.0.1")
		os.Exit(1)
	}

	targetIp = os.Args[1]

	start, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Invalid start range")
		os.Exit(1)
	}
	
	end, err := strconv.Atoi(os.Args[3])
	if err != nil || end > 65535 {
		fmt.Println("Invalid end range")
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

	var wg sync.WaitGroup
	ports := make(chan int, workersAmount)

	for i := 0; i < workersAmount; i++ {
		wg.Add(1)
		go scanPorts(ports, &wg)
	}

	fmt.Println("Scanning for ports on " + targetIp)

	for i := start; i < end+1; i++ {
		ports <- i
	}

	close(ports)

	wg.Wait()
}

func scanPorts(ports <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for port := range ports {

		portToString := strconv.Itoa(port)

		connection, error := net.DialTimeout("tcp", targetIp+":"+portToString, 2*time.Second)

		if error == nil {
			connection.Close()
			fmt.Println("Porta " + portToString + " aberta!")
		}
	}
}