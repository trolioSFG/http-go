package main

import (
	"fmt"
//	"os"
	"io"
	"strings"
	"net"
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
		c := getLinesChannel(conn)

		for line := range c {
			fmt.Printf("%s\n", line)
		}

		fmt.Println("Connection closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	c := make(chan string)

	go func() {

		buf := make([]byte, 8)
		var currentLine string

		for {
			n, err := f.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Printf("Error reading: %v\n", err)
				break
			}
			parts := strings.Split(string(buf[:n]), "\n")
			for i := 0; i < len(parts) - 1 && len(parts) > 1; i++ {
				currentLine = currentLine + parts[i]
				// fmt.Printf("read: %s\n", currentLine)
				c <- currentLine
				currentLine = ""
			}

			currentLine = currentLine + parts[len(parts)-1]

			// fmt.Printf("read: %s\n", buf[:n])
		}

		if currentLine != "" {
			c <- currentLine
		}

		f.Close()
		close(c)
	}()

	return c
}

