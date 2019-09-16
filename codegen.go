package main

import (
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path"
	"sync"
	"time"
)

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

	meta := MetaInfo{
		Date:    time.Now().Format("2006-01-02"),
		Version: "1.2.0",
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

	wg := sync.WaitGroup{}
	wg.Add(8)

	go func() {
		GenerateVersion(options)
		wg.Done()
	}()
	go func() {
		GenerateDirectory(options)
		wg.Done()
	}()
	go func() {
		GenerateState(options)
		wg.Done()
	}()
	go func() {
		GenerateStore(options)
		wg.Done()
	}()
	go func() {
		GenerateEnum(options)
		wg.Done()
	}()
	go func() {
		GenerateStats(options)
		wg.Done()
	}()
	go func() {
		GenerateDataStore(options)
		wg.Done()
	}()
	go func() {
		GenerateGob(options)
		wg.Done()
	}()

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
	Date    string
	Version string
}

func mapKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
