package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/rjeczalik/notify"
)

// config represents a toml file used to configure jenga
type config struct {
	InputDirPath  string
	OutputDirPath string
	TemplatePath  string
}

const AppVersion = "v0.1.1"

func run() error {
	version := flag.Bool("v", false, "prints current jenga version")
	dev := flag.Bool("dev", false, "watchs source folder and rebuilds on changes")
	configPath := flag.String("config", "./jenga.toml", "path to jenga.toml config")
	flag.Parse()

	args := flag.Args()
	if len(args) != 0 {
		fmt.Printf("unknown arguments: %v\n", args)
		fmt.Println("use jenga -h for accepted arguments")
		os.Exit(1)
	}

	if *version {
		fmt.Println(AppVersion)
		os.Exit(0)
	}

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

	b := &builder{
		inputFilePaths: inputFilePaths,
		outputDirPath:  cfg.OutputDirPath,
		template:       template,
	}

	fmt.Printf("\033[0;34mconfig\033[0m   = %q\n", *configPath)
	fmt.Printf("\033[0;34mtemplate\033[0m = %q\n", cfg.TemplatePath)
	fmt.Printf("\033[0;34minput\033[0m    = %q\n", cfg.InputDirPath)
	fmt.Printf("\033[0;34moutput\033[0m   = %q\n", cfg.OutputDirPath)

	if *dev {
		if err := watch(b, cfg); err != nil {
			return fmt.Errorf("failed to watch files: %w", err)
		}
	}

	if err := b.build(); err != nil {
		return fmt.Errorf("failed to build files: %w", err)
	}

	return nil
}

func watch(b *builder, cfg *config) error {
	c := make(chan notify.EventInfo, 1)
	if err := notify.Watch(cfg.InputDirPath, c, notify.Write, notify.Remove); err != nil {
		log.Fatal(err)
	}
	defer notify.Stop(c)

	srv := http.FileServer(http.Dir(cfg.OutputDirPath))
	go http.ListenAndServe(":3000", srv)

	fmt.Println()
	fmt.Println("listening @ \033[0;34mhttp://localhost:3000\033[0m")
	fmt.Println()
	fmt.Println("watching \033[0;34minput\033[0m for changes...")

	if err := b.build(); err != nil {
		return fmt.Errorf("failed to build files: %w", err)
	}

	for range c {
		if err := b.build(); err != nil {
			return fmt.Errorf("failed to build files: %w", err)
		}
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

// getInputFilePaths returns file paths to every .md in input directory in reverse alphabetical order
func getInputFilePaths(inputDirPath string) ([]string, error) {
	var inputFilePaths []string

	inputDir, err := os.Open(inputDirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open %q %w", inputDirPath, err)
	}

	defer inputDir.Close()

	inputFiles, err := inputDir.ReadDir(-1)
	if err != nil {
		return nil, fmt.Errorf("failed to read inputDir %w", err)
	}

	if len(inputFiles) == 0 {
		return nil, errors.New("inputDirPath is empty")
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
		return nil, errors.New("failed to find InputDirPath in config")
	}

	if cfg.OutputDirPath == "" {
		return nil, errors.New("failed to find OutputDirPath in config")
	}

	if cfg.TemplatePath == "" {
		return nil, errors.New("failed to find TemplatePath in config")
	}

	return &cfg, nil
}
