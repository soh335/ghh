package main

import (
	"encoding/json"
	"os"
)

func ReadConfig(filePath string) (map[string]string, error) {
	m := make(map[string]string)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(file).Decode(&m); err != nil {
		return nil, err
	}
	return m, nil
}

func WriteConfig(filePath string, m map[string]string) error {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return json.NewEncoder(file).Encode(m)
}
