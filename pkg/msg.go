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

type Record struct {
	json *os.File
	yaml *os.File
	toml *os.File
}

func NewDataRecord(dir string) *Record {
	return newRecord(dir, "data")
}

func NewAlertRecord(dir string) *Record {
	return newRecord(dir, "alerts")
}

func newRecord(dir, baseFilename string) *Record {
	r := &Record{}
	r.json = openFile(path.Join(dir, baseFilename+".json"))
	r.yaml = openFile(path.Join(dir, baseFilename+".yaml"))
	r.toml = openFile(path.Join(dir, baseFilename+".toml"))
	return r
}

func (record *Record) Save(data interface{}) {
	// Save as JSON
	jsonData, _ := json.MarshalIndent(data, "", "  ")
	record.json.Write(jsonData)

	// Save as YAML
	yamlData, _ := yaml.Marshal(data)
	record.yaml.Write(yamlData)

	// Save as TOML
	var buffer bytes.Buffer
	encoder := toml.NewEncoder(&buffer)
	encoder.Encode(data)
	record.toml.Write(buffer.Bytes())
}

func openFile(filename string) *os.File {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	CheckError(err)
	return file
}

func CheckError(err error) error {
	if err != nil {
		return err
	}
	return nil
}

func HandleMessage(rawMessage []byte) {
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
		handleAlert(alertPayload)
	case "data":
		var dataPayload Data
		err = json.Unmarshal(msg.Payload, &dataPayload)
		if err != nil {
			println("Error parsing data payload", err)
			return
		}
		handleData(dataPayload)
	default:
		fmt.Println("Unknown message type:", msg.Type)
		return
	}
}

func handleAlert(payload Alert) {
	timestamp := time.Unix(payload.Date, 0)
	fmt.Printf("Alert received at %s: %s\n", timestamp, payload.Event)
	record := NewAlertRecord("alerts")
	defer closeFiles(record)
	record.Save(payload)
}

func handleData(payload Data) {
	fmt.Printf("Data received: %s = %f\n", payload.Name, payload.Value)
	record := NewDataRecord("data")
	defer closeFiles(record)
	record.Save(payload)
}

func closeFiles(r *Record) {
	r.json.Close()
	r.yaml.Close()
	r.toml.Close()
}
