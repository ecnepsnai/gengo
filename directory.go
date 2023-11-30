package main

import (
	"fmt"
	"os"
	"path"
	"sort"
	"text/template"

	"github.com/ecnepsnai/gengo/templates"
)

type TDictionaryGenerator struct{}

var DictionaryGenerator = &TDictionaryGenerator{}

func (g *TDictionaryGenerator) Generate(options Options) (*GeneratorResult, error) {
	directoryFileName := fmt.Sprintf("%sdirectory.go", options.FilePrefix)

	var directories []Directory
	if !loadConfig(options.ConfigDir, "directory", &directories) {
		return nil, nil
	}

	sort.Slice(directories, func(l, r int) bool {
		left := directories[l]
		right := directories[r]

		return left.Name < right.Name
	})

	t, _ := template.New("directory").Parse(templates.DirectoryGo)
	f, err := os.OpenFile(path.Join(options.TempDir, directoryFileName), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating %s: %s", directoryFileName, err.Error())
		return nil, err
	}
	defer f.Close()

	err = t.ExecuteTemplate(f, "main", struct {
		GenGo       MetaInfo
		PackageName string
		Directories []Directory
	}{
		GenGo:       options.MetaInfo,
		PackageName: options.PackageName,
		Directories: directories,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating %s: %s", directoryFileName, err.Error())
		return nil, err
	}

	return &GeneratorResult{
		GoFiles: []string{
			directoryFileName,
		},
	}, nil
}

// Directory describes a directory object
type Directory struct {
	Name           string      `json:"name" yaml:"name"`
	DirectoryName  string      `json:"dir_name" yaml:"dir_name"`
	Required       bool        `json:"required" yaml:"required"`
	SubDirectories []Directory `json:"subdirs" yaml:"subdirs"`
	IsData         bool        `json:"is_data" yaml:"is_data"`
}
