package utils

import (
	"encoding/json"
)

type result_type map[string]interface{}

func JsonToMap(json_data []byte) result_type {
	var result result_type
	json.Unmarshal([]byte(json_data), &result)

	return result
}
