package pkg

import (
	"bufio"
	"fmt"
	"net"
)

// startTCP starts a TCP listener
func startTCP() {
	listener, err := net.Listen("tcp", ":8080") // Listens on port 8080
	if err != nil {
		LogBuffer.WriteString("Error: ")
		println("Error starting TCP:", err)
		return
	}
	defer listener.Close() // Close the listener when the function returns

	for {
		conn, err := listener.Accept() // Wait for a connection
		if err != nil {
			LogBuffer.WriteString("Error: ")
			println("Error starting TCP connection:", err)
			continue
		}
		go handleMsg(conn) // Handle connection in a new goroutine
	}
}

// startUDP starts a UDP listener
func startUDP() {
	addr, err := net.ResolveUDPAddr("udp", ":8081") // Listens on port 8081
	if err != nil {
		LogBuffer.WriteString("Error: ")
		println("Error starting UDP server:", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr) // Listen on the port
	if err != nil {
		LogBuffer.WriteString("Error: ")
		println("Error starting UDP listener:", err)
		return
	}
	defer conn.Close() // Close the listener when the function returns

	buf := make([]byte, 1024)

	for {
		n, _, err := conn.ReadFromUDP(buf) // Read from the connection
		if err != nil {
			LogBuffer.WriteString("Error: ")
			println("Error starting UDP connection:", err)
			continue
		}

		message := string(buf[:n]) // Convert the message to a string
		fmt.Println("Received Message:", message)
		handleMsg(conn) // Handle the message
	}
}

// handleTCP handles a TCP connection
func handleMsg(conn net.Conn) {
	if tcpConn, ok := conn.(*net.TCPConn); ok { // Check if the connection is TCP
		defer tcpConn.Close()
		scanner := bufio.NewScanner(tcpConn) // Create a scanner
		for scanner.Scan() {                 // Scan the connection
			rawMessage := scanner.Bytes() // Get the message as bytes
			HandleMessage(rawMessage)     // Handle the message
			fmt.Println("TCP Message Received:", scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			LogBuffer.WriteString("Error: ")
			println("Error scanning:", err)
			return
		}
		return
	}

	// Handle UDP messages
	if udpConn, ok := conn.(*net.UDPConn); ok {
		buf := make([]byte, 1024)
		n, _, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			LogBuffer.WriteString("Error: ")
			println("Error starting UDP connection")
			return
		}

		rawMessage := buf[:n]
		HandleMessage(rawMessage)
		fmt.Println("UDP Message Received:", string(rawMessage))
		return
	}

	// Unsupported connection type
	LogBuffer.WriteString("Unsupported connection type: ")
	println("Unsupported connection type: ")
}
