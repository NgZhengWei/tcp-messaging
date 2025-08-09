package main

import (
	"fmt"
	"log"

	"github.com/NgZhengWei/tcp-messaging/impl"
)

func main() {
	server := impl.NewServer()

	for {
		conn, err := server.Listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		fmt.Printf("Client connected: %s\n", conn.RemoteAddr())
		go server.HandleConnection(conn)
	}
}
