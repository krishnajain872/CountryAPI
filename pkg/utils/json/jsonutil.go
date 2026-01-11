package jsonutil

import (
	"encoding/json"
	"fmt"
)

// Marshal converts any struct to JSON string
func Marshal[T any](v T) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("json marshal failed: %w", err)
	}
	return string(data), nil
}

// Unmarshal converts JSON string to struct
func Unmarshal[T any](data string) (*T, error) {
	var v T
	if err := json.Unmarshal([]byte(data), &v); err != nil {
		return nil, fmt.Errorf("json unmarshal failed: %w", err)
	}
	return &v, nil
}
