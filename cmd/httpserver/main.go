package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"github.com/trolioSFG/http-go/internal/server"
	"github.com/trolioSFG/http-go/internal/request"
	"github.com/trolioSFG/http-go/internal/response"
)

const port = 42069

func fooHandler(w io.Writer, req *request.Request) *server.HandlerError {
	hnd := &server.HandlerError{}
	
	if req.RequestLine.RequestTarget == "/yourproblem" {
		hnd.Status = response.StatusBadRequest
		hnd.Msg = []byte("Your problem is not my problem\r\n")
		return hnd
	}

	if req.RequestLine.RequestTarget == "/myproblem" {
		hnd.Status = response.StatusError
		hnd.Msg = []byte("Woopsie, my bad\r\n")
		return hnd
	}

	hnd.Status = response.StatusOK
	w.Write([]byte("All good, frfr\r\n"))
	return hnd
}

func main() {
	server, err := server.Serve(port, fooHandler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefull stopped")
}

