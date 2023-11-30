package main

import (
	"fmt"
	"os"
	"path"
	"sort"
	"text/template"

	"github.com/ecnepsnai/gengo/templates"
)

type TStatsGenerator struct{}

var StatsGenerator = &TStatsGenerator{}

func (g *TStatsGenerator) Generate(options Options) (*GeneratorResult, error) {
	statsFileName := fmt.Sprintf("%sstats.go", options.FilePrefix)

	var stats Stats
	if !loadConfig(options.ConfigDir, "stats", &stats) {
		return nil, nil
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
	f, err := os.OpenFile(path.Join(options.TempDir, statsFileName), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating %s: %s", statsFileName, err.Error())
		return nil, err
	}
	defer f.Close()

	err = t.ExecuteTemplate(f, "main", struct {
		GenGo       MetaInfo
		PackageName string
		Stats       Stats
	}{
		GenGo:       options.MetaInfo,
		PackageName: options.PackageName,
		Stats:       stats,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating %s: %s", statsFileName, err.Error())
		return nil, err
	}

	return &GeneratorResult{
		GoFiles: []string{
			statsFileName,
		},
	}, nil
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
