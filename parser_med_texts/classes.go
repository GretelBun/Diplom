package main

import (
	"encoding/json"
	"log"
)

type Class struct {
	Name     string   `json:"class"`
	Synonyms []string `json:"synonyms"`
}

func ParseJson(data []byte) (classes []Class) {
	err := json.Unmarshal(data, &classes)
	if err != nil {
		log.Fatal(err)
	}
	return
}
