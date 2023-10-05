package main

import "encoding/json"

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

func main() {
	go startTCP()
	startUDP()
}
