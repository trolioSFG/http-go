package main

import (
	"fmt"
//	"io"
	"net/http"
	"log"
	"os"
	"os/signal"
	"strings"
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

	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		ChunkedHandler(w, req)
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


func ChunkedHandler(w *response.Writer, r *request.Request) {
	url := "https://httpbin.org/" + strings.TrimPrefix(r.RequestLine.RequestTarget, "/httpbin/")
	rsp, err := http.Get(url)
	if err != nil || rsp.StatusCode > 200 {
		httpbinErr := err
		err := w.WriteStatusLine(response.StatusBadRequest)
		if err != nil {
			return
		}

		str := fmt.Sprintf(`<html><head><title>400 Bad Request</title></head>
		<body><h1>Bad Request</h1><p>Your request honestly kinda sucked.</p><p>HTTPBIN err: %v</p>
		<p>HTTPBIN StatusCode: %d</p></body></html>`, httpbinErr, rsp.StatusCode)
		msg := []byte(str)

		hdr := headers.NewHeaders()
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

	defer rsp.Body.Close()
	w.WriteStatusLine(response.StatusOK)
	hdr := headers.NewHeaders()
	hdr["Content-Type"] = rsp.Header["Content-Type"][0]
	hdr["Transfer-Encoding"] ="chunked"
	w.WriteHeaders(hdr)

	buf := make([]byte, 20)
	finished := false
	bytesRead := 0
	for !finished {
		bytesRead, err = rsp.Body.Read(buf)
		if bytesRead > 0 {
			w.WriteChunkedBody(buf[:bytesRead])
		} else {
			finished = true
		}
	}

	w.WriteChunkedBodyDone()
	

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

