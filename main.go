/*
Jenga is a static site builder focused on single page blogs.  It uses a TOML
file for configuration ( jenga.toml ).  It expects your source files to be
written in standard Markdown with the extension .md

Usage:

	jenga -config '/path/to/config'

The config flag is optional and if not recieved, jenga will look for the config
in the current directory
*/
package main

import (
	"log"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
