package main

import (
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"text/template"

	"github.com/ecnepsnai/cbgen/templates"
)

const stateFileName = "cbgen_state.go"

// GenerateState generates the state store
func GenerateState(options Options) {
	var states []StateProperty
	if !loadConfig(options.ConfigDir, "state", &states) {
		return
	}

	t, _ := template.New("state").Parse(templates.StateGo)
	f, err := os.OpenFile(path.Join(options.OutputDir, stateFileName+"~"), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Error generating state file: %s", err.Error())
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
		CodeGen     MetaInfo
		PackageName string
		Properties  []StateProperty
		Types       []stateType
		Imports     []string
	}{
		CodeGen:     options.MetaInfo,
		PackageName: options.PackageName,
		Properties:  states,
		Types:       types,
		Imports:     imports,
	})
	if err != nil {
		log.Fatalf("Error generating state file: %s", err.Error())
	}
	err = os.Rename(path.Join(options.OutputDir, stateFileName+"~"), path.Join(options.OutputDir, stateFileName))
	if err != nil {
		log.Fatalf("Error generating state file: %s", err.Error())
	}

	goFmt(path.Join(options.OutputDir, stateFileName))
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
