package main

import (
	"log"
	"os"
	"sort"
	"text/template"

	"github.com/ecnepsnai/cbgen/templates"
)

const gobFileName = "cbgen_gob.go"

// GenerateGob generate gob
func GenerateGob(options Options) {
	var gobs []Gob
	if !loadConfig("gob", &gobs) {
		return
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

	t, _ := template.New("gob").Parse(templates.Gob)
	f, err := os.OpenFile(gobFileName+"~", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Error generating gob file: %s", err.Error())
	}
	defer f.Close()

	err = t.ExecuteTemplate(f, "main", struct {
		CodeGen     MetaInfo
		PackageName string
		Gobs        []Gob
		Imports     []string
	}{
		CodeGen:     options.MetaInfo,
		PackageName: options.PackageName,
		Gobs:        gobs,
		Imports:     mapKeys(imports),
	})
	if err != nil {
		log.Fatalf("Error generating gob file: %s", err.Error())
	}
	err = os.Rename(gobFileName+"~", gobFileName)
	if err != nil {
		log.Fatalf("Error generating gob file: %s", err.Error())
	}

	goFmt(gobFileName)
}

// Gob describes an gob type
type Gob struct {
	Type   string `json:"type" yaml:"type"`
	Import string `json:"import" yaml:"import"`
}
