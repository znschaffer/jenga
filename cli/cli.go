package cli

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/zschaffer/jenga/config"
	"github.com/zschaffer/jenga/generator"
)

func Run() {
	// TODO: Read Config
	cfg, err := readConfig()

	if err != nil {
		log.Fatalf("can't read config: %v", err)
	}

	sourceFiles, err := getContentFolders(cfg.SourceFiles)

	if err != nil {
		log.Fatalf("can't get content folders: %v", err)
	}

	// TODO: Create new generator and pass in config
	gen := generator.New(&generator.Config{
		SourceFiles: sourceFiles,
		Destination: cfg.Destination,
		Config:      cfg,
	})

	if err := gen.Generate(); err != nil {
		log.Fatalf("can't generate files: %v", err)
	}

	fmt.Println("done generating!")
}

func getContentFolders(path string) ([]string, error) {
	var result []string
	dir, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
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

func readConfig() (*config.Config, error) {
	configDir := "/Users/zane/.config"

	configPath := filepath.Join(configDir, "jenga.toml")

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %s", err)
	}

	cfg := config.Config{}

	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %v", err)
	}

	return &cfg, nil
}
