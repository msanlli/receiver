package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

var LogBuffer bytes.Buffer

func CheckError(err error) {
	if err != nil {
		LogBuffer.WriteString("Error: ")
		os.Exit(0)
	}
}

func HandleMessage(rawMessage []byte) {
	var msg Message
	err := json.Unmarshal(rawMessage, &msg) // Parse the message
	if err != nil {
		LogBuffer.WriteString("Error parsing message:")
		println("Error parsing message:", err)
		return
	}

	switch msg.Type { // Check the message type
	case "alert":
		var alertPayload Alert
		err = json.Unmarshal(msg.Payload, &alertPayload) // Parse the payload
		if err != nil {
			LogBuffer.WriteString("Error parsing alert payload:")
			println("Error parsing alert payload:", err)
			return
		}
		handleAlert(alertPayload) // Handle the alert
	case "data":
		var dataPayload Data
		err = json.Unmarshal(msg.Payload, &dataPayload) // Parse the payload
		if err != nil {
			LogBuffer.WriteString("Error parsing data payload:")
			println("Error parsing data payload", err)
			return
		}
		handleData(dataPayload) // Handle the data
	default:
		fmt.Println("Unknown message type:", msg.Type)
		return
	}
}

// saveAsJSON saves data as JSON
func saveAsJSON(directory string, filename string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ") // Indent the JSON
	if err != nil {
		return err
	}
	var filenameE = path.Join(directory, filename+".json")
	file, err := os.OpenFile(filenameE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 064)
	CheckError(err)
	defer file.Close()
	_, err = file.Write(jsonData)
	return err
}

// saveAsYAML saves data as YAML
func saveAsYAML(directory string, filename string, data interface{}) error {

	yamlData, err := yaml.Marshal(data) // Convert the data to YAML
	if err != nil {
		return err
	}
	var filenameE = path.Join(directory, filename+".yaml")
	file, err := os.OpenFile(filenameE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 064)
	CheckError(err)
	defer file.Close()
	_, err = file.Write(yamlData)
	return err
}

// saveAsTOML saves data as TOML
func saveAsTOML(directory string, filename string, data interface{}) error {
	var buffer bytes.Buffer             // Create a buffer to write to
	encoder := toml.NewEncoder(&buffer) // Create a new encoder
	if err := encoder.Encode(data); err != nil {
		return err
	}
	var filenameE = path.Join(directory, filename+".toml")
	file, err := os.OpenFile(filenameE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 064)
	CheckError(err)
	defer file.Close()
	_, err = file.Write(buffer.Bytes())
	return err
}

// handleMessage handles a message
func handleAlert(payload Alert) {
	timestamp := time.Unix(payload.Date, 0)
	fmt.Printf("Alert received at %s: %s\n", timestamp, payload.Event)

	// Save as JSON
	if err := saveAsJSON("alerts", "alerts", payload); err != nil {
		LogBuffer.WriteString("Error saving JSON:")
		println("Error saving JSON:", err)
	}

	// Save as YAML
	if err := saveAsYAML("alerts", "alerts", payload); err != nil {
		LogBuffer.WriteString("Error saving YAML:")
		println("Error YAML", err)
	}

	// Save as TOML
	if err := saveAsTOML("alerts", "alerts", payload); err != nil {
		LogBuffer.WriteString("Error saving TOML:")
		println("Error saving TOML:", err)
	}
}

// handleData handles a data payload
func handleData(payload Data) {
	fmt.Printf("Data received: %s = %f\n", payload.Name, payload.Value)

	// Save as JSON
	if err := saveAsJSON("data", "data", payload); err != nil {
		LogBuffer.WriteString("Error saving JSON:")
		println("Error saving JSON:", err)
	}

	// Save as YAML
	if err := saveAsYAML("data", "data", payload); err != nil {
		LogBuffer.WriteString("Error saving YAML:")
		println("Error savings YAML:", err)
	}

	// Save as TOML
	if err := saveAsTOML("data", "data", payload); err != nil {
		LogBuffer.WriteString("Error saving TOML:")
		println("Error saving TOML:", err)
	}
}
