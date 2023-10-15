package pkg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Alert struct {
	Date  int64  `json:"date"`
	Event string `json:"event"`
}

type Data struct {
	Name  string  `json:"name"`
	Value float32 `json:"value"`
}

// startTCP starts a TCP listener
func startTCP() {
	fmt.Println("Starting TCP listener")
	addrTCP, err := net.ResolveTCPAddr("tcp", ":8080")
	if err != nil {
		println("Error starting UDP server:", err)
		return
	}

	listener, err := net.ListenTCP("tcp", addrTCP)
	if err != nil {
		println("Error starting TCP:", err)
		return
	}
	defer listener.Close()

	for {
		fmt.Printf("Listening on %s\n", listener.Addr().String())
		conn, err := listener.Accept()
		if err != nil {
			println("Error starting TCP connection:", err)
			continue
		}
		handleTCP(conn) // Handle connection in a new goroutine
	}
}

// startUDP starts a UDP listener
func startUDP() {
	fmt.Println("Starting UDP listener")
	addrUDP, err := net.ResolveUDPAddr("udp", ":8081") // Listens on port 8081
	if err != nil {
		println("Error starting UDP server:", err)
		return
	}

	listener, err := net.ListenUDP("udp", addrUDP) // Listen on the port
	if err != nil {
		println("Error starting UDP listener:", err)
		return
	}
	defer listener.Close() // Close the connection when the function returns

	for {
		fmt.Printf("Listening on %s\n", listener.LocalAddr().String())
		handleUDP(listener) // Handle the received data
	}
}

// handleTCP handles a TCP connection
func handleTCP(conn net.Conn) {
	scanner := bufio.NewScanner(conn) // Create a scanner
	for scanner.Scan() {              // Scan the connection
		rawMessage := scanner.Bytes() // Get the message as bytes
		HandleMessage(rawMessage)     // Handle the message
		fmt.Println("TCP Message Received:", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		println("Error scanning:", err)
		return
	}
}

func handleUDP(conn net.Conn) {
	// Handle UDP messages
	if udpConn, ok := conn.(*net.UDPConn); ok {
		buf := make([]byte, 1024)
		n, _, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			println("Error starting UDP connection")
			return
		}

		rawMessage := buf[:n]
		HandleMessage(rawMessage)
		fmt.Println("UDP Message Received:", string(rawMessage))
		return
	}

	// Unsupported connection type
	println("Unsupported connection type: ")
}
