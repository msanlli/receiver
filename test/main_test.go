package pkg

import (
	"os"
	"testing"

	pkg "receiver.com/m/pkg"
)

func TestHandleMessage(t *testing.T) {
	// Sample input messages
	alertMessage := `{
		"Type": "alert",
		"Payload": {
			"Event": "Fire",
			"Date": 1633801090
		}
	}`
	dataMessage := `{
		"Type": "data",
		"Payload": {
			"Name": "Temperature",
			"Value": 32.5
		}
	}`

	tests := []struct {
		input string
	}{
		{alertMessage},
		{dataMessage},
	}

	for _, tt := range tests {
		// Create temporary files for JSON, YAML, and TOML
		jsonFile, err := os.CreateTemp("", "test-*.json")
		if err != nil {
			t.Fatalf("Failed to create temp JSON file: %v", err)
		}

		yamlFile, err := os.CreateTemp("", "test-*.yaml")
		if err != nil {
			t.Fatalf("Failed to create temp YAML file: %v", err)
		}

		tomlFile, err := os.CreateTemp("", "test-*.toml")
		if err != nil {
			t.Fatalf("Failed to create temp TOML file: %v", err)
		}

		// Call HandleMessage (Note: this assumes it uses the above files to write logs.
		// If not, then we need some mechanism to redirect its output to these files)
		pkg.HandleMessage([]byte(tt.input))

		// Cleanup: Close and remove temporary files
		jsonFile.Close()
		yamlFile.Close()
		tomlFile.Close()

		os.Remove(jsonFile.Name())
		os.Remove(yamlFile.Name())
		os.Remove(tomlFile.Name())
	}
}
