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

type generator struct {
	files    []string
	output   string
	template *template.Template
}

func readFile(path string) (template.HTML, error) {

	extensions := parser.CommonExtensions | parser.Attributes | parser.Mmark
	parser := parser.NewWithExtensions(extensions)

	input, err := os.ReadFile(path)
	if err != nil {
		return template.HTML(""), fmt.Errorf("error while reading file %v: %v", path, err)
	}
	html := markdown.ToHTML(input, parser, nil)
	return template.HTML(html), nil
}

func (g *generator) generate() error {
	fmt.Println("generating...")

	var posts []template.HTML
	for _, file := range g.files {
		post, err := readFile(file)
		if err != nil {
			return err
		}
		posts = append(posts, post)
	}

	// reverse post order
	for i, j := 0, len(posts)-1; i < j; i, j = i+1, j-1 {
		posts[i], posts[j] = posts[j], posts[i]
	}

	if err := writeFile(posts, g.output, g.template); err != nil {
		return err
	}
	return nil
}

func writeFile(posts []template.HTML, destination string, temp *template.Template) error {
	filePath := filepath.Join(destination, "index.html")
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %v", filePath, err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	defer writer.Flush()
	if err := temp.Execute(writer, posts); err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	return nil
}
