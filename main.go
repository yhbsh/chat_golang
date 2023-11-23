package main

import (
	"fmt"
	"net"
	"sync"
)

var (
	clients = make([]net.Conn, 0)
	mu      sync.Mutex // Mutex to protect access to clients slice
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("[INFO]: client connected from: %s\n", conn.RemoteAddr().String())

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	defer removeClient(conn)

	buffer := make([]byte, 1024)

	addClient(conn)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println(err)
			return // Exit the loop if an error occurs
		}

		mu.Lock()
		for _, c := range clients {
			if c != conn { // Avoid sending the message back to the sender
				c.Write(buffer[:n]) // Send only the actual data read
			}
		}
		mu.Unlock()
	}
}

func addClient(conn net.Conn) {
	mu.Lock()
	clients = append(clients, conn)
	mu.Unlock()
}

func removeClient(conn net.Conn) {
	mu.Lock()
	for i, c := range clients {
		if c == conn {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
	mu.Unlock()
}
