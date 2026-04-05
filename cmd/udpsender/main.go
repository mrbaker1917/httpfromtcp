package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	udpAdr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Printf("error resolving UDP address: %s\n", err)
	}
	conn, err := net.DialUDP("udp", nil, udpAdr)
	if err != nil {
		log.Printf("error dialing UDP: %s\n", err)
	}
	defer conn.Close()

	r := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := r.ReadString('\n')
		if err != nil {
			log.Fatalf("error in reading string: %s\n", err)
		}
		conn.Write([]byte(line))
	}
}
