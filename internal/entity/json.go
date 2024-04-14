package entity

import (
	"encoding/json"
)

// MapToJSON преобразует карту в формат JSON
func MapToJSON(content map[string]interface{}) (string, error) {
	jsonData, err := json.Marshal(content)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
