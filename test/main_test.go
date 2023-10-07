// blackbox_test.go
package main

import (
	"encoding/json"
	"net"
	"testing"

	pkg "receiver.com/m/pkg"
)

func TestTCPHandling(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	message := &pkg.Message{
		Type:    "alert",
		Payload: json.RawMessage(`{"date": 1674832322, "event": "test event"}`),
	}

	msgBytes, err := json.Marshal(message)
	if err != nil {
		t.Fatalf("Failed to marshal message: %v", err)
	}

	_, err = conn.Write(msgBytes)
	if err != nil {
		t.Fatalf("Failed to write to server: %v", err)
	}

}

func TestUDPHandling(t *testing.T) {
	conn, err := net.Dial("udp", "localhost:8081")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	message := &pkg.Message{
		Type:    "data",
		Payload: json.RawMessage(`{"name": "test data", "value": 42.0}`),
	}

	msgBytes, err := json.Marshal(message)
	if err != nil {
		t.Fatalf("Failed to marshal message: %v", err)
	}

	_, err = conn.Write(msgBytes)
	if err != nil {
		t.Fatalf("Failed to write to server: %v", err)
	}
}
