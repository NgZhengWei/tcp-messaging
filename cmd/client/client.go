// made this to learn tcp
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/NgZhengWei/tcp-messaging/config"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server:", conn.RemoteAddr())

	stdinReader := bufio.NewReader(os.Stdin)
	connReader := bufio.NewReader(conn)

	username, err := getUsername(stdinReader)
	if err != nil {
		log.Fatalf("Error getting username: %v", err)
		return
	}

	// send the username to the server
	_, err = conn.Write([]byte(username + "\n"))
	if err != nil {
		log.Printf("Error sending username: %v", err)
		return
	}
	res, err := connReader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading response: %v", err)
		return
	}
	if strings.TrimSpace(res) == config.ErrUsernameTaken {
		fmt.Println("Username is already taken. Please try again.")
		return
	}

	fmt.Println("Welcome,", username)

	// loop to read from stdin and send to the connection
	for {
		fmt.Print("> ")

		// read user input from stdin
		input, err := stdinReader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading from stdin: %v", err)
			return
		}

		input = strings.TrimSpace(input)
		if input == "exit" {
			fmt.Println("Exiting...")
			return
		}

		// write user input to the connection
		_, err = conn.Write([]byte(input + "\n"))
		if err != nil {
			log.Printf("Error writing to connection: %v", err)
			return
		}

		// read response from the connection
		response, err := connReader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading from connection: %v", err)
			return
		}

		fmt.Println("Received response:", response)
	}
}

func getUsername(stdinReader *bufio.Reader) (string, error) {
	fmt.Print("Enter your username: ")
	username, err := stdinReader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading username: %v", err)
		return "", err
	}

	username = strings.TrimSpace(username)
	if username == "" {
		log.Fatal("Username cannot be empty")
		return "", fmt.Errorf("username cannot be empty")
	}

	return username, nil
}
