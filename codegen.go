package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

// Version the current version of Codegen
var Version = "v1.10.0"

func main() {
	if len(os.Args) <= 1 {
		printHelpAndExit()
	}

	args := os.Args[1:]

	var packageName string

	i := 0
	for i < len(args) {
		arg := args[i]
		if arg == "-n" || arg == "--name" {
			packageName = args[i+1]
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
		PackageName: packageName,
		MetaInfo:    meta,
	})
}

func printHelpAndExit() {
	fmt.Printf("Usage: %s -n <package name>\n", os.Args[0])
	fmt.Printf("-n --name\tPackage name\n")
	os.Exit(1)
}

// Options describes CodeGen options
type Options struct {
	PackageName string
	MetaInfo    MetaInfo
}

// Generate run the generator with the given options. Outputs files in the current working directory
func Generate(options Options) {
	type genFunc func(Options)
	ops := []genFunc{
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
		MinimumVersion string `json:"minimum_version" yaml:"minimum_version"`
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
