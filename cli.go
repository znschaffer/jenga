package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type config struct {
	Input        string
	Output       string
	TemplatePath string
}

func getTemplate(path string) (*template.Template, error) {
	temp, err := template.ParseFiles(path)
	if err != nil {
		return nil, fmt.Errorf("can't parse template files: %v", err)
	}

	return temp, nil
}

func run() {

	configPath := flag.String("config", "jenga.toml", "/path/to/jenga.toml")

	flag.Parse()

	cfg, err := parseConfig(*configPath)
	if err != nil {
		log.Fatalf("can't read config: %v", err)
	}

	files, err := parseInput(cfg.Input)

	if err != nil {
		log.Fatalf("can't get content folders: %v", err)
	}

	template, err := getTemplate(cfg.TemplatePath)
	if err != nil {
		log.Fatalf("can't get template: %v", err)
	}

	gen := &generator{files: files, output: cfg.Output, template: template}

	if err := gen.generate(); err != nil {
		log.Fatalf("can't generate files: %v", err)
	}

	fmt.Println("done generating!")
}

func parseInput(path string) ([]string, error) {
	var result []string
	dir, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("can't open path %s: %v", path, err)
	}

	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	for _, file := range files {
		if !file.IsDir() && file.Name()[0] != '.' {
			result = append(result, filepath.Join(path, file.Name()))
		}
	}

	return result, nil
}

func parseConfig(path string) (*config, error) {
	cfg := config{}
	_, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %v", err)
	}

	return &cfg, nil
}
