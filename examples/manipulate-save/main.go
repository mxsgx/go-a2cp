package main

import (
	"log"

	a2cp "github.com/mxsgx/go-a2cp"
)

func main() {
	doc, err := a2cp.ParseFile("testdata/examples/sample-apache2.conf")
	if err != nil {
		log.Fatal(err)
	}

	doc.Append(a2cp.Directive{Name: "ServerTokens", Args: []string{"Prod"}})

	vhosts := doc.FindBlocks("VirtualHost")
	if len(vhosts) > 0 {
		vhosts[0].Append(a2cp.Directive{Name: "ServerAdmin", Args: []string{"admin@example.com"}})
	}

	if err := doc.Save("examples/manipulate-save/apache2.generated.conf"); err != nil {
		log.Fatal(err)
	}

	log.Println("wrote examples/manipulate-save/apache2.generated.conf")
}
