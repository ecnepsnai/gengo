package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
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
	sort.Slice(stats.Counters, func(l, r int) bool {
		left := stats.Counters[l]
		right := stats.Counters[r]

		return left.Name < right.Name
	})
	sort.Slice(stats.Timers, func(l, r int) bool {
		left := stats.Timers[l]
		right := stats.Timers[r]

		return left.Name < right.Name
	})

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
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Timer describes a Timer object
type Timer struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Stats describes a Stats object
type Stats struct {
	Counters []Counter `json:"counters"`
	Timers   []Timer   `json:"timers"`
}
