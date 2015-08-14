package main

import (
	"log"
	"net"
	//	"os"
)

type Client struct {
	conn net.Conn
	ch   chan<- string
}

func main() {
	ln, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Server Listening on %s", ln.Addr())

	msgchan := make(chan string)
	addchan := make(chan Client)
	rmchan := make(chan Client)

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

func handleMessages(msgchan <-chan string, addchan, rmchan <-chan Client) {
	clients := make(map[net.Conn]chan<- string)
	for {
		select {
		case msg := <-msgchan:
			for _, ch := range clients {
				go func(mch chan<- string) {
					mch <- "\033[1;33;40m" + msg + "\033[m\r\n" 
				}(ch)
			}
		case client:= <-addchan:
			clients[client.conn] = client.ch
			}
		}
	}
}

func handleConnection(c net.Conn, msgchan chan<- string, addchan, rmchan chan<- Client) {
	buf := make([]byte, 4096)
	log.Printf("New connection from %s", c.RemoteAddr())

	// Add the client to the channel
	ch := make(chan string)
	addchan <- Client{c, ch}

	defer func() {
		rmchan <- Client{c, ch}
	}()

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
