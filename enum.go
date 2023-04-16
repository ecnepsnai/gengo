package main

import (
	"log"
	"os"
	"sort"
	"text/template"

	"github.com/ecnepsnai/cbgen/templates"
)

const enumGoFileName = "cbgen_enum.go"
const enumTsFileName = "cbgen_enum.ts"

// GenerateEnum generate enums and schemas
func GenerateEnum(options Options) {
	var enums []Enum
	if !loadConfig("enum", &enums) {
		return
	}

	exportTs := false
	sort.Slice(enums, func(l, r int) bool {
		left := enums[l]
		right := enums[r]

		if left.IncludeTypeScript || right.IncludeTypeScript {
			exportTs = true
		}

		return left.Name < right.Name
	})

	goTemplate, _ := template.New("enum").Parse(templates.EnumGo)
	f, err := os.OpenFile(enumGoFileName+"~", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Error generating enum go file: %s", err.Error())
	}

	err = goTemplate.ExecuteTemplate(f, "main", struct {
		CodeGen     MetaInfo
		PackageName string
		Enums       []Enum
	}{
		CodeGen:     options.MetaInfo,
		PackageName: options.PackageName,
		Enums:       enums,
	})
	f.Close()
	if err != nil {
		log.Fatalf("Error generating enum go file: %s", err.Error())
	}
	err = os.Rename(enumGoFileName+"~", enumGoFileName)
	if err != nil {
		log.Fatalf("Error generating enum go file: %s", err.Error())
	}

	goFmt(enumGoFileName)

	if !exportTs {
		return
	}

	tsTemplate, _ := template.New("enum").Parse(templates.EnumTs)
	f, err = os.OpenFile(enumTsFileName+"~", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Error generating enum ts file: %s", err.Error())
	}

	tsEnums := []Enum{}
	for _, enum := range enums {
		if !enum.IncludeTypeScript {
			continue
		}
		values := enum.Values
		if enum.Type == "string" {
			for i, v := range values {
				values[i].Value = "'" + v.Value[1:len(v.Value)-1] + "'"
			}
		}
		tsEnums = append(tsEnums, Enum{
			Name:        enum.Name,
			Type:        enum.Type,
			Description: enum.Description,
			Values:      values,
		})
	}

	err = tsTemplate.ExecuteTemplate(f, "main", struct {
		Version string
		TsEnums []Enum
	}{
		Version: options.MetaInfo.Version,
		TsEnums: tsEnums,
	})
	f.Close()
	if err != nil {
		log.Fatalf("Error generating enum ts file: %s", err.Error())
	}
	err = os.Rename(enumTsFileName+"~", enumTsFileName)
	if err != nil {
		log.Fatalf("Error generating enum ts file: %s", err.Error())
	}
}

// Enum describes an enum type
type Enum struct {
	Name              string      `json:"name" yaml:"name"`
	Type              string      `json:"type" yaml:"type"`
	Description       string      `json:"description" yaml:"description"`
	Values            []EnumValue `json:"values" yaml:"values"`
	IncludeTypeScript bool        `json:"include_typescript" yaml:"include_typescript"`
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
