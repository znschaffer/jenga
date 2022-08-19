package main

import (
	"html/template"
	"reflect"
	"testing"
)

func TestRun(t *testing.T) {
}

func TestParseInput(t *testing.T) {
}

func TestGetTemplate(t *testing.T) {
	got, err := getTemplate("./testdata/template.html")
	if err != nil {
		t.Errorf("failed to get template: %v", err)
	}

	want, _ := template.ParseFiles("./testdata/template.html")

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestGetInputFilePaths(t *testing.T) {
	t.Run("correctly return file paths to markdown files in input directory in reverse order", func(t *testing.T) {
		got, err := getInputFilePaths("./testdata/src")
		if err != nil {
			t.Errorf("failed to return file paths: %v", err)
		}

		want := []string{
			"testdata/src/test3.md",
			"testdata/src/test2.md",
			"testdata/src/test1.md",
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("return empty string array and err when passed path with no markdown files", func(t *testing.T) {
		_, err := getInputFilePaths("./testdata/empty")

		if err == nil {
			t.Errorf("error is empty, want: inputDirPath is empty")
		}
	})

	t.Run("return error when passed a broken path", func(t *testing.T) {
		_, err := getInputFilePaths("./testdata/path/does/not/exist")

		if err == nil {
			t.Errorf("error is empty, want: failed to open")
		}
	})
}

func TestGetConfig(t *testing.T) {
	t.Run("correctly parse passed config", func(t *testing.T) {
		got, err := getConfig("./testdata/jenga.toml")
		if err != nil {
			t.Errorf("failed to parse config: %v", err)
		}

		want := &config{
			InputDirPath:  "./src",
			OutputDirPath: "./build",
			TemplatePath:  "./template.html",
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("fails when passed bad path", func(t *testing.T) {
		_, err := getConfig("")

		if err == nil {
			t.Errorf("should return error when passed a bad path")
		}
	})
}
