package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"text/template"
)

// GenerateEnum generate enums and schemas
func GenerateEnum(options Options) {
	enumConfig := path.Join(".", "enum.json")
	enumFile := path.Join(".", "cbgen_enum.go")

	if _, err := os.Stat(enumConfig); err != nil {
		return
	}

	var enums []Enum
	data, err := ioutil.ReadFile(enumConfig)
	if err != nil {
		log.Fatalf("Error reading enum configuration: %s", err.Error())
	}
	if err = json.Unmarshal(data, &enums); err != nil {
		log.Fatalf("Error reading enum configuration: %s", err.Error())
	}
	sort.Slice(enums, func(l, r int) bool {
		left := enums[l]
		right := enums[r]

		return left.Name < right.Name
	})

	t := template.Must(template.ParseFiles(getTemplateFile("enum.tmpl")))
	f, err := os.OpenFile(enumFile+"~", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Error generating enum file: %s", err.Error())
	}
	defer f.Close()

	err = t.ExecuteTemplate(f, "main", struct {
		CodeGen     MetaInfo
		PackageName string
		Enums       []Enum
	}{
		CodeGen:     options.MetaInfo,
		PackageName: options.PackageName,
		Enums:       enums,
	})
	if err != nil {
		log.Fatalf("Error generating enum file: %s", err.Error())
	}
	err = os.Rename(enumFile+"~", enumFile)
	if err != nil {
		log.Fatalf("Error generating enum file: %s", err.Error())
	}

	goFmt(enumFile)
}

// Enum describes an enum type
type Enum struct {
	Name   string      `json:"name"`
	Type   string      `json:"type"`
	Values []EnumValue `json:"values"`
}

// EnumValue describes an single enum value
type EnumValue struct {
	Key         string `json:"key"`
	Description string `json:"description"`
	Value       string `json:"value"`
	Name        string `json:"name"`
}

// Title return the title of the enum value
func (v EnumValue) Title() string {
	if v.Name != "" {
		return v.Name
	}
	return v.Key
}
