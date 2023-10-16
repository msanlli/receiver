package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

// The struct is used to open, parse and save data to files.
type Record struct {
	json *os.File
	yaml *os.File
	toml *os.File
}

func (r *Record) CloseAll() {
	if r.json != nil {
		r.json.Close()
	}
	if r.yaml != nil {
		r.yaml.Close()
	}
	if r.toml != nil {
		r.toml.Close()
	}
}

func NewAlertRecord(dir []string) (*Record, error) {
	alertDir := dir[0:3]
	return NewRecord(alertDir)
}

func NewDataRecord(dir []string) (*Record, error) {
	dataDir := dir[3:6]
	return NewRecord(dataDir)
}

// NewRecord opens all the files and returns a Record struct.
func NewRecord(files []string) (*Record, error) {
	r := &Record{}

	for _, file := range files {
		var f *os.File
		var err error

		f, err = os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			r.CloseAll()
			return nil, err
		}

		switch {
		case strings.HasSuffix(file, ".json"):
			r.json = f
		case strings.HasSuffix(file, ".yaml"):
			r.yaml = f
		case strings.HasSuffix(file, ".toml"):
			r.toml = f
		}
	}

	return r, nil
}

// SaveMessage saves data in all three file formats.
func SaveMessage(data interface{}, fileType string, record *Record) {
	// Save as JSON
	jsonData, err := json.MarshalIndent(data, "", "")
	if err == nil && record.json != nil {
		_, err = record.json.Write(jsonData)
		if err != nil {
			fmt.Println("Error writing JSON data:", err)
		}
		record.json.WriteString("\n")
	}

	// Save as YAML
	yamlData, err := yaml.Marshal(data)
	if err == nil && record.yaml != nil {
		_, err = record.yaml.Write(yamlData)
		if err != nil {
			fmt.Println("Error writing YAML data:", err)
		}
	}

	// Save as TOML
	var buffer bytes.Buffer
	encoder := toml.NewEncoder(&buffer)
	if err := encoder.Encode(data); err == nil && record.toml != nil {
		_, err = record.toml.Write(buffer.Bytes())
		if err != nil {
			fmt.Println("Error writing TOML data:", err)
		}
	}
}

// HandleMessage handles a message received from the network
// and distinguises between alert and data messages.
func HandleMessage(rawMessage []byte, alertfiles, datafiles *Record) {
	var msg Message
	err := json.Unmarshal(rawMessage, &msg)
	if err != nil {
		println("Error parsing message:", err)
		return
	}

	switch msg.Type {
	case "alert":
		var alertPayload Alert
		err = json.Unmarshal(msg.Payload, &alertPayload)
		if err != nil {
			println("Error parsing alert payload:", err)
			return
		}
		timestamp := time.Unix(alertPayload.Date, 0)
		fmt.Printf("Alert received at %s: %s\n", timestamp, alertPayload.Event)

		SaveMessage(alertPayload, "alert", alertfiles)

	case "data":
		var dataPayload Data
		err = json.Unmarshal(msg.Payload, &dataPayload)
		if err != nil {
			println("Error parsing data payload", err)
			return
		}

		fmt.Printf("Data received: %s = %f\n", dataPayload.Name, dataPayload.Value)

		SaveMessage(dataPayload, "data", datafiles)

	default:
		fmt.Println("Unknown message type:", msg.Type)
		return
	}
}
