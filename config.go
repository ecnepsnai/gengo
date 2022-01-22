package main

import (
	"encoding/json"
	"os"

	"gopkg.in/yaml.v2"
)

func loadConfig(fileName string, obj interface{}) bool {
	var f *os.File
	var err error
	isJSON := false

	if fileExists(fileName + ".json") {
		f, err = os.Open(fileName + ".json")
		isJSON = true
	} else if fileExists(fileName + ".yml") {
		f, err = os.Open(fileName + ".yml")
	} else if fileExists(fileName + ".yaml") {
		f, err = os.Open(fileName + ".yaml")
	} else {
		return false
	}

	if err != nil {
		panic(err)
	}

	if isJSON {
		if err := json.NewDecoder(f).Decode(obj); err != nil {
			panic(err)
		}
		return true
	}
	if err := yaml.NewDecoder(f).Decode(obj); err != nil {
		panic(err)
	}
	return true
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}
