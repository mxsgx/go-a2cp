# Apache2 Config Parser

`go-a2cp` is a small Go library for parsing and manipulating Apache2 `.conf` files.

## Features

- Parses directives (for example `Listen 80`)
- Parses nested block directives (for example `<VirtualHost *:80> ... </VirtualHost>`)
- Handles `#` comments (including inline comments)
- Supports quoted arguments (`"..."` and `'...'`)
- Supports line continuation using trailing `\\`
- Supports AST manipulation (append/insert/remove/find)
- Supports creating configs from scratch (builder-style API)
- Renders and writes modified config back to file

## Install

```bash
go get github.com/mxsgx/go-a2cp
```

## Compatibility

- Go 1.26+
- OS: Linux, macOS, Windows

## Quick Start

```go
package main

import (
	"fmt"

	"github.com/mxsgx/go-a2cp"
)

func main() {
	src := `
<VirtualHost *:80>
    ServerName example.com
    DocumentRoot "/var/www/html"
</VirtualHost>
`

	doc, err := a2cp.ParseString(src)
	if err != nil {
		panic(err)
	}

	for _, stmt := range doc.Statements {
		switch s := stmt.(type) {
		case a2cp.Directive:
			fmt.Printf("directive: %s %v\n", s.Name, s.Args)
		case *a2cp.Block:
			fmt.Printf("block: <%s %v> children=%d\n", s.Name, s.Args, len(s.Children))
		}
	}
}
```

## Public API

- `ParseString(src string) (*Document, error)`
- `ParseReader(r io.Reader) (*Document, error)`
- `ParseFile(path string) (*Document, error)`
- `NewDocument() *Document`
- `NewDirective(name string, args ...string) Directive`
- `NewBlock(name string, args ...string) *Block`
- `(*Document).String() string`
- `(*Document).WriteTo(w io.Writer) (int64, error)`
- `(*Document).Save(path string) error`
- `(*Document).SaveWithMode(path string, mode os.FileMode) error`
- `(*Document).AddDirective(name string, args ...string) *Document`
- `(*Document).AddBlock(name string, args ...string) *Block`
- `(*Block).AddDirective(name string, args ...string) *Block`
- `(*Block).AddBlock(name string, args ...string) *Block`
- AST nodes:
  - `Directive`
  - `Block`
  - `Document`
  - `Position`

## Testing

```bash
go test ./...
```

## Development

```bash
go test -v
```

Keep pull requests small and include tests for behavior changes.

## Examples Folder

Runnable examples are available in `examples/`:

- `examples/parse-string`: parse config from in-memory string
- `examples/parse-file`: parse a `.conf` file from disk
- `examples/manipulate-save`: modify AST and save generated config
- `examples/from-scratch`: build a full config from empty document and save it

Run each example from repository root:

```bash
go run ./examples/parse-string
go run ./examples/parse-file
go run ./examples/manipulate-save
go run ./examples/from-scratch
```

## Build From Scratch

```go
package main

import a2cp "github.com/mxsgx/go-a2cp"

func main() {
	doc := a2cp.NewDocument()
	doc.AddDirective("Listen", "8080")

	vh := doc.AddBlock("VirtualHost", "*:8080")
	vh.AddDirective("ServerName", "scratch.local")
	vh.AddDirective("DocumentRoot", "/var/www/scratch")

	_ = doc.Save("apache2.from-scratch.conf")
}
```

## Example (Inline)

```go
package main

import (
	"log"

	"github.com/mxsgx/go-a2cp"
)

func main() {
	doc, err := a2cp.ParseFile("apache2.conf")
	if err != nil {
		log.Fatal(err)
	}

	// Add a top-level directive.
	doc.Append(a2cp.Directive{Name: "ServerTokens", Args: []string{"Prod"}})

	// Edit the first VirtualHost block if present.
	vhosts := doc.FindBlocks("VirtualHost")
	if len(vhosts) > 0 {
		vhosts[0].Append(a2cp.Directive{Name: "ServerAdmin", Args: []string{"admin@example.com"}})
	}

	if err := doc.Save("apache2.generated.conf"); err != nil {
		log.Fatal(err)
	}
}
```
