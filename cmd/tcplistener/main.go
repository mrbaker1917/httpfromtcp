package main

import (
	"fmt"
	"log"
	"net"

	"github.com/mrbaker1917/httpfromtcp/internal/request"
)

const port = ":42069"

func main() {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Error in listening to TCP traffic: %s\n", err.Error())
	}
	defer ln.Close()

	fmt.Println("Listening for TCP traffic on", port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error in accepting connection: %v\n", err)
			break
		}
		fmt.Printf("Accepted connection from: %s\n", conn.RemoteAddr())

		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("Error occurred in requesting data: %v\n", err)
		}
		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for k, v := range req.Headers {
			fmt.Printf("- %s: %s\n", k, v)
		}
	}
}
