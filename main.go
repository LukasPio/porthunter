package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {

	targetIp := os.Args[1]

	ports := [7]int{21, 22, 80, 3000, 3306, 8080, 443}

	for i := 0; i < len(ports); i++ {

		currentPort := strconv.Itoa(ports[i])

		connection, error := net.DialTimeout("tcp", targetIp + ":" + currentPort, 2*time.Second)

		if error == nil {
			connection.Close()
			fmt.Println("Porta " + currentPort + " aberta!")
		}
	}

}