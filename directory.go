package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"text/template"

	"github.com/ecnepsnai/cbgen/templates"
)

// GenerateDirectory generates directory file
func GenerateDirectory(options Options) {
	directoryConfig := path.Join(".", "directory.json")
	directoryFile := path.Join(".", "cbgen_directory.go")

	if _, err := os.Stat(directoryConfig); err != nil {
		return
	}

	var directories []Directory
	data, err := ioutil.ReadFile(directoryConfig)
	if err != nil {
		log.Fatalf("Error reading directory configuration: %s", err.Error())
	}
	if err = json.Unmarshal(data, &directories); err != nil {
		log.Fatalf("Error reading directory configuration: %s", err.Error())
	}
	sort.Slice(directories, func(l, r int) bool {
		left := directories[l]
		right := directories[r]

		return left.Name < right.Name
	})

	t, _ := template.New("directory").Parse(templates.Directory)
	f, err := os.OpenFile(directoryFile+"~", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Error generating directory file: %s", err.Error())
	}
	defer f.Close()

	err = t.ExecuteTemplate(f, "main", struct {
		CodeGen     MetaInfo
		PackageName string
		Directories []Directory
	}{
		CodeGen:     options.MetaInfo,
		PackageName: options.PackageName,
		Directories: directories,
	})
	if err != nil {
		log.Fatalf("Error generating directory file: %s", err.Error())
	}
	err = os.Rename(directoryFile+"~", directoryFile)
	if err != nil {
		log.Fatalf("Error generating directory file: %s", err.Error())
	}

	goFmt(directoryFile)
}

// Directory describes a directory object
type Directory struct {
	Name           string      `json:"name"`
	DirectoryName  string      `json:"dir_name"`
	Required       bool        `json:"required"`
	SubDirectories []Directory `json:"subdirs"`
	IsData         bool        `json:"is_data"`
}
