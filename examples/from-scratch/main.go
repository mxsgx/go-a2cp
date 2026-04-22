package main

import (
	"log"

	a2cp "github.com/mxsgx/go-a2cp"
)

func main() {
	doc := a2cp.NewDocument()
	doc.AddDirective("ServerTokens", "Prod")
	doc.AddDirective("Listen", "8080")

	vh := doc.AddBlock("VirtualHost", "*:8080")
	vh.AddDirective("ServerName", "scratch.local")

	if err := vh.AddComment(" You can change DocumentRoot below"); err != nil {
		log.Fatal(err)
	}
	vh.AddDirective("DocumentRoot", "/var/www/scratch")
	if err := vh.AddComment(" change the /var/www/scratch path", a2cp.WithInlineComment()); err != nil {
		log.Fatal(err)
	}

	dir := vh.AddBlock("Directory", "/var/www/scratch")
	dir.AddDirective("Require", "all", "granted")

	if err := doc.Save("examples/from-scratch/apache2.from-scratch.conf"); err != nil {
		log.Fatal(err)
	}

	log.Println("wrote examples/from-scratch/apache2.from-scratch.conf")
}
