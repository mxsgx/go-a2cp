package main

import (
	"fmt"
	"log"

	a2cp "github.com/mxsgx/go-a2cp"
)

func main() {
	src := `
Listen 80
<VirtualHost *:80>
    ServerName example.com
    DocumentRoot "/var/www/html"
</VirtualHost>
`

	doc, err := a2cp.ParseString(src)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("top-level statements: %d\n", len(doc.Statements))
	for _, stmt := range doc.Statements {
		switch s := stmt.(type) {
		case a2cp.Directive:
			fmt.Printf("directive: %s %v\n", s.Name, s.Args)
		case *a2cp.Block:
			fmt.Printf("block: <%s %v> children=%d\n", s.Name, s.Args, len(s.Children))
		}
	}
}
