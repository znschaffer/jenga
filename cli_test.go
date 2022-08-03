package main

import (
	"reflect"
	"testing"
)

func TestRun(t *testing.T) {

}

func TestParseInput(t *testing.T) {

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
