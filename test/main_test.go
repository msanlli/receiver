package test

import (
	"os"
	"testing"

	"receiver.com/m/pkg"
)

func readFileContent(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func TestHandleMessage(t *testing.T) {
	tests := []struct {
		input        []byte
		expectedJSON string
		expectedYAML string
		expectedTOML string
		jsonFilePath string
		yamlFilePath string
		tomlFilePath string
	}{
		{
			input: []byte(`{
                "Type": "alert",
                "Payload": "{\"Date\":42,\"Event\":\"Test\"}"
            }`),
			expectedJSON: "Alert received",
			expectedYAML: "Alert received",
			expectedTOML: "Alert received",
			jsonFilePath: "./alert/alert.json",
			yamlFilePath: "./alert/alert.yaml",
			tomlFilePath: "./alert/alert.toml",
		},
		{
			input: []byte(`{
                "Type": "data",
                "Payload": "{\"Name\":\"Test\",\"Value\":42.42}"
            }`),
			expectedJSON: "Alert received",
			expectedYAML: "Alert received",
			expectedTOML: "Alert received",
			jsonFilePath: "./data/data.json",
			yamlFilePath: "./data/data.yaml",
			tomlFilePath: "./data/data.toml",
		},
	}

	for _, test := range tests {
		pkg.HandleMessage(test.input)

		// Check the content of the JSON file for validation
		contentJSON, err := readFileContent(test.jsonFilePath)
		if err != nil {
			t.Fatal(err)
		}

		if contentJSON != test.expectedJSON {
			t.Errorf("(JSON) For input %s, expected %s, but got %s", string(test.input), test.expectedJSON, contentJSON)
		}

		contentYAML, err := readFileContent(test.yamlFilePath)
		if err != nil {
			t.Fatal(err)
		}

		if contentYAML != test.expectedYAML {
			t.Errorf("(YAML) For input %s, expected %s, but got %s", string(test.input), test.expectedYAML, contentYAML)
		}

		contentTOML, err := readFileContent(test.jsonFilePath)
		if err != nil {
			t.Fatal(err)
		}

		if contentTOML != test.expectedTOML {
			t.Errorf("(TOML) For input %s, expected %s, but got %s", string(test.input), test.expectedTOML, contentTOML)
		}
		// Similarly, you can add checks for the yaml and toml files if needed

		// Delete the files after checking
		//os.Remove(test.jsonFilePath)
		//os.Remove(test.yamlFilePath)
		//os.Remove(test.tomlFilePath)
	}
}
