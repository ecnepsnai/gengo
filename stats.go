package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"text/template"
)

// GenerateStats generates the stats file
func GenerateStats(options Options) {
	statsConfig := path.Join(".", "stats.json")
	statsFile := path.Join(".", "cbgen_stats.go")

	if _, err := os.Stat(statsConfig); err != nil {
		return
	}

	var stats Stats
	data, err := ioutil.ReadFile(statsConfig)
	if err != nil {
		log.Fatalf("Error reading stats configuration: %s", err.Error())
	}
	if err = json.Unmarshal(data, &stats); err != nil {
		log.Fatalf("Error reading stats configuration: %s", err.Error())
	}

	t := template.Must(template.ParseFiles(getTemplateFile("stats.tmpl")))
	f, err := os.OpenFile(statsFile+"~", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Error generating stats file: %s", err.Error())
	}
	defer f.Close()

	err = t.ExecuteTemplate(f, "main", struct {
		CodeGen     MetaInfo
		PackageName string
		Stats       Stats
	}{
		CodeGen:     options.MetaInfo,
		PackageName: options.PackageName,
		Stats:       stats,
	})
	if err != nil {
		log.Fatalf("Error generating stats file: %s", err.Error())
	}
	err = os.Rename(statsFile+"~", statsFile)
	if err != nil {
		log.Fatalf("Error generating stats file: %s", err.Error())
	}

	goFmt(statsFile)
}

// Counter describes a Counter object
type Counter struct {
	Name        string
	Description string
}

// Timer describes a Timer object
type Timer struct {
	Name        string
	Description string
}

// Stats describes a Stats object
type Stats struct {
	NonvolatileCounters []Counter
	VolatileCounters    []Counter
	Timers              []Timer
}
