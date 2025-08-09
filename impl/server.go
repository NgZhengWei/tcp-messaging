package impl

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/NgZhengWei/tcp-messaging/config"
)

type Server struct {
	Listener      net.Listener
	ClientManager ClientManager
}

type ClientManager struct {
	clients map[string]Client // username - Client
	sync.RWMutex
}

type Client struct {
	username string
	conn     net.Conn
}

func NewServer() *Server {
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic(err)
	}
	fmt.Print("Server started on port 8000...")

	server := &Server{
		Listener: listener,
		ClientManager: ClientManager{
			clients: make(map[string]Client),
		},
	}

	return server
}

// HandleConnection handles incoming client connections
func (s *Server) HandleConnection(conn net.Conn) {
	defer conn.Close()

	var (
		username string
		err      error
	)

	// read from the connection into a buffer
	reader := bufio.NewReader(conn)

	username, err = checkUsername(&s.ClientManager, reader, conn)
	if err != nil {
		if err.Error() == config.ErrUsernameTaken {
			conn.Write([]byte(config.ErrUsernameTaken + "\n"))
			fmt.Println("Username is already taken. Closing connection.")
			return
		}

		log.Printf("Error reading username: %v", err)
		return
	} else {
		conn.Write([]byte("Welcome, " + username + "\n"))
		fmt.Printf("Client connected: %s with username: %s\n", conn.RemoteAddr(), username)
	}

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Printf("Client disconnected: %s\n", conn.RemoteAddr())
			} else {
				log.Printf("Error reading from connection: %v", err)
			}

			return
		}

		cleanMsg := strings.TrimSpace(message)
		fmt.Println("Received message:", cleanMsg)

		conn.Write([]byte(cleanMsg + "\n"))
	}
}

func checkUsername(cm *ClientManager, reader *bufio.Reader, conn net.Conn) (string, error) {
	username, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading username: %v", err)
		return "", err
	}
	username = strings.TrimSpace(username)

	cm.Lock()
	defer cm.Unlock()
	if _, exists := cm.clients[username]; exists {
		return "", fmt.Errorf(config.ErrUsernameTaken)
	}
	cm.clients[username] = Client{
		username: username,
		conn:     conn,
	}

	return username, nil
}
