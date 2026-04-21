package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/mrbaker1917/httpfromtcp/internal/request"
	"github.com/mrbaker1917/httpfromtcp/internal/response"
)

type Server struct {
	closed   atomic.Bool
	listener net.Listener
	handler  Handler
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w *response.Writer, req *request.Request)

func Serve(port int, handler Handler) (*Server, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("Error in listening to TCP traffic: %s\n", err.Error())
	}
	s := &Server{
		listener: ln,
		handler:  handler,
	}
	go s.listen()

	return s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	err := s.listener.Close()
	if err != nil {
		return fmt.Errorf("Unable to close server: %s", err)
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
	w := response.NewWriter(conn)
	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Printf("Error in reading request from reader: %v", err)
		err = w.WriteStatusLine(response.BadRequest)
		if err != nil {
			log.Printf("Error writing status line: %v", err)
		}
		h := response.GetDefaultHeaders(23)
		err = w.WriteHeaders(h)
		if err != nil {
			log.Printf("Error writing headers: %v", err)
		}
		_, err = w.WriteBody([]byte("Could not parse request"))
		if err != nil {
			log.Printf("Error writing body: %v", err)
		}
		return
	}
	s.handler(w, req)
}
