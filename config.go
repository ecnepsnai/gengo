package main

import (
	"encoding/json"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

func loadConfig(configDir, fileName string, obj interface{}) bool {
	var f *os.File
	var err error
	isJSON := false

	filePath := path.Join(configDir, fileName)

	if fileExists(filePath + ".json") {
		f, err = os.Open(filePath + ".json")
		isJSON = true
	} else if fileExists(filePath + ".yml") {
		f, err = os.Open(filePath + ".yml")
	} else if fileExists(filePath + ".yaml") {
		f, err = os.Open(filePath + ".yaml")
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
