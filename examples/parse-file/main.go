package main

import (
	"fmt"
	"log"

	a2cp "github.com/mxsgx/go-a2cp"
)

func main() {
	doc, err := a2cp.ParseFile("testdata/examples/sample-apache2.conf")
	if err != nil {
		log.Fatal(err)
	}

	listens := doc.FindDirectives("Listen")
	vhosts := doc.FindBlocks("VirtualHost")

	fmt.Printf("listen directives: %d\n", len(listens))
	fmt.Printf("virtualhost blocks: %d\n", len(vhosts))
}
