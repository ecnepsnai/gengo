package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"text/template"
)

// GenerateDataStore generates the data store file
func GenerateDataStore(options Options) {
	storeConfig := path.Join(".", "data_store.json")
	storeFile := path.Join(".", "cbgen_data_store.go")

	if _, err := os.Stat(storeConfig); err != nil {
		return
	}

	var stores []DataStore
	data, err := ioutil.ReadFile(storeConfig)
	if err != nil {
		log.Fatalf("Error reading data store configuration: %s", err.Error())
	}
	if err = json.Unmarshal(data, &stores); err != nil {
		log.Fatalf("Error reading data store configuration: %s", err.Error())
	}

	t := template.Must(template.ParseFiles(getTemplateFile("data_store.tmpl")))
	f, err := os.OpenFile(storeFile+"~", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
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
	err = os.Rename(storeFile+"~", storeFile)
	if err != nil {
		log.Fatalf("Error generating data store file: %s", err.Error())
	}

	goFmt(storeFile)
}

// DataStore describes a data store type
type DataStore struct {
	Name          string `json:"name"`
	LowercaseName string
	TitlecaseName string
	Object        string `json:"object"`
	Unordered     bool   `json:"unordered"`
}
