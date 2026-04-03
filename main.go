package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const inputFilePath = "messages.txt"

func getLinesChannel(f io.ReadCloser) <-chan string {

	ch := make(chan string)
	currentLineContents := ""
	go func() {
		defer f.Close()
		defer close(ch)
		for {
			buffer := make([]byte, 8)
			n, err := f.Read(buffer)
			if err != nil {
				if currentLineContents != "" {
					ch <- currentLineContents
					currentLineContents = ""
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

func main() {

	f, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("Unable to open file %s: %v\n", inputFilePath, err)
	}

	ch := getLinesChannel(f)
	for line := range ch {
		fmt.Printf("read: %s\n", line)
	}

}
