package main

import (
	"log"
	"os"
	"strings"
	"text/template"
)

// GenerateVersion generate a version file
func GenerateVersion(options Options) {
	if options.PackageVersion == "" {
		return
	}

	versionFile := "cbgen_version.go"

	t := template.Must(template.ParseFiles(getTemplateFile("version.tmpl")))
	f, err := os.OpenFile(versionFile+"~", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Error generating version file: %s", err.Error())
	}
	defer f.Close()

	err = t.ExecuteTemplate(f, "main", struct {
		CodeGen     MetaInfo
		PackageName string
		AppName     string
		AppVersion  string
	}{
		CodeGen:     options.MetaInfo,
		PackageName: options.PackageName,
		AppName:     strings.Title(options.PackageName),
		AppVersion:  options.PackageVersion,
	})
	if err != nil {
		log.Fatalf("Error generating version file: %s", err.Error())
	}
	err = os.Rename(versionFile+"~", versionFile)
	if err != nil {
		log.Fatalf("Error generating version file: %s", err.Error())
	}

	goFmt(versionFile)
}
