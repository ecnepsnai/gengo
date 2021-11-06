package main

import (
	"encoding/json"
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"sync"
)

// Version the current version of Codegen
var Version = "v1.8.1"

func main() {
	if len(os.Args) <= 1 {
		printHelpAndExit()
	}

	args := os.Args[1:]

	var packageName string
	var packageVersion string

	i := 0
	for i < len(args) {
		arg := args[i]
		if arg == "-n" || arg == "--name" {
			packageName = args[i+1]
			i++
		} else if arg == "-v" || arg == "--version" {
			packageVersion = args[i+1]
			i++
		} else {
			fmt.Printf("Unknown argument '%s'\n", arg)
			printHelpAndExit()
		}
		i++
	}

	assertCodegenVersion()

	meta := MetaInfo{
		Version: Version,
	}

	Generate(Options{
		PackageName:    packageName,
		PackageVersion: packageVersion,
		MetaInfo:       meta,
	})
}

func printHelpAndExit() {
	fmt.Printf("Usage: %s -n <package name> [-v <package version]\n", os.Args[0])
	fmt.Printf("-n --name\tPackage name\n")
	fmt.Printf("-v --version\tPackage version. Including will generate a version go file\n")
	os.Exit(1)
}

// Options describes CodeGen options
type Options struct {
	PackageName    string
	PackageVersion string
	MetaInfo       MetaInfo
}

var gopath string

// Generate run the generator with the given options. Outputs files in the current working directory
func Generate(options Options) {
	gopath = os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	type genFunc func(Options)
	ops := []genFunc{
		GenerateVersion,
		GenerateDirectory,
		GenerateState,
		GenerateStore,
		GenerateEnum,
		GenerateStats,
		GenerateDataStore,
		GenerateGob,
	}

	wg := sync.WaitGroup{}
	wg.Add(len(ops))

	for i := range ops {
		op := ops[i]
		go func() {
			op(options)
			wg.Done()
		}()
	}

	wg.Wait()
}

func getTemplateFile(templateName string) string {
	gopath := os.Getenv("GOPATH")
	return path.Join(gopath, "src", "github.com", "ecnepsnai", "cbgen", "templates", templateName)
}

func goFmt(path string) {
	exec.Command("go", "fmt", path).Run()
}

// MetaInfo describes meta information about CodeGen
type MetaInfo struct {
	Version string
}

func mapKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func assertCodegenVersion() {
	codegenConfigPath := "codegen.json"
	if _, err := os.Stat(codegenConfigPath); err != nil {
		return
	}

	codegenConfig := struct {
		MinimumVersion string `json:"minimum_version"`
	}{}
	f, err := os.OpenFile(codegenConfigPath, os.O_RDONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&codegenConfig); err != nil {
		return
	}

	versionStrToNumber := func(in string) int {
		v := strings.ReplaceAll(in[1:], ".", "")
		i, err := strconv.Atoi(v)
		if err != nil {
			return -1
		}
		return i
	}

	currentVersionNumber := versionStrToNumber(Version)
	minimumVersionNumber := versionStrToNumber(codegenConfig.MinimumVersion)

	if minimumVersionNumber > currentVersionNumber {
		fmt.Fprintf(os.Stderr, "Incorrect Codegen version installed.\nWanted: %s\nInstalled: %s\n", codegenConfig.MinimumVersion, Version)
		os.Exit(1)
	}
}
