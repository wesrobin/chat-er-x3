package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	msg := "hello there"
	_, err = conn.Write([]byte(msg))
	if err != nil {
		log.Fatal(err)
	}

	rspBuf := make([]byte, 1024)
	n, err := conn.Read(rspBuf)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Rsp", string(rspBuf[:n]))

	for {
		scanner := bufio.NewScanner(os.Stdin)

		var inp string
		for scanner.Scan() {
			tkn := scanner.Text()
			if tkn == "" {
				break
			}
			inp += tkn
		}

		if inp == "quit" {
			return
		}

		_, err := conn.Write([]byte(inp))
		if err != nil {
			log.Fatal(err)
		}

		n, err := conn.Read(rspBuf)
		if err != nil {
			log.Fatal(err)
		}

		log.Print("Rsp: ", string(rspBuf[:n]))
	}
}
