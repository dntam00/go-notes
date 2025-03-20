package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8085")
	if err != nil {
		fmt.Printf("Failed to connect to server: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server at port 8085")

	// Create a reader to read input from the console
	reader := bufio.NewReader(os.Stdin)

	for {
		// Read input from the console
		fmt.Print("Enter message: ")
		message, _ := reader.ReadString('\n')

		// Send the message to the server
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Printf("Failed to send message: %v\n", err)
			return
		}

		// Read the response from the server
		response := make([]byte, 1024)
		n, err := conn.Read(response)
		if err != nil {
			fmt.Printf("Failed to read response: %v\n", err)
			return
		}

		// Print the server's response
		fmt.Printf("Server response: %s\n", string(response[:n]))
	}
}
