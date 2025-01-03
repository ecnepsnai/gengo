package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"sync"
)

// Version the current version of GenGo
var Version = "v1.14.1"

func main() {
	if len(os.Args) <= 1 {
		printHelpAndExit()
	}

	args := os.Args[1:]

	packageName := "main"
	configDir := "."
	goOutputDir := "."
	tsOutputDir := "."
	quiet := false

	i := 0
	for i < len(args) {
		arg := args[i]
		if arg == "-n" || arg == "--name" {
			packageName = args[i+1]
			i++
		} else if arg == "-c" || arg == "--config-dir" {
			configDir = args[i+1]
			i++
		} else if arg == "-g" || arg == "--go-output-dir" {
			goOutputDir = args[i+1]
			i++
		} else if arg == "-t" || arg == "--ts-output-dir" {
			tsOutputDir = args[i+1]
			i++
		} else if arg == "-q" || arg == "--quiet" {
			quiet = true
		} else {
			fmt.Printf("Unknown argument '%s'\n", arg)
			printHelpAndExit()
		}
		i++
	}

	loadGenGoConfig(configDir)

	meta := MetaInfo{
		Version: Version,
	}

	tempDir, err := os.MkdirTemp(configDir, "gengo")
	if err != nil {
		panic(err)
	}

	Generate(Options{
		PackageName: packageName,
		ConfigDir:   configDir,
		TempDir:     tempDir,
		GoOutputDir: goOutputDir,
		TsOutputDir: tsOutputDir,
		FilePrefix:  GenGoConfig.FilePrefix,
		Quiet:       quiet,
		MetaInfo:    meta,
	})

	os.RemoveAll(tempDir)
}

func printHelpAndExit() {
	fmt.Printf("Usage: %s [Options]\n", os.Args[0])
	fmt.Printf("-n --name <name>           Package name, defaults to 'main'\n")
	fmt.Printf("-c --config-dir <dir>      Config dir, defaults to current dir\n")
	fmt.Printf("-g --go-output-dir <dir>   Output dir for go files, defaults to current dir\n")
	fmt.Printf("-t --ts-output-dir <dir>   Output dir for ts files, defaults to current dir\n")
	fmt.Printf("-q --quiet                 Don't print out names of generated files\n")
	os.Exit(1)
}

// Options describes GenGo options
type Options struct {
	PackageName string
	ConfigDir   string
	TempDir     string
	GoOutputDir string
	TsOutputDir string
	FilePrefix  string
	Quiet       bool
	MetaInfo    MetaInfo
}

// Generate run the generator with the given options. Outputs files in the current working directory
func Generate(options Options) {
	generators := []IGenerator{
		DataStoreGenerator,
		DictionaryGenerator,
		EnumGenerator,
		GobGenerator,
		StateGenerator,
		StatsGenerator,
		StoreGenerator,
	}

	wg := sync.WaitGroup{}
	wg.Add(len(generators))

	success := true

	for i := range generators {
		generator := generators[i]
		go func() {
			defer wg.Done()

			result, err := generator.Generate(options)
			if err != nil {
				success = false
				return
			}
			if result == nil {
				return
			}

			for _, goFile := range result.GoFiles {
				if err := goFmt(path.Join(options.TempDir, goFile)); err != nil {
					success = false
					return
				}

				if err := os.Rename(path.Join(options.TempDir, goFile), path.Join(options.GoOutputDir, goFile)); err != nil {
					success = false
					fmt.Fprintf(os.Stderr, "Error generating %s: %s\n", goFile, err.Error())
					return
				}

				if !options.Quiet {
					fmt.Println(goFile)
				}
			}

			for _, tsFile := range result.TsFiles {
				if err := os.Rename(path.Join(options.TempDir, tsFile), path.Join(options.TsOutputDir, tsFile)); err != nil {
					success = false
					fmt.Fprintf(os.Stderr, "Error generating %s: %s\n", tsFile, err.Error())
					return
				}

				if !options.Quiet {
					fmt.Println(tsFile)
				}
			}
		}()
	}

	wg.Wait()

	if !success {
		os.Exit(1)
	}
}

func goFmt(path string) error {
	output, err := exec.Command("gofmt", "-l", "-w", path).CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", output)
		return fmt.Errorf("gofmt error")
	}
	return nil
}

// MetaInfo describes meta information about GenGo
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
