package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	const serverIsClosed = "Server is closed. Try to open later, see you soon!"

	conn, err := net.Dial("tcp", "127.0.0.1:8081")
	if err != nil {
		fmt.Println(serverIsClosed)
		return
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Number to send: ")
		text, _ := reader.ReadString('\n')

		fmt.Fprintf(conn, text)

		message, _ := bufio.NewReader(conn).ReadString('\n')
		if message == "" {
			message = serverIsClosed
		}

		fmt.Print("Answer from server: " + message)
		if strings.Contains(message, "ee you soon!") {
			break
		}
	}
}
