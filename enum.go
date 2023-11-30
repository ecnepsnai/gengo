package main

import (
	"fmt"
	"os"
	"path"
	"sort"
	"text/template"

	"github.com/ecnepsnai/gengo/templates"
)

type TEnumGenerator struct{}

var EnumGenerator = &TEnumGenerator{}

func (g *TEnumGenerator) Generate(options Options) (*GeneratorResult, error) {
	enumGoFileName := fmt.Sprintf("%senum.go", options.FilePrefix)
	enumTsFileName := fmt.Sprintf("%senum.ts", options.FilePrefix)

	var enums []Enum
	if !loadConfig(options.ConfigDir, "enum", &enums) {
		return nil, nil
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
	f, err := os.OpenFile(path.Join(options.TempDir, enumGoFileName), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating %s: %s", enumGoFileName, err.Error())
		return nil, err
	}

	err = goTemplate.ExecuteTemplate(f, "main", struct {
		GenGo       MetaInfo
		PackageName string
		Enums       []Enum
	}{
		GenGo:       options.MetaInfo,
		PackageName: options.PackageName,
		Enums:       enums,
	})
	f.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating %s: %s", enumGoFileName, err.Error())
		return nil, err
	}

	if !exportTs {
		return &GeneratorResult{
			GoFiles: []string{
				enumGoFileName,
			},
		}, nil
	}

	tsTemplate, _ := template.New("enum").Parse(templates.EnumTs)
	f, err = os.OpenFile(path.Join(options.TempDir, enumTsFileName), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating %s: %s", enumTsFileName, err.Error())
		return nil, err
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
		fmt.Fprintf(os.Stderr, "Error generating %s: %s", enumTsFileName, err.Error())
		return nil, err
	}
	return &GeneratorResult{
		GoFiles: []string{
			enumGoFileName,
		},
		TsFiles: []string{
			enumTsFileName,
		},
	}, nil
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
