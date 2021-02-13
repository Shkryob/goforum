package utils

import (
	"encoding/json"
	"github.com/labstack/gommon/log"
)

type result_type map[string]interface{}

func JsonToMap(json_data []byte) result_type {
	var result result_type
	err := json.Unmarshal([]byte(json_data), &result)
	if err != nil {
		log.Info(err)
	}

	return result
}
