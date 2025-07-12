package server

import (
	"net"
	"strconv"
	"sync/atomic"
	"log"
	"github.com/trolioSFG/http-go/internal/response"
)


type Server struct {
	listener net.Listener
	closed atomic.Bool
}

func Serve(port int) (*Server, error) {
	// Alternative
	// .. fmt.Sprintf(":%d", port)
	l, err := net.Listen("tcp", ":" + strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	
	srv := &Server{
		listener: l,
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
	// msg := "Hello World!"
	response.WriteStatusLine(conn, response.StatusOK)
	response.WriteHeaders(conn, response.GetDefaultHeaders(0))
	conn.Write([]byte("\r\n"))

	conn.Close()
}

