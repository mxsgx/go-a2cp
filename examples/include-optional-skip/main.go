package main

import (
	"fmt"
	"log"

	a2cp "github.com/mxsgx/go-a2cp"
)

func main() {
	doc, err := a2cp.ParseFile(
		"testdata/examples/include-optional-skip/apache2.main.conf",
		a2cp.WithIncludeResolution("testdata/examples/include-optional-skip"),
	)
	if err != nil {
		log.Fatal(err)
	}

	listens := doc.FindDirectives("Listen")
	vhosts := doc.FindBlocks("VirtualHost")

	fmt.Printf("top-level statements: %d\n", len(doc.Statements))
	fmt.Printf("listen directives: %d\n", len(listens))
	fmt.Printf("virtualhost blocks: %d\n", len(vhosts))
}
