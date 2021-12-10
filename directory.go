package main

import (
	"log"
	"os"
	"sort"
	"text/template"

	"github.com/ecnepsnai/cbgen/templates"
)

const directoryFileName = "cbgen_directory.go"

// GenerateDirectory generates directory file
func GenerateDirectory(options Options) {
	var directories []Directory
	if !loadConfig("directory", &directories) {
		return
	}

	sort.Slice(directories, func(l, r int) bool {
		left := directories[l]
		right := directories[r]

		return left.Name < right.Name
	})

	t, _ := template.New("directory").Parse(templates.Directory)
	f, err := os.OpenFile(directoryFileName+"~", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
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
	err = os.Rename(directoryFileName+"~", directoryFileName)
	if err != nil {
		log.Fatalf("Error generating directory file: %s", err.Error())
	}

	goFmt(directoryFileName)
}

// Directory describes a directory object
type Directory struct {
	Name           string      `json:"name" yaml:"name"`
	DirectoryName  string      `json:"dir_name" yaml:"dir_name"`
	Required       bool        `json:"required" yaml:"required"`
	SubDirectories []Directory `json:"subdirs" yaml:"subdirs"`
	IsData         bool        `json:"is_data" yaml:"is_data"`
}
