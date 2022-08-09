package main

import (
	"bufio"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"golang.org/x/sync/errgroup"
)

type builder struct {
	inputFilePaths []string
	outputDirPath  string
	template       *template.Template
}

// readFile reads in a markdown file, converts it to HTML and returns the HTML string
func readFile(filePath string) (template.HTML, error) {
	extensions := parser.CommonExtensions | parser.Attributes | parser.Mmark
	parser := parser.NewWithExtensions(extensions)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return template.HTML(""), fmt.Errorf("error while reading file %s: %w", filePath, err)
	}

	html := markdown.ToHTML(data, parser, nil)
	return template.HTML(html), nil
}

// build reads all the files in the inputFilePaths slice, then passes them to writeOutputFile to build an index.html
func (b *builder) build() error {
	var inputData []template.HTML
	temporaryMap := make(map[string]template.HTML)
	g := new(errgroup.Group)

	for _, inputFilePath := range b.inputFilePaths {
		path := inputFilePath
		g.Go(func() error {
			data, err := readFile(path)
			if err == nil {
				temporaryMap[path] = data
			}
			return err
		})
	}

	if err := g.Wait(); err != nil {
		return err
	} else {
		for _, inputFilePath := range b.inputFilePaths {
			inputData = append([]template.HTML{temporaryMap[inputFilePath]}, inputData...)
		}

	}

	fmt.Printf("building files to %s\n", b.outputDirPath)
	if err := writeOutputFile(inputData, b.outputDirPath, b.template); err != nil {
		return fmt.Errorf("failed to write to output file: %w", err)
	}
	return nil
}

// writeOutputFile creates an index.html file at outputDirPath using a template filled with inputData
func writeOutputFile(inputData []template.HTML, outputDirPath string, t *template.Template) error {
	filePath := filepath.Join(outputDirPath, "index.html")
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %+s: %v", filePath, err)
	}

	defer file.Close()

	writer := bufio.NewWriter(file)

	defer writer.Flush()

	if err := t.Execute(writer, inputData); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	return nil
}
