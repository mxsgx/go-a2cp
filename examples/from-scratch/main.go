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
	vh.AddDirective("DocumentRoot", "/var/www/scratch")

	dir := vh.AddBlock("Directory", "/var/www/scratch")
	dir.AddDirective("Require", "all", "granted")

	if err := doc.Save("examples/from-scratch/apache2.from-scratch.conf"); err != nil {
		log.Fatal(err)
	}

	log.Println("wrote examples/from-scratch/apache2.from-scratch.conf")
}
