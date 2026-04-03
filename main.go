package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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

	currentLine := ""

	for {

		b := make([]byte, 8)
		n, err := f.Read(b)

		str := string(b[:n])
		parts := strings.Split(str, "\n")

		for i := 0; i < len(parts)-1; i++ {
			fmt.Printf("read: %s%s\n", currentLine, parts[i])
			currentLine = ""
		}
		currentLine += parts[len(parts)-1]

		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}
			fmt.Printf("Error reading chunk: %v\n", err)
		}
	}
	fmt.Printf("read: %s\n", currentLine)
}
