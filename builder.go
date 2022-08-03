package main

import (
	"bufio"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
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
		return template.HTML(""), fmt.Errorf("error while reading file %v: %v", filePath, err)
	}

	html := markdown.ToHTML(data, parser, nil)
	return template.HTML(html), nil
}

// build reads all the files in the inputFilePaths slice, then passes them to writeOutputFile to build an index.html
func (g *builder) build() error {
	var inputData []template.HTML

	for _, inputFilePath := range g.inputFilePaths {
		inputFileData, err := readFile(inputFilePath)
		if err != nil {
			return err
		}
		inputData = append(inputData, inputFileData)
	}

	for i, j := 0, len(inputData)-1; i < j; i, j = i+1, j-1 {
		inputData[i], inputData[j] = inputData[j], inputData[i]
	}

	fmt.Println("~ building ~")

	if err := writeOutputFile(inputData, g.outputDirPath, g.template); err != nil {
		return fmt.Errorf("failed to write to output file: %v", err)
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

	fmt.Println(inputData)
	if err := t.Execute(writer, inputData); err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	return nil
}
