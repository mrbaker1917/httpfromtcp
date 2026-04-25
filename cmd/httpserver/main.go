package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/mrbaker1917/httpfromtcp/internal/headers"
	"github.com/mrbaker1917/httpfromtcp/internal/request"
	"github.com/mrbaker1917/httpfromtcp/internal/response"
	"github.com/mrbaker1917/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		proxyHandler(w, req)
		return
	}
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		body := []byte("<html><head><title>400 Bad Request</title></head><body><h1>Bad Request</h1><p>Your request honestly kinda sucked.</p></body></html>")
		w.WriteStatusLine(response.BadRequest)
		h := response.GetDefaultHeaders(len(body))
		h["Content-Type"] = "text/html"
		w.WriteHeaders(h)
		w.WriteBody(body)
	case "/myproblem":
		body := []byte("<html><head><title>500 Internal Server Error</title></head><body><h1>Internal Server Error</h1><p>Okay, you know what? This one is on me.</p></body></html>")
		w.WriteStatusLine(response.ServerError)
		h := response.GetDefaultHeaders(len(body))
		h["Content-Type"] = "text/html"
		w.WriteHeaders(h)
		w.WriteBody(body)
	default:
		body := []byte("<html><head><title>200 OK</title></head><body><h1>Success!</h1><p>Your request was an absolute banger.</p></body></html>")
		w.WriteStatusLine(response.OK)
		h := response.GetDefaultHeaders(len(body))
		h["Content-Type"] = "text/html"
		w.WriteHeaders(h)
		w.WriteBody(body)
	}
}

func proxyHandler(w *response.Writer, req *request.Request) {
	s := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	reqLine := "https://httpbin.org/" + s
	resp, err := http.Get(reqLine)
	if err != nil {
		log.Printf("Error in getting response from httpbin: %v", err)
		return
	}
	defer resp.Body.Close()
	err = w.WriteStatusLine(response.OK)
	if err != nil {
		log.Printf("Error in writing status line: %v", err)
		return
	}
	h := response.GetDefaultHeaders(0)
	delete(h, "Content-Length")
	h["Transfer-Encoding"] = "chunked"
	h["Trailer"] = "X-Content-SHA256, X-Content-Length"
	err = w.WriteHeaders(h)
	if err != nil {
		log.Printf("Error in writing headers: %v", err)
		return
	}
	buf := make([]byte, 1024)
	var fullBody []byte
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			w.WriteChunkedBody(buf[:n])
			fullBody = append(fullBody, buf[:n]...)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error in writing reading response: %v", err)
			break
		}
	}
	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		log.Printf("Error in ending reading of chunk body: %v", err)
	}

	hsh := sha256.Sum256(fullBody)
	bodLength := len(fullBody)
	heads := headers.NewHeaders()
	heads["X-Content-SHA256"] = fmt.Sprintf("%x", hsh)
	heads["X-Content-Length"] = fmt.Sprintf("%d", bodLength)
	err = w.WriteTrailers(heads)
	if err != nil {
		log.Printf("Error in writing trailers: %v", err)
	}
}
