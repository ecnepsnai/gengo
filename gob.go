package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"text/template"
)

// GenerateGob generate gob
func GenerateGob(options Options) {
	gobConfig := path.Join(".", "gob.json")
	gobFile := path.Join(".", "cbgen_gob.go")

	if _, err := os.Stat(gobConfig); err != nil {
		return
	}

	var gobs []Gob
	data, err := ioutil.ReadFile(gobConfig)
	if err != nil {
		log.Fatalf("Error reading gob configuration: %s", err.Error())
	}
	if err = json.Unmarshal(data, &gobs); err != nil {
		log.Fatalf("Error reading gob configuration: %s", err.Error())
	}
	var imports = map[string]bool{}
	for _, gob := range gobs {
		if gob.Import != "" {
			imports[gob.Import] = true
		}
	}

	t := template.Must(template.ParseFiles(getTemplateFile("gob.tmpl")))
	f, err := os.OpenFile(gobFile+"~", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
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
	err = os.Rename(gobFile+"~", gobFile)
	if err != nil {
		log.Fatalf("Error generating gob file: %s", err.Error())
	}

	goFmt(gobFile)
}

// Gob describes an gob type
type Gob struct {
	Type   string `json:"type"`
	Import string `json:"import"`
}
