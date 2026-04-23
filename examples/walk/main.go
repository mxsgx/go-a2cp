package main

import (
	"fmt"
	"log"
	"strings"

	a2cp "github.com/mxsgx/go-a2cp"
)

func main() {
	src := `
Listen 80
<VirtualHost *:80>
    ServerName example.com
    <Directory "/var/www/html">
        Require all granted
    </Directory>
</VirtualHost>
`

	doc, err := a2cp.ParseString(src)
	if err != nil {
		log.Fatal(err)
	}

	doc.Walk(func(stmt a2cp.Statement, depth int) bool {
		indent := strings.Repeat("  ", depth)
		switch s := stmt.(type) {
		case a2cp.Directive:
			fmt.Printf("%sdirective: %s %v\n", indent, s.Name, s.Args)
		case *a2cp.Block:
			fmt.Printf("%sblock: <%s %v>\n", indent, s.Name, s.Args)
		case a2cp.Comment:
			fmt.Printf("%scomment: #%s\n", indent, s.Text)
		}
		return true
	})
}
