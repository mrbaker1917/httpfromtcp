package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
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

type Handler func(w io.Writer, req *request.Request) *HandlerError

func (he *HandlerError) HandleError(w io.Writer) error {
	err := response.WriteStatusLine(w, he.StatusCode)
	if err != nil {
		return err
	}
	h := response.GetDefaultHeaders(len(he.Message))
	err = response.WriteHeaders(w, h)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(he.Message))
	return err
}

func Serve(port int, handler Handler) (*Server, error) {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))

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
	for !s.closed.Load() {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error getting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		hErr := &HandlerError{
			StatusCode: response.BadRequest,
			Message:    "Bad request\n",
		}
		hErr.HandleError(conn)
		return
	}
	b := bytes.Buffer{}
	hErr := s.handler(&b, req)
	if hErr != nil {
		hErr.HandleError(conn)
		return
	}
	h := response.GetDefaultHeaders(b.Len())
	err = response.WriteStatusLine(conn, response.OK)
	if err != nil {
		fmt.Printf("Error in writing the status line: %v", err)
	}
	err = response.WriteHeaders(conn, h)
	if err != nil {
		fmt.Printf("Error in writing headers: %v", err)
	}
	_, err = conn.Write(b.Bytes())
	if err != nil {
		fmt.Printf("Error in writing response body: %v", err)
	}
}
