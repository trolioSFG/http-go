package server

import (
	"bytes"
	"io"
//	"bufio"
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

type Handler func(w io.Writer, req *request.Request) *HandlerError

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

	// Parse the request from the connection
	req, err := request.RequestFromReader(conn)
	if err != nil {
		WriteError(conn, &HandlerError{
			Status: response.StatusBadRequest,
			Msg: []byte(err.Error()),
		})
		return
	}

	// log.Printf(req.RequestLine.String())

	w := bytes.NewBuffer([]byte{})
	hndErr := s.handler(w, req)

	if hndErr.Status != response.StatusOK {
		WriteError(conn, hndErr)
		return
	}

	response.WriteStatusLine(conn, response.StatusOK)
	response.WriteHeaders(conn, response.GetDefaultHeaders(w.Len()))
	conn.Write([]byte("\r\n"))
	conn.Write(w.Bytes())
}


