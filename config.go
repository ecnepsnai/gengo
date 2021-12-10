package main

import (
	"encoding/json"
	"os"

	"gopkg.in/yaml.v2"
)

func loadConfig(fileName string, obj interface{}) bool {
	if fileExists(fileName + ".json") {
		f, err := os.Open(fileName + ".json")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if err := json.NewDecoder(f).Decode(obj); err != nil {
			panic(err)
		}

		return true
	} else if fileExists(fileName + ".yml") {
		f, err := os.Open(fileName + ".yml")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if err := yaml.NewDecoder(f).Decode(obj); err != nil {
			panic(err)
		}

		return true
	} else if fileExists(fileName + ".yaml") {
		f, err := os.Open(fileName + ".yaml")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if err := yaml.NewDecoder(f).Decode(obj); err != nil {
			panic(err)
		}

		return true
	}

	return false
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}
