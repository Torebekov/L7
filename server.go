package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

const (
	numHandlers = 1
	invNum       = "Invalid number."
	serverClosed = "Server is closed."
	seeYou       = "See you soon!"
)

func handleConnection(conn net.Conn, numHandling, stopHandlers <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() { <-numHandling }()
	defer conn.Close()

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		message := scanner.Text()
		fmt.Println("Message Received:", message)
		if message == "close" {
			conn.Write([]byte(seeYou + "\n"))
			break
		}
		num, err := strconv.Atoi(message)
		if err != nil {
			select {
			case <-stopHandlers:
				conn.Write([]byte(invNum + serverClosed + seeYou + "\n"))
				return
			default:
				conn.Write([]byte(invNum + "\n"))
			}
			continue
		}
		select {
		case <-stopHandlers:
			conn.Write([]byte(strconv.Itoa(num*num) + ". " + serverClosed + seeYou + "\n"))
			return
		default:
			conn.Write([]byte(strconv.Itoa(num*num) + "\n"))
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error:", err)
	}
}
func main() {
	fmt.Println("Launching server...")

	ln, _ := net.Listen("tcp", ":8081")

	var wg sync.WaitGroup
	numHandling := make(chan struct{}, numHandlers)
	stopHandlers := make(chan struct{})

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func(stop <-chan os.Signal, stopHandlers chan struct{}) {
		<-stop
		close(stopHandlers)
		err := ln.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(stop, stopHandlers)

Loop:
	for {
		numHandling <- struct{}{}
		conn, err := ln.Accept()
		if err != nil {
			select {
			case <-stopHandlers:
				break Loop
			default:
				log.Fatal(err)
			}
		}
		wg.Add(1)
		fmt.Println("Connection is handling")
		go handleConnection(conn, numHandling, stopHandlers, &wg)
	}
	wg.Wait()
	fmt.Println(serverClosed)
}
