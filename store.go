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

type TStoreGenerator struct{}

var StoreGenerator = &TStoreGenerator{}

func (g *TStoreGenerator) Generate(options Options) (*GeneratorResult, error) {
	storeFileName := fmt.Sprintf("%sstore.go", options.FilePrefix)

	var stores []Store
	if !loadConfig(options.ConfigDir, "store", &stores) {
		return nil, nil
	}

	t, _ := template.New("store").Parse(templates.StoreGo)
	f, err := os.OpenFile(path.Join(options.TempDir, storeFileName), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating %s: %s", storeFileName, err.Error())
		return nil, err
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
		GenGo        MetaInfo
		PackageName  string
		Stores       []Store
		ExtraImports []string
	}{
		GenGo:        options.MetaInfo,
		PackageName:  options.PackageName,
		Stores:       stores,
		ExtraImports: extraImports,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating %s: %s", storeFileName, err.Error())
		return nil, err
	}

	return &GeneratorResult{
		GoFiles: []string{
			storeFileName,
		},
	}, nil
}

// Store describes a store type
type Store struct {
	Name          string `json:"name" yaml:"name"`
	BucketName    string `json:"bucket_name" yaml:"bucket_name"`
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
