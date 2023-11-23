package main

import (
	"fmt"
	"net"
)

type Client struct {
	conn net.Conn
	ch   chan []byte
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ln.Close()

	newClients := make(chan Client)
	deadClients := make(chan Client)
	messages := make(chan []byte)

	go func() {
		clients := make(map[net.Conn]Client)
		for {
			select {
			case msg := <-messages:
				// Broadcast message to all clients
				for _, cli := range clients {
					cli.ch <- msg
				}
			case newCli := <-newClients:
				// Add new client
				clients[newCli.conn] = newCli
				fmt.Printf("[INFO]: client connected from: %s\n", newCli.conn.RemoteAddr().String())
			case deadCli := <-deadClients:
				// Remove client
				delete(clients, deadCli.conn)
				close(deadCli.ch)
				fmt.Printf("[INFO]: client disconnected from: %s\n", deadCli.conn.RemoteAddr().String())
			}
		}
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		clientCh := make(chan []byte)
		client := Client{conn, clientCh}
		newClients <- client

		go handleConnection(client, messages, deadClients)
	}
}

func handleConnection(client Client, messages chan []byte, deadClients chan Client) {
	defer func() {
		deadClients <- client
		client.conn.Close()
	}()

	buffer := make([]byte, 1024)

	for {
		n, err := client.conn.Read(buffer)
		if err != nil {
			return
		}

		messages <- buffer[:n]
	}
}
