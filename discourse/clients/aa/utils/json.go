package utils

import (
	"encoding/json"
	"log"
)

func PrettifyJson(data []byte) string {
	var j interface{}
	json.Unmarshal(data, &j)
	pj, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return string(pj)
}
