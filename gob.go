package main

import (
	"fmt"
	"os"
	"path"
	"sort"
	"text/template"

	"github.com/ecnepsnai/gengo/templates"
)

type TGobGenerator struct{}

var GobGenerator = &TGobGenerator{}

func (g *TGobGenerator) Generate(options Options) (*GeneratorResult, error) {
	gobFileName := fmt.Sprintf("%sgob.go", options.FilePrefix)

	var gobs []Gob
	if !loadConfig(options.ConfigDir, "gob", &gobs) {
		return nil, nil
	}

	var imports = map[string]bool{}
	for _, gob := range gobs {
		if gob.Import != "" {
			imports[gob.Import] = true
		}
	}
	sort.Slice(gobs, func(l, r int) bool {
		left := gobs[l]
		right := gobs[r]

		return left.Type < right.Type
	})

	t, _ := template.New("gob").Parse(templates.GobGo)
	f, err := os.OpenFile(path.Join(options.TempDir, gobFileName), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating %s: %s", gobFileName, err.Error())
		return nil, err
	}
	defer f.Close()

	err = t.ExecuteTemplate(f, "main", struct {
		GenGo       MetaInfo
		PackageName string
		Gobs        []Gob
		Imports     []string
	}{
		GenGo:       options.MetaInfo,
		PackageName: options.PackageName,
		Gobs:        gobs,
		Imports:     mapKeys(imports),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating %s: %s", gobFileName, err.Error())
		return nil, err
	}

	return &GeneratorResult{
		GoFiles: []string{
			gobFileName,
		},
	}, nil
}

// Gob describes an gob type
type Gob struct {
	Type   string `json:"type" yaml:"type"`
	Import string `json:"import" yaml:"import"`
}
