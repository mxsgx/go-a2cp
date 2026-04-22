package main

import (
	"fmt"
	"log"

	a2cp "github.com/mxsgx/go-a2cp"
)

func main() {
	doc, err := a2cp.ParseFile("testdata/examples/backslash-comments.conf")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("rendered config with backslash continuation + comments:")
	fmt.Println(doc.String())
}
