package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"text/template"
)

// GenerateState generates the state store
func GenerateState(options Options) {
	stateConfig := path.Join(".", "state.json")
	stateFile := path.Join(".", "cbgen_state.go")

	if _, err := os.Stat(stateConfig); err != nil {
		return
	}

	var states []StateProperty
	data, err := ioutil.ReadFile(stateConfig)
	if err != nil {
		log.Fatalf("Error reading state configuration: %s", err.Error())
	}
	if err = json.Unmarshal(data, &states); err != nil {
		log.Fatalf("Error reading state configuration: %s", err.Error())
	}

	t := template.Must(template.ParseFiles(getTemplateFile("state.tmpl")))
	f, err := os.OpenFile(stateFile+"~", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Error generating state file: %s", err.Error())
	}
	defer f.Close()

	var tmap = map[string]bool{}
	for _, property := range states {
		tmap[property.Type] = true
	}

	types := make([]string, len(tmap))
	i := 0
	for k := range tmap {
		types[i] = k
		i++
	}

	err = t.ExecuteTemplate(f, "main", struct {
		CodeGen     MetaInfo
		PackageName string
		Properties  []StateProperty
		Types       []string
	}{
		CodeGen:     options.MetaInfo,
		PackageName: options.PackageName,
		Properties:  states,
		Types:       types,
	})
	if err != nil {
		log.Fatalf("Error generating state file: %s", err.Error())
	}
	err = os.Rename(stateFile+"~", stateFile)
	if err != nil {
		log.Fatalf("Error generating state file: %s", err.Error())
	}

	goFmt(stateFile)
}

// StateProperty describes a state property
type StateProperty struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Default string `json:"default"`
}
