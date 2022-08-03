package generator

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/gomarkdown/markdown"
	"github.com/zschaffer/jenga/config"
)

type Generator struct {
	Config *Config
}

type Config struct {
	SourceFiles []string
	Destination string
	Config      *config.Config
}

func New(config *Config) *Generator {
	return &Generator{Config: config}
}

func getTemplate(path string) (*template.Template, error) {
	temp, err := template.ParseFiles(path)
	if err != nil {
		return nil, fmt.Errorf("can't parse template files: %v", err)
	}

	return temp, nil
}

func readFile(path string) (template.HTML, error) {
	input, err := os.ReadFile(path)
	if err != nil {
		return template.HTML(""), fmt.Errorf("error while reading file %v: %v", path, err)
	}
	html := markdown.ToHTML(input, nil, nil)
	return template.HTML(html), nil
}

func (g *Generator) Generate() error {
	fmt.Println("generating...")

	templatePath := filepath.Join("static", "template.html")
	temp, err := getTemplate(templatePath)
	if err != nil {
		return fmt.Errorf("can't get template: %v", err)
	}

	sourceFiles := g.Config.SourceFiles
	destination := g.Config.Destination

	// if err := createDestination(destinationFolder); err != nil {
	// 	return err
	// }

	var posts []template.HTML
	for _, sourceFile := range sourceFiles {
		post, err := readFile(sourceFile)
		if err != nil {
			return err
		}
		posts = append(posts, post)
	}

	if err := writeFile(posts, destination, temp); err != nil {
		return err
	}
	return nil
}

func writeFile(posts []template.HTML, destination string, temp *template.Template) error {
	filePath := filepath.Join(destination, "index.html")
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %v", filePath, err)
	}
	defer f.Close()

	if err := temp.Execute(f, posts); err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	return nil
}
