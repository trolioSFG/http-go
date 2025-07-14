package server

import (
//	"bytes"
//	"io"
//	"bufio"
	"fmt"
	"net"
	"strconv"
	"sync/atomic"
	"log"
	"github.com/trolioSFG/http-go/internal/response"
	"github.com/trolioSFG/http-go/internal/request"
)


type Server struct {
	listener net.Listener
	closed atomic.Bool
	handler Handler
}

type Handler func(w *response.Writer, req *request.Request)

/**
type HandlerError struct {
	Status response.StatusCode
	Msg []byte
}

func WriteError(w io.Writer, herr *HandlerError) error {
	err := response.WriteStatusLine(w, herr.Status)
	if err != nil {
		return err
	}
	err = response.WriteHeaders(w, response.GetDefaultHeaders(len(herr.Msg)))
	if err != nil {
		return err
	}

	_, err = w.Write([]byte("\r\n"))
	if err != nil {
		return err
	}
	_, err = w.Write(herr.Msg)
	if err != nil {
		return err
	}

	return nil
}
**/


func Serve(port int, handler Handler) (*Server, error) {
	// Alternative
	// .. fmt.Sprintf(":%d", port)
	l, err := net.Listen("tcp", ":" + strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	
	srv := &Server{
		listener: l,
		handler: handler,
	}

	go srv.listen()

	return srv, nil
}

func (s *Server) Close() error {
	// Atomic
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}


func (s *Server) handle(conn net.Conn) {

	defer conn.Close()

	// Buf: bytes.NewBuffer([]byte{}),
	w := &response.Writer{
		Buf: conn,
		State: response.StateStatus,
	}
	// Parse the request from the connection
	req, err := request.RequestFromReader(conn)
	if err != nil {
		w.WriteStatusLine(response.StatusBadRequest)
		body := []byte(fmt.Sprintf("Error parsing request: %v", err))
		w.WriteHeaders(response.GetDefaultHeaders(len(body)))
		w.WriteBody(body)
		return
	}		


	s.handler(w, req)
}


