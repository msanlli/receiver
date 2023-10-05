package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

func handleMessage(rawMessage []byte) {
	var msg Message
	err := json.Unmarshal(rawMessage, &msg) // Parse the message
	if err != nil {
		fmt.Println("Error parsing message:", err)
		return
	}

	switch msg.Type { // Check the message type
	case "alert":
		var alertPayload Alert
		err = json.Unmarshal(msg.Payload, &alertPayload) // Parse the payload
		if err != nil {
			fmt.Println("Error parsing alert payload:", err)
			return
		}
		handleAlert(alertPayload) // Handle the alert
	case "data":
		var dataPayload Data
		err = json.Unmarshal(msg.Payload, &dataPayload) // Parse the payload
		if err != nil {
			fmt.Println("Error parsing data payload:", err)
			return
		}
		handleData(dataPayload) // Handle the data
	default:
		fmt.Println("Unknown message type:", msg.Type)
		return
	}
}

// getFilename generates a filename with the given base name and extension incorporating the current UTC time.
func getFilename(baseName, extension string) string {
	currentTime := time.Now().UTC()
	timestamp := currentTime.Format("20060102-150405") // YYYYMMDD-HHMMSS format
	return fmt.Sprintf("%s-%s.%s", baseName, timestamp, extension)
}

// saveAsJSON saves data as JSON
func saveAsJSON(baseFilename string, data interface{}) error {
	filename := getFilename(baseFilename, "json")
	jsonData, err := json.MarshalIndent(data, "", "  ") // Indent the JSON
	if err != nil {
		return err
	}
	return os.WriteFile(filename, jsonData, 0644)
}

// saveAsYAML saves data as YAML
func saveAsYAML(baseFilename string, data interface{}) error {
	filename := getFilename(baseFilename, "yaml")
	yamlData, err := yaml.Marshal(data) // Convert the data to YAML
	if err != nil {
		return err
	}
	return os.WriteFile(filename, yamlData, 0644)
}

// saveAsTOML saves data as TOML
func saveAsTOML(baseFilename string, data interface{}) error {
	filename := getFilename(baseFilename, "toml")
	var buffer bytes.Buffer             // Create a buffer to write to
	encoder := toml.NewEncoder(&buffer) // Create a new encoder
	if err := encoder.Encode(data); err != nil {
		return err
	}
	return os.WriteFile(filename, buffer.Bytes(), 0644) // Write the buffer to a file
}

// handleMessage handles a message
func handleAlert(payload Alert) {
	timestamp := time.Unix(payload.Date, 0)
	fmt.Printf("Alert received at %s: %s\n", timestamp, payload.Event)

	// Save as JSON
	if err := saveAsJSON("alert.json", payload); err != nil {
		fmt.Println("Error saving JSON:", err)
	}

	// Save as YAML
	if err := saveAsYAML("alert.yaml", payload); err != nil {
		fmt.Println("Error saving YAML:", err)
	}

	// Save as TOML
	if err := saveAsTOML("alert.toml", payload); err != nil {
		fmt.Println("Error saving TOML:", err)
	}
}

// handleData handles a data payload
func handleData(payload Data) {
	fmt.Printf("Data received: %s = %f\n", payload.Name, payload.Value)

	// Save as JSON
	if err := saveAsJSON("data.json", payload); err != nil {
		fmt.Println("Error saving JSON:", err)
	}

	// Save as YAML
	if err := saveAsYAML("data.yaml", payload); err != nil {
		fmt.Println("Error saving YAML:", err)
	}

	// Save as TOML
	if err := saveAsTOML("data.toml", payload); err != nil {
		fmt.Println("Error saving TOML:", err)
	}
}
