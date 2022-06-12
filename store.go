package main

import (
	"log"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/ecnepsnai/cbgen/templates"
)

const storeFileName = "cbgen_store.go"

// GenerateStore generates the store file
func GenerateStore(options Options) {
	var stores []Store
	if !loadConfig("store", &stores) {
		return
	}

	t, _ := template.New("store").Parse(templates.Store)
	f, err := os.OpenFile(storeFileName+"~", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
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
		//lint:ignore SA1019 Puncuation not supported anyways
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
	err = os.Rename(storeFileName+"~", storeFileName)
	if err != nil {
		log.Fatalf("Error generating store file: %s", err.Error())
	}

	goFmt(storeFileName)
}

// Store describes a store type
type Store struct {
	Name          string `json:"name" yaml:"name"`
	LowercaseName string
	TitlecaseName string
	Interfaces    []string   `json:"gobs" yaml:"gobs"`
	Gobs          []StoreGob `json:"-" yaml:"-"`
	ExtraImports  []string   `json:"extra_imports" yaml:"extra_imports"`
}

// StoreGob describes a object to encode/decode using gob
type StoreGob struct {
	Name string
	Type string
}
