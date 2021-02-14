package utils

import (
	"encoding/json"
	"github.com/labstack/gommon/log"
)

type resultType map[string]interface{}

func JsonToMap(jsonData []byte) resultType {
	var result resultType
	err := json.Unmarshal([]byte(jsonData), &result)
	if err != nil {
		log.Info(err)
	}

	return result
}
