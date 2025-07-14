package main

import (
//	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"github.com/trolioSFG/http-go/internal/server"
	"github.com/trolioSFG/http-go/internal/request"
	"github.com/trolioSFG/http-go/internal/response"
	"github.com/trolioSFG/http-go/internal/headers"
)

const port = 42069

func fooHandler(w *response.Writer, req *request.Request) {

	hdr := headers.NewHeaders()

	if req == nil {
		err := w.WriteStatusLine(response.StatusBadRequest)
		if err != nil {
			return
		}

		str := `<html><head><title>400 Bad Request</title></head>
<body><h1>Bad Request</h1><p>Your request honestly kinda sucked.</p></body></html>`
		msg := []byte(str)

		hdr["Content-Type"] = "text/html"
		hdr["Content-Length"] = strconv.Itoa(len(msg))

		err = w.WriteHeaders(hdr)
		if err != nil {
			return
		}
		_, err = w.WriteBody(msg)
		if err != nil {
			return
		}

		return
	}



	if req.RequestLine.RequestTarget == "/yourproblem" {
		w.WriteStatusLine(response.StatusBadRequest)
		msg := []byte(`<html><head><title>400 BadRequest</title></head>
<body><h1>Bad Request</h1><p>Your request honestly kinda sucked.</p></body></html>`)
		hdr["Content-Type"] = "text/html"
		hdr["Content-Length"] = strconv.Itoa(len(msg))
		w.WriteHeaders(hdr)
		w.WriteBody(msg)
		return
	}

	if req.RequestLine.RequestTarget == "/myproblem" {
		w.WriteStatusLine(response.StatusError)
		msg := []byte(`<html><head><title>500 Internal Server Error</title></head>
<body><h1>Internal Server Error</h1><p>Okay, you know what? This one is on me.</p></body></html>`)
		hdr["Content-Type"] = "text/html"
		hdr["Content-Length"] = strconv.Itoa(len(msg))
		w.WriteHeaders(hdr)
		w.WriteBody(msg)
		return
		/**
		hnd.Status = response.StatusError
		hnd.Msg = []byte("r\n")
		return hnd
		**/
	}

	w.WriteStatusLine(response.StatusOK)
	msg := []byte(`<html><head><title>200 OK</title></head>
<body><h1>Success!</h1><p>Your request was an absolute banger.</p></body></html>`)
	hdr["Content-Type"] = "text/html"
	hdr["Content-Length"] = strconv.Itoa(len(msg))
	w.WriteHeaders(hdr)
	w.WriteBody(msg)

	return
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

