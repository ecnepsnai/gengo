package main

import (
	"log"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/ecnepsnai/cbgen/templates"
)

const dataStoreFileName = "cbgen_data_store.go"

// GenerateDataStore generates the data store file
func GenerateDataStore(options Options) {
	var stores []DataStore
	if !loadConfig("data_store", &stores) {
		return
	}

	t, _ := template.New("data_store").Parse(templates.DataStore)
	f, err := os.OpenFile(dataStoreFileName+"~", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Error generating data store file: %s", err.Error())
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
		CodeGen      MetaInfo
		PackageName  string
		Stores       []DataStore
		ExtraImports []string
	}{
		CodeGen:      options.MetaInfo,
		PackageName:  options.PackageName,
		Stores:       stores,
		ExtraImports: extraImports,
	})
	if err != nil {
		log.Fatalf("Error generating data store file: %s", err.Error())
	}
	err = os.Rename(dataStoreFileName+"~", dataStoreFileName)
	if err != nil {
		log.Fatalf("Error generating data store file: %s", err.Error())
	}

	goFmt(dataStoreFileName)
}

// DataStore describes a data store type
type DataStore struct {
	Name          string `json:"name" yaml:"name"`
	LowercaseName string
	TitlecaseName string
	Object        string `json:"object" yaml:"object"`
	Unordered     bool   `json:"unordered" yaml:"unordered"`
}
