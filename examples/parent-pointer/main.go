package main

import (
	"fmt"
	"log"

	a2cp "github.com/mxsgx/go-a2cp"
)

func main() {
	src := `<VirtualHost *:443>
    ServerName example.com
    <Directory /var/www/html>
        <IfModule mod_authz_core.c>
            Require all granted
        </IfModule>
    </Directory>
</VirtualHost>`

	doc, err := a2cp.ParseString(src)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("parsed parent links")
	showParents(doc)

	fmt.Println("")
	fmt.Println("mutation parent links")
	if err := mutateAndShow(doc); err != nil {
		log.Fatal(err)
	}
}

func showParents(doc *a2cp.Document) {
	vhosts := doc.FindBlocks("VirtualHost")
	if len(vhosts) == 0 {
		fmt.Println("no VirtualHost blocks")
		return
	}

	vh := vhosts[0]
	dirs := vh.FindBlocks("Directory")
	if len(dirs) == 0 {
		fmt.Println("no Directory blocks")
		return
	}

	dir := dirs[0]
	mods := dir.FindBlocks("IfModule")
	if len(mods) == 0 {
		fmt.Println("no IfModule blocks")
		return
	}

	mod := mods[0]

	fmt.Printf("VirtualHost depth=%d isRoot=%t\n", vh.Depth(), vh.IsRoot())
	fmt.Printf("Directory depth=%d isRoot=%t\n", dir.Depth(), dir.IsRoot())
	fmt.Printf("IfModule depth=%d isRoot=%t\n", mod.Depth(), mod.IsRoot())
	fmt.Printf("IfModule parent is Directory: %t\n", mod.Parent == dir)
}

func mutateAndShow(doc *a2cp.Document) error {
	vh := doc.AddBlock("VirtualHost", "*:8080")
	dir := vh.AddBlock("Directory", "/srv/www")

	fmt.Printf("new VirtualHost depth=%d parentIsDoc=%t\n", vh.Depth(), vh.Parent == doc)
	fmt.Printf("new Directory depth=%d parentIsVHost=%t\n", dir.Depth(), dir.Parent == vh)

	removed, err := vh.Remove(0)
	if err != nil {
		return err
	}

	removedBlock, ok := removed.(*a2cp.Block)
	if ok {
		fmt.Printf("removed block parent nil=%t\n", removedBlock.Parent == nil)
	}

	return nil
}
