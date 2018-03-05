package utils

import (
	"encoding/json"
)

func GetAsJsonString(v interface{}) (string) {
	j, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		panic (err)
	}
	return string(j)
}
