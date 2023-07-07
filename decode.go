package gonacos

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

// json解码
func DecodeJson(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}

// yaml解码
func DecodeYaml(data string, v interface{}) error {
	return yaml.Unmarshal([]byte(data), v)
}
