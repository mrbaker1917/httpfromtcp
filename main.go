package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	ln, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("Error in setting up TCP listener: %v", err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error in accepting connection: %v\n", err)
			break
		}
		fmt.Printf("Connection has been accepted: %s\n", conn.RemoteAddr())

		ch := getLinesChannel(conn)
		for line := range ch {
			fmt.Printf("%s\n", line)
		}
		fmt.Println("Connection closed.")
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
