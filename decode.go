package gonacos

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

// DecodeJson json解码
func DecodeJson(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}

// DecodeYaml yaml解码
func DecodeYaml(data string, v interface{}) error {
	return yaml.Unmarshal([]byte(data), v)
}
