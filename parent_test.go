package a2cp

import "testing"

func TestParentsSetAfterParse(t *testing.T) {
	doc, err := ParseString(`<VirtualHost *:80>
    <Directory /var/www/html>
        Require all granted
    </Directory>
</VirtualHost>`)
	if err != nil {
		t.Fatalf("ParseString() error = %v", err)
	}

	root := doc.FindBlocks("VirtualHost")
	if len(root) != 1 {
		t.Fatalf("root block count = %d, want 1", len(root))
	}

	vh := root[0]
	if vh.Parent != doc {
		t.Fatalf("root block parent not document")
	}
	if !vh.IsRoot() {
		t.Fatalf("root block IsRoot() = false, want true")
	}

	dirs := vh.FindBlocks("Directory")
	if len(dirs) != 1 {
		t.Fatalf("nested block count = %d, want 1", len(dirs))
	}

	dir := dirs[0]
	if dir.Parent != vh {
		t.Fatalf("nested block parent mismatch")
	}
	if dir.IsRoot() {
		t.Fatalf("nested block IsRoot() = true, want false")
	}
}

func TestParentSetAfterAddBlock(t *testing.T) {
	doc := NewDocument()
	vh := doc.AddBlock("VirtualHost", "*:8080")
	dir := vh.AddBlock("Directory", "/var/www/html")

	if vh.Parent != doc {
		t.Fatalf("document AddBlock parent mismatch")
	}
	if dir.Parent != vh {
		t.Fatalf("block AddBlock parent mismatch")
	}
}

func TestBlockDepthNested(t *testing.T) {
	doc := NewDocument()
	l0 := doc.AddBlock("VirtualHost", "*:8080")
	l1 := l0.AddBlock("Directory", "/var/www/html")
	l2 := l1.AddBlock("IfModule", "mod_ssl.c")

	if got := l0.Depth(); got != 0 {
		t.Fatalf("level0 Depth() = %d, want 0", got)
	}
	if got := l1.Depth(); got != 1 {
		t.Fatalf("level1 Depth() = %d, want 1", got)
	}
	if got := l2.Depth(); got != 2 {
		t.Fatalf("level2 Depth() = %d, want 2", got)
	}
}

func TestParentNilAfterRemove(t *testing.T) {
	doc := NewDocument()
	vh := doc.AddBlock("VirtualHost", "*:8080")
	dir := vh.AddBlock("Directory", "/var/www/html")

	removed, err := vh.Remove(0)
	if err != nil {
		t.Fatalf("Block.Remove() error = %v", err)
	}
	removedBlock, ok := removed.(*Block)
	if !ok {
		t.Fatalf("removed statement type = %T, want *Block", removed)
	}
	if removedBlock != dir {
		t.Fatalf("removed block mismatch")
	}
	if dir.Parent != nil {
		t.Fatalf("removed child parent not nil")
	}

	removedRoot, err := doc.Remove(0)
	if err != nil {
		t.Fatalf("Document.Remove() error = %v", err)
	}
	removedRootBlock, ok := removedRoot.(*Block)
	if !ok {
		t.Fatalf("removed root statement type = %T, want *Block", removedRoot)
	}
	if removedRootBlock != vh {
		t.Fatalf("removed root block mismatch")
	}
	if vh.Parent != nil {
		t.Fatalf("removed root parent not nil")
	}
}
