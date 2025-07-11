package main

import (
	"net"
	"os"
	"fmt"
	"bufio"
	"log"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Printf("Error resolving: %v\n", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Printf("Error dialing: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	
	for {
		fmt.Print("> ")
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading: %v\n", err)
			continue
		}

		n, err := conn.Write([]byte(msg))
		if err != nil {
			log.Printf("Error writing: %v\n", err)
			continue
		}
		log.Printf("%d bytes written\n", n)
	}
}



