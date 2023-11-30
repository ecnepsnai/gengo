package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type TGenGoConfig struct {
	MinimumVersion string `json:"minimum_version" yaml:"minimum_version"`
	FilePrefix     string `json:"file_prefix" yaml:"file_prefix"`
}

var GenGoConfig *TGenGoConfig

func loadGenGoConfig(configDir string) {
	if !loadConfig(configDir, "codegen", GenGoConfig) {
		GenGoConfig = &TGenGoConfig{
			MinimumVersion: Version,
			FilePrefix:     "gengo_",
		}
	}

	versionStrToNumber := func(in string) int {
		v := strings.ReplaceAll(in[1:], ".", "")
		i, err := strconv.Atoi(v)
		if err != nil {
			return -1
		}
		return i
	}

	currentVersionNumber := versionStrToNumber(Version)
	minimumVersionNumber := versionStrToNumber(GenGoConfig.MinimumVersion)

	if minimumVersionNumber > currentVersionNumber {
		fmt.Fprintf(os.Stderr, "Incorrect GenGo version installed.\nWanted: %s\nInstalled: %s\n", GenGoConfig.MinimumVersion, Version)
		os.Exit(1)
	}
}

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
