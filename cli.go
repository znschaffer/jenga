package main

import (
	"flag"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// config represents a toml file used to configure jenga
type config struct {
	InputDirPath  string
	OutputDirPath string
	TemplatePath  string
}

func run() error {
	configPath := flag.String("config", "jenga.toml", "/path/to/jenga.toml")
	flag.Parse()

	cfg, err := getConfig(*configPath)
	if err != nil {
		return fmt.Errorf("failed to get config (%q) %w", *configPath, err)
	}

	inputFilePaths, err := getInputFilePaths(cfg.InputDirPath)
	if err != nil {
		return fmt.Errorf("failed to get input file paths (%q) %w", cfg.InputDirPath, err)
	}

	template, err := getTemplate(cfg.TemplatePath)
	if err != nil {
		return fmt.Errorf("failed to get template (%q) %w", cfg.TemplatePath, err)
	}

	g := &builder{
		inputFilePaths: inputFilePaths,
		outputDirPath:  cfg.OutputDirPath,
		template:       template,
	}
	if err := g.build(); err != nil {
		return fmt.Errorf("failed to build files %w", err)
	}

	return nil
}

// getTemplate returns a pointer to a parsed Template from templatePath
func getTemplate(templatePath string) (*template.Template, error) {
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template (%q) %w", templatePath, err)
	}

	return t, nil
}

// getInputFilePaths returns file paths to every .md in input directory
func getInputFilePaths(inputDirPath string) ([]string, error) {
	var inputFilePaths []string

	inputDir, err := os.Open(inputDirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open %q %w", inputDirPath, err)
	}

	defer inputDir.Close()

	inputFiles, err := inputDir.Readdir(-1)
	if err != nil {
		return nil, fmt.Errorf("failed to read inputDir %w", err)
	}

	for _, file := range inputFiles {
		if !file.IsDir() && file.Name()[0] != '.' {
			inputFilePaths = append(inputFilePaths, filepath.Join(inputDirPath, file.Name()))
		}
	}

	return inputFilePaths, nil
}

// getConfig decodes a toml file at variable path, returning a config struct
func getConfig(path string) (*config, error) {
	cfg := config{}

	_, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to decode toml config %w", err)
	}

	if cfg.InputDirPath == "" {
		return nil, fmt.Errorf("failed to find InputDirPath in config")
	}

	if cfg.OutputDirPath == "" {
		return nil, fmt.Errorf("failed to find OutputDirPath in config")
	}

	if cfg.TemplatePath == "" {
		return nil, fmt.Errorf("failed to find TemplatePath in config")
	}

	return &cfg, nil
}
