package main

import (
	"log"
	"net"
	//	"os"
)

func main() {
	ln, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Server Listening on %s", ln.Addr())

	msgchan := make(chan string)
	go printMessages(msgchan)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn, msgchan)
	}
}

func handleConnection(c net.Conn, msgchan chan<- string) {
	buf := make([]byte, 4096)
	log.Printf("New connection from %s", c.RemoteAddr())
	for {
		n, err := c.Read(buf)
		if err != nil || n == 0 {
			c.Close()
			break
		}

		msgchan <- string(buf[0:n])
		n, err = c.Write(buf[0:n])
		if err != nil {
			c.Close()
			break
		}
	}
	log.Printf("Connection from %v closed.", c.RemoteAddr())
}

func printMessages(msgchan <-chan string) {
	for msg := range msgchan {
		log.Printf("new message: %s", msg)
	}
}
