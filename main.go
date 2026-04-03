package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

const inputFilePath = "messages.txt"

func main() {

	f, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("Unable to open file %s: %v\n", inputFilePath, err)
	}
	defer f.Close()

	fmt.Printf("Reading data from %s\n", inputFilePath)
	fmt.Println("=====================================")

	for {
		b := make([]byte, 8)

		n, err := f.Read(b)
		if n > 0 {
			fmt.Printf("read: %s\n", string(b[:n]))
		}
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}
			fmt.Printf("Error reading chunk: %v", err)
		}
	}
}
