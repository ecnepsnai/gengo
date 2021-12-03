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

	"github.com/ecnepsnai/cbgen/templates"
)

// GenerateStore generates the store file
func GenerateStore(options Options) {
	storeConfig := path.Join(".", "store.json")
	storeFile := path.Join(".", "cbgen_store.go")

	if _, err := os.Stat(storeConfig); err != nil {
		return
	}

	var stores []Store
	data, err := ioutil.ReadFile(storeConfig)
	if err != nil {
		log.Fatalf("Error reading store configuration: %s", err.Error())
	}
	if err = json.Unmarshal(data, &stores); err != nil {
		log.Fatalf("Error reading store configuration: %s", err.Error())
	}

	t, _ := template.New("store").Parse(templates.Store)
	f, err := os.OpenFile(storeFile+"~", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Error generating store file: %s", err.Error())
	}
	defer f.Close()

	for i, store := range stores {
		var gobs []StoreGob
		for _, intf := range store.Interfaces {
			var name string
			if strings.Contains(intf, ".") {
				components := strings.Split(intf, ".")
				name = components[len(components)-1]
			} else {
				name = intf
			}

			name = strings.Replace(name, "{}", "", -1)
			object := strings.Replace(intf, "{}", "", -1)

			gobs = append(gobs, StoreGob{
				Name: name,
				Type: object,
			})
		}
		stores[i].Gobs = gobs
	}

	var extraImports []string
	for i, store := range stores {
		store.LowercaseName = strings.ToLower(store.Name)
		store.TitlecaseName = strings.Title(store.Name)
		stores[i] = store

		if store.ExtraImports != nil {
			extraImports = append(extraImports, store.ExtraImports...)
		}
	}
	sort.Slice(stores, func(l, r int) bool {
		left := stores[l]
		right := stores[r]
		return left.Name < right.Name
	})

	err = t.ExecuteTemplate(f, "main", struct {
		CodeGen      MetaInfo
		PackageName  string
		Stores       []Store
		ExtraImports []string
	}{
		CodeGen:      options.MetaInfo,
		PackageName:  options.PackageName,
		Stores:       stores,
		ExtraImports: extraImports,
	})
	if err != nil {
		log.Fatalf("Error generating store file: %s", err.Error())
		defer os.Remove(f.Name())
	}
	err = os.Rename(storeFile+"~", storeFile)
	if err != nil {
		log.Fatalf("Error generating store file: %s", err.Error())
	}

	goFmt(storeFile)
}

// Store describes a store type
type Store struct {
	Name          string `json:"name"`
	LowercaseName string
	TitlecaseName string
	Interfaces    []string   `json:"gobs"`
	Gobs          []StoreGob `json:"-"`
	ExtraImports  []string   `json:"extra_imports"`
}

// StoreGob describes a object to encode/decode using gob
type StoreGob struct {
	Name string
	Type string
}
