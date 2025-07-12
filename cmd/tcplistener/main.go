package main

import (
	"fmt"
//	"os"
//	"io"
//	"strings"
	"net"
	"github.com/trolioSFG/http-go/internal/request"
)

func main() {
	// fmt.Println("I hope I get the job!")
	// file, err := os.Open("messages.txt")
	lstn, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Printf("Error opening port 42069: %v", err)
	}
	defer lstn.Close()

	for {
		conn, err := lstn.Accept()
		if err != nil {
			fmt.Printf("Accept error: %v\n", err)
			continue
		}

		fmt.Println("Connection accepted")
		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("Request error: %v\n", err)
			continue
		}

		fmt.Printf("Request line:\n")
		fmt.Printf("- Method: %v\n", req.RequestLine.Method)
		fmt.Printf("- Target: %v\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %v\n", req.RequestLine.HttpVersion)
		fmt.Printf("Headers:\n")
		for k, v := range req.Headers {
			fmt.Printf("- %v: %v\n", k, v)
		}

		fmt.Println("Connection closed")
	}
}

