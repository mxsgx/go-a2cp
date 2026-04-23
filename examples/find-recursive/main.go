package main

import (
	"fmt"
	"log"

	a2cp "github.com/mxsgx/go-a2cp"
)

func main() {
	src := `
<VirtualHost *:443>
    ServerName root.example
    <Directory "/srv/site">
        ServerName app.example
        <Directory "/srv/site/private">
            ServerName private.example
        </Directory>
    </Directory>
</VirtualHost>
`

	doc, err := a2cp.ParseString(src)
	if err != nil {
		log.Fatal(err)
	}

	directives := doc.FindDirectivesRecursive("servername")
	fmt.Printf("recursive ServerName directives: %d\n", len(directives))
	for _, d := range directives {
		fmt.Printf("- %s\n", d.Args[0])
	}

	blocks := doc.FindBlocksRecursive("directory")
	fmt.Printf("recursive Directory blocks: %d\n", len(blocks))
	for _, b := range blocks {
		fmt.Printf("- <%s %v>\n", b.Name, b.Args)
	}
}
