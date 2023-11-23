package main

import (
	"fmt"
	"net"
)

type Client struct {
	conn net.Conn
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ln.Close()

	clients := make(map[net.Conn]bool)
	messages := make(chan string)
	newClients := make(chan Client)
	deadClients := make(chan Client)

	// Handle new messages and client connections/disconnections
	go func() {
		for {
			select {
			case msg := <-messages:
				// Broadcast the message to all clients
				for conn := range clients {
					fmt.Fprint(conn, msg)
				}

			case newClient := <-newClients:
				// Add new client
				clients[newClient.conn] = true
				fmt.Printf("[INFO]: client connected from: %s\n", newClient.conn.RemoteAddr().String())

			case deadClient := <-deadClients:
				// Remove disconnected client
				delete(clients, deadClient.conn)
				fmt.Printf("[INFO]: client disconnected from: %s\n", deadClient.conn.RemoteAddr().String())
			}
		}
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		newClients <- Client{conn: conn}
		go handleConnection(conn, messages, deadClients)
	}
}

func handleConnection(conn net.Conn, messages chan<- string, deadClients chan<- Client) {
	defer func() {
		deadClients <- Client{conn: conn}
		conn.Close()
	}()

	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			return
		}
		messages <- string(buffer[:n])
	}
}
