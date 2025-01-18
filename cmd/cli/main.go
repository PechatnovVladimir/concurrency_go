package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:12345")
	if err != nil {
		fmt.Printf("Failed to connect to server: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server. Enter commands (SET, GET, DEL) or type EXIT to quit.")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		command := scanner.Text()
		if strings.ToUpper(command) == "EXIT" {
			fmt.Println("Exiting...")
			break
		}

		_, err := conn.Write([]byte(command + "\n"))
		if err != nil {
			fmt.Printf("Failed to send command: %v\n", err)
			break
		}

		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Printf("Failed to read response: %v\n", err)
			break
		}
		fmt.Printf("%s", response)
	}
}
