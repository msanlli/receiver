package main

import (
	"bufio"
	"fmt"
	"net"
)

func startTCPServer() {
	listener, err := net.Listen("tcp", ":8080") // Listens on port 8080
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection Error:", err)
			continue
		}
		go handleTCP(conn) // Handle connection in a new goroutine
	}
}

func startUDPServer() {
	addr, err := net.ResolveUDPAddr("udp", ":8081") // Listens on port 8081
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	buf := make([]byte, 1024)

	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Reading Error:", err)
			continue
		}

		message := string(buf[:n])
		fmt.Println("Received Message:", message)
	}
}

func handleTCP(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		fmt.Println("Message Received:", message)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading:", err)
	}
}
