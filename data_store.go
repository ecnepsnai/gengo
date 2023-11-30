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

type TDataStoreGenerator struct{}

var DataStoreGenerator = &TDataStoreGenerator{}

func (g *TDataStoreGenerator) Generate(options Options) (*GeneratorResult, error) {
	dataStoreFileName := fmt.Sprintf("%sdata_store.go", options.FilePrefix)

	var stores []DataStore
	if !loadConfig(options.ConfigDir, "data_store", &stores) {
		return nil, nil
	}

	t, _ := template.New("data_store").Parse(templates.DataStoreGo)
	f, err := os.OpenFile(path.Join(options.TempDir, dataStoreFileName), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating %s: %s", dataStoreFileName, err.Error())
		return nil, err
	}
	defer f.Close()

	var extraImports []string
	for i, store := range stores {
		store.LowercaseName = strings.ToLower(store.Name)
		store.TitlecaseName = strings.Title(store.Name)
		stores[i] = store
	}
	sort.Slice(stores, func(l, r int) bool {
		left := stores[l]
		right := stores[r]

		return left.Name < right.Name
	})

	err = t.ExecuteTemplate(f, "main", struct {
		GenGo        MetaInfo
		PackageName  string
		Stores       []DataStore
		ExtraImports []string
	}{
		GenGo:        options.MetaInfo,
		PackageName:  options.PackageName,
		Stores:       stores,
		ExtraImports: extraImports,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating %s: %s", dataStoreFileName, err.Error())
		return nil, err
	}

	return &GeneratorResult{
		GoFiles: []string{
			dataStoreFileName,
		},
	}, nil
}

// DataStore describes a data store type
type DataStore struct {
	Name          string `json:"name" yaml:"name"`
	LowercaseName string
	TitlecaseName string
	Object        string `json:"object" yaml:"object"`
	Unordered     bool   `json:"unordered" yaml:"unordered"`
}
