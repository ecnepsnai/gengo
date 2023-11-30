package main

type GeneratorResult struct {
	GoFiles []string
	TsFiles []string
}

type IGenerator interface {
	Generate(options Options) (*GeneratorResult, error)
}
