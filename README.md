# Jenga
![build](https://github.com/zschaffer/jenga/actions/workflows/go.yml/badge.svg)



A tool for building static single page blogs in Go.

## Details

Jenga is a no frills, fast-enough static site builder written in [Go](https://golang.org/). It is optimized for single-page infinite scrolling blogs.
Jenga takes a source directory of markdown files and an HTML template and spits outs a full HTML blog.

### Supported Platforms

In the releases tab you will find pre-built binaries for Linux, Windows and macOS (Intel and Apple Silicon). Otherwise, Jenga can compile and run anywhere Go can!

## Build and Install from Source

### Prerequisities

- [Git](https://git-scm.com/)
- [Go](https://golang.org/)

Clone the source from GitHub and install:

```bash
git clone https://github.com/zschaffer/jenga.git
cd jenga
go install
```

## Usage

Jenga has some basic setting up in order to get going; sort of like the real game!

### Set up your own `template.html` or copy it from the releases tab

Jenga uses Go's [`html/template`](https://pkg.go.dev/html/template) library for template construction. Read their [doc's](https://pkg.go.dev/html/template) for more information on how to manipulate your data. The basic thing required in your `template.html` is a `{{.}}` block to render the data converted from your `.md` files.

The included template.html file looks something like this:

```html
<body>
  <!-- Wrap everything in a div -->
  <div>
    <!-- Map over all your input .md files -->
    {{range .}}

    <!-- Wrap each input file in a div tag -->
    <div>{{.}}</div>

    <!-- End the map -->
    {{end}}
  </div>
</body>
```

### Set up your config `jenga.toml` file or copy it from the releases tab

Jenga uses [TOML]() as a configuration language (for now). [TOML]() is structured with `keys` and `values` like `mykey = myvalue`.

Jenga has three keys it looks for in order to run:

```toml
 InputFileDir = "/path/to/your/markdown/files"
 OutputDir = "/path/to/your/output/folder"
 TemplatePath = "/path/to/your/template.html"
```

> Note: TOML is case-sensitive so make sure you get those keys right!

### Run Jenga

Jenga only takes one flag, `-config`, that indicates to Jenga where your `jenga.toml` file is.

```bash
jenga    #If you don't pass any flags - default is "./jenga.toml"

# OR

jenga -config="$HOME/.config/jenga.toml"
```

Running Jenga will read through your source files, convert them to markdown, punch them into your template, and output them your your output directory - once you see the `build is finished!` you're good to go!

From here you can use Cloudflare Page's or GitHub Pages to host your new site. Just point them at your new build folder and bam!

## Dependencies

Jenga takes advantage of Go's super handy standard library for most things.

Other than that, Jenga currently relies on:

- [Burnt Sushi's toml](https://github.com/BurntSushi/toml.git)
- [gomarkdown's markdown](https://github.com/gomarkdown/markdown.git)

Shoutout to them because otherwise this would've been a lot trickier.
