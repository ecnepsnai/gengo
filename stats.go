package main

import (
	"log"
	"os"
	"path"
	"sort"
	"text/template"

	"github.com/ecnepsnai/cbgen/templates"
)

const statsFileName = "cbgen_stats.go"

// GenerateStats generates the stats file
func GenerateStats(options Options) {
	var stats Stats
	if !loadConfig(options.ConfigDir, "stats", &stats) {
		return
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

	t, _ := template.New("stats").Parse(templates.StatsGo)
	f, err := os.OpenFile(path.Join(options.OutputDir, statsFileName+"~"), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
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
	err = os.Rename(path.Join(options.OutputDir, statsFileName+"~"), path.Join(options.OutputDir, statsFileName))
	if err != nil {
		log.Fatalf("Error generating stats file: %s", err.Error())
	}

	goFmt(path.Join(options.OutputDir, statsFileName))
}

// Counter describes a Counter object
type Counter struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
}

// TimedCounter describes a TimedCounter object
type TimedCounter struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	MaxMinutes  int    `json:"max_minutes" yaml:"max_minutes"`
}

// Timer describes a Timer object
type Timer struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
}

// Stats describes a Stats object
type Stats struct {
	Counters      []Counter      `json:"counters" yaml:"counters"`
	TimedCounters []TimedCounter `json:"timed_counters" yaml:"timed_counters"`
	Timers        []Timer        `json:"timers" yaml:"timers"`
}
