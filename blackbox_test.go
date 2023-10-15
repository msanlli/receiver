package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	pkg "receiver.com/m/pkg"
)

func checkFiles(file string, expectedContent string) error {
	contentBytes, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	contentLines := strings.Split(strings.TrimSpace(string(contentBytes)), "\n")
	expectedLines := strings.Split(strings.TrimSpace(expectedContent), "\n")

	if len(contentLines) != len(expectedLines) {
		return fmt.Errorf("different number of lines: expected %d, got %d", len(expectedLines), len(contentLines))
	}

	for i, line := range contentLines {
		if strings.TrimSpace(line) != strings.TrimSpace(expectedLines[i]) {
			return fmt.Errorf("line %d mismatch: expected %s, got %s", i+1, expectedLines[i], line)
		}
	}
	return nil
}

func TestBlackbox(t *testing.T) {

	go pkg.Main()

	time.Sleep(5 * time.Second)

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

	expectedContents := map[string]map[string]string{
		"alert": {
			"json": `{
				"date": 1673782920,
				"event": "fire detected"
				}
				{
				"date": 1673782920,
				"event": "fire detected"
				}
				`,
			"yaml": `date: 1673782920
			event: fire detected
			date: 1673782920
			event: fire detected			
			`,
			"toml": `Date = 1673782920
			Event = 'fire detected'
			Date = 1673782920
			Event = 'fire detected'
			`,
		},

		"data": {
			"json": `{
				"name": "temperature",
				"value": 23.5
				}
				{
				"name": "temperature",
				"value": 23.5
				}
				`,
			"yaml": `name: temperature
			value: 23.5
			name: temperature
			value: 23.5
			`,
			"toml": `Name = 'temperature'
			Value = 23.5
			Name = 'temperature'
			Value = 23.5
			`,
		},
	}

	if err := checkFiles("./alert/alert.json", expectedContents["alert"]["json"]); err != nil {
		t.Error(err)
	}
	if err := checkFiles("./data/data.json", expectedContents["data"]["json"]); err != nil {
		t.Error(err)
	}
	if err := checkFiles("./alert/alert.yaml", expectedContents["alert"]["yaml"]); err != nil {
		t.Error(err)
	}
	if err := checkFiles("./data/data.yaml", expectedContents["data"]["yaml"]); err != nil {
		t.Error(err)
	}
	if err := checkFiles("./alert/alert.toml", expectedContents["alert"]["toml"]); err != nil {
		t.Error(err)
	}
	if err := checkFiles("./data/data.toml", expectedContents["data"]["toml"]); err != nil {
		t.Error(err)
	}
}

func Benchmark(b *testing.B) {
	// Start the program in the background
	go main()

	// Give some time for servers to start
	time.Sleep(5 * time.Second)

	// Run the benchmark loop
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
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
				b.Fatalf("failed to connect to TCP server: %v", err)
			}
			_, err = conn.Write([]byte(msg))
			if err != nil {
				b.Fatalf("failed to send TCP message: %v", err)
			}
			conn.Close()
		}

		// Send UDP messages
		for _, msg := range udpMessages {
			conn, err := net.Dial("udp", "localhost:8081")
			if err != nil {
				b.Fatalf("failed to connect to UDP server: %v", err)
			}
			_, err = conn.Write([]byte(msg))
			if err != nil {
				b.Fatalf("failed to send UDP message: %v", err)
			}
			conn.Close()
		}
	}
}
