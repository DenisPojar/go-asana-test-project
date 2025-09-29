package storage

import (
	"encoding/json"
	"os"
)

// SaveJSON saves any Go value to a JSON file
func SaveJSON(fileName string, data interface{}) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // pretty print
	return encoder.Encode(data)
}
