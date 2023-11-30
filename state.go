package main

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
	"text/template"

	"github.com/ecnepsnai/gengo/templates"
)

type TStateGenerator struct{}

var StateGenerator = &TStateGenerator{}

func (g *TStateGenerator) Generate(options Options) (*GeneratorResult, error) {
	stateFileName := fmt.Sprintf("%sstate.go", options.FilePrefix)

	var states []StateProperty
	if !loadConfig(options.ConfigDir, "state", &states) {
		return nil, nil
	}

	t, _ := template.New("state").Parse(templates.StateGo)
	f, err := os.OpenFile(path.Join(options.TempDir, stateFileName), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating %s: %s", stateFileName, err.Error())
		return nil, err
	}
	defer f.Close()

	var tmap = map[string]bool{}
	var impmap = map[string]bool{}
	i := 0
	for i < len(states) {
		property := states[i]
		tmap[property.UnsafeType] = true
		states[i].Type = stateType{
			Name: strings.ReplaceAll(strings.ReplaceAll(property.UnsafeType, "[]", "Arr"), ".", ""),
			Type: property.UnsafeType,
		}

		for _, imp := range property.Import {
			impmap[imp] = true
		}

		i++
	}

	types := make([]stateType, len(tmap))
	i = 0
	for k := range tmap {
		types[i] = stateType{
			Type: k,
			Name: strings.ReplaceAll(strings.ReplaceAll(k, "[]", "Arr"), ".", ""),
		}
		i++
	}
	sort.Slice(types, func(l, r int) bool {
		left := types[l]
		right := types[r]
		return left.Name < right.Name
	})

	imports := make([]string, len(impmap))
	i = 0
	for imp := range impmap {
		imports[i] = imp
		i++
	}

	err = t.ExecuteTemplate(f, "main", struct {
		GenGo       MetaInfo
		PackageName string
		Properties  []StateProperty
		Types       []stateType
		Imports     []string
	}{
		GenGo:       options.MetaInfo,
		PackageName: options.PackageName,
		Properties:  states,
		Types:       types,
		Imports:     imports,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating %s: %s", stateFileName, err.Error())
		return nil, err
	}

	return &GeneratorResult{
		GoFiles: []string{
			stateFileName,
		},
	}, nil
}

// StateProperty describes a state property
type StateProperty struct {
	Name       string    `json:"name" yaml:"name"`
	UnsafeType string    `json:"type" yaml:"type"`
	Type       stateType `json:"-" yaml:"-"`
	Default    string    `json:"default" yaml:"default"`
	Import     []string  `json:"import" yaml:"import"`
}

type stateType struct {
	Type string
	Name string
}
