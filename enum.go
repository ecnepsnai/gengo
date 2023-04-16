package main

import (
	"log"
	"os"
	"sort"
	"text/template"

	"github.com/ecnepsnai/cbgen/templates"
)

const enumFileName = "cbgen_enum.go"

// GenerateEnum generate enums and schemas
func GenerateEnum(options Options) {
	var enums []Enum
	if !loadConfig("enum", &enums) {
		return
	}

	sort.Slice(enums, func(l, r int) bool {
		left := enums[l]
		right := enums[r]

		return left.Name < right.Name
	})

	t, _ := template.New("enum").Parse(templates.EnumGo)
	f, err := os.OpenFile(enumFileName+"~", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
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
	err = os.Rename(enumFileName+"~", enumFileName)
	if err != nil {
		log.Fatalf("Error generating enum file: %s", err.Error())
	}

	goFmt(enumFileName)
}

// Enum describes an enum type
type Enum struct {
	Name   string      `json:"name" yaml:"name"`
	Type   string      `json:"type" yaml:"type"`
	Values []EnumValue `json:"values" yaml:"values"`
}

// EnumValue describes an single enum value
type EnumValue struct {
	Key         string `json:"key" yaml:"key"`
	Description string `json:"description" yaml:"description"`
	Value       string `json:"value" yaml:"value"`
	Name        string `json:"name" yaml:"name"`
}

// Title return the title of the enum value
func (v EnumValue) Title() string {
	if v.Name != "" {
		return v.Name
	}
	return v.Key
}
