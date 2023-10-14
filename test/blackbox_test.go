package test_main

import (
	"net"
	"os"
	"testing"
	"time"

	pkg "receiver.com/m/pkg"
)

func Blackbox(t *testing.T) {

	go pkg.Main()

	time.Sleep(2 * time.Second)

	tcpMessages := []string{
		`{"type":"alert","payload":{"date":1673782920,"event":"fire detected"}}`,
		`{"type":"data","payload":{"name":"temperature","value":23.5}}`,
	}

	udpMessages := []string{
		`{"type":"alert","payload":{"date":1673782920,"event":"fire detected"}}`,
		`{"type":"data","payload":{"name":"temperature","value":23.5}}`,
	}

	// Send TCP messages
	for _, msg := range tcpMessages {
		conn, err := net.Dial("tcp", "localhost:8080")
		if err != nil {
			t.Fatalf("failed to connect to TCP server: %v", err)
		}
		_, err = conn.Write([]byte(msg))
		if err != nil {
			t.Fatalf("failed to send TCP message: %v", err)
		}
		conn.Close()
	}

	// Send UDP messages
	for _, msg := range udpMessages {
		conn, err := net.Dial("udp", "localhost:8081")
		if err != nil {
			t.Fatalf("failed to connect to UDP server: %v", err)
		}
		_, err = conn.Write([]byte(msg))
		if err != nil {
			t.Fatalf("failed to send UDP message: %v", err)
		}
		conn.Close()
	}
	alertFiles := []string{
		"./alert/alert.json",
		"./alert/alert.yaml",
		"./alert/alert.toml",
	}

	dataFiles := []string{
		"./data/data.json",
		"./data/data.yaml",
		"./data/data.toml",
	}

	expectedAlertContent := `{"date":1673782920,"event":"fire detected"}`
	expectedDataContent := `{"name":"temperature","value":23.5}`

	for _, filePath := range alertFiles {
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("failed to read %s: %v", filePath, err)
		}
		if string(content) != expectedAlertContent {
			t.Errorf("unexpected content in %s. got: %s, want: %s", filePath, content, expectedAlertContent)
		}
	}

	for _, filePath := range dataFiles {
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("failed to read %s: %v", filePath, err)
		}
		if string(content) != expectedDataContent {
			t.Errorf("unexpected content in %s. got: %s, want: %s", filePath, content, expectedDataContent)
		}
	}
}
