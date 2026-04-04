package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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

		ch := getLinesChannel(conn)
		for line := range ch {
			fmt.Println(line)
		}
		fmt.Printf("Connection to %s closed.\n", conn.RemoteAddr())
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {

	ch := make(chan string)
	go func() {
		defer f.Close()
		defer close(ch)
		currentLineContents := ""
		for {
			buffer := make([]byte, 8)
			n, err := f.Read(buffer)
			if err != nil {
				if currentLineContents != "" {
					ch <- currentLineContents
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				return
			}
			str := string(buffer[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				ch <- currentLineContents + parts[i]
				currentLineContents = ""

			}
			currentLineContents += parts[len(parts)-1]
		}
	}()
	return ch
}
