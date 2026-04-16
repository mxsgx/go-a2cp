package a2cp

import "testing"

func TestDocumentManipulation(t *testing.T) {
	doc := &Document{}
	doc.Append(Directive{Name: "Listen", Args: []string{"80"}})
	doc.Append(Directive{Name: "ServerTokens", Args: []string{"Prod"}})

	if len(doc.Statements) != 2 {
		t.Fatalf("statements = %d, want 2", len(doc.Statements))
	}

	if err := doc.Insert(1, Directive{Name: "ServerName", Args: []string{"example.com"}}); err != nil {
		t.Fatalf("Insert() error = %v", err)
	}

	if len(doc.FindDirectives("servername")) != 1 {
		t.Fatalf("FindDirectives(servername) mismatch")
	}

	removed, err := doc.Remove(0)
	if err != nil {
		t.Fatalf("Remove() error = %v", err)
	}
	if d, ok := removed.(Directive); !ok || d.Name != "Listen" {
		t.Fatalf("removed statement mismatch: %#v", removed)
	}
}

func TestBlockManipulation(t *testing.T) {
	b := &Block{Name: "VirtualHost", Args: []string{"*:80"}}
	b.Append(Directive{Name: "ServerName", Args: []string{"example.com"}})
	b.Append(&Block{Name: "Directory", Args: []string{"/var/www/html"}})

	if got := len(b.FindDirectives("ServerName")); got != 1 {
		t.Fatalf("FindDirectives() = %d, want 1", got)
	}
	if got := len(b.FindBlocks("directory")); got != 1 {
		t.Fatalf("FindBlocks() = %d, want 1", got)
	}

	if _, err := b.Remove(1); err != nil {
		t.Fatalf("Remove() error = %v", err)
	}
	if got := len(b.FindBlocks("directory")); got != 0 {
		t.Fatalf("FindBlocks() after remove = %d, want 0", got)
	}
}

func TestRoundTripRenderAndParse(t *testing.T) {
	doc := &Document{}
	vh := &Block{Name: "VirtualHost", Args: []string{"*:443"}}
	vh.Append(Directive{Name: "ServerName", Args: []string{"example.com"}})
	vh.Append(Directive{Name: "DocumentRoot", Args: []string{"/var/www/site root"}})
	doc.Append(vh)

	rendered := doc.String()
	parsed, err := ParseString(rendered)
	if err != nil {
		t.Fatalf("ParseString(rendered) error = %v", err)
	}
	if len(parsed.Statements) != 1 {
		t.Fatalf("parsed statements = %d, want 1", len(parsed.Statements))
	}
}

func TestRoundTripFixtureFile(t *testing.T) {
	doc, err := ParseFile("testdata/roundtrip/virtualhost.conf")
	if err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}

	rendered := doc.String()
	parsed, err := ParseString(rendered)
	if err != nil {
		t.Fatalf("ParseString(rendered) error = %v", err)
	}

	if len(parsed.Statements) != len(doc.Statements) {
		t.Fatalf("statement count mismatch: got %d, want %d", len(parsed.Statements), len(doc.Statements))
	}
}

func TestBuildFromScratch(t *testing.T) {
	doc := NewDocument()
	doc.AddDirective("Listen", "8080")

	vh := doc.AddBlock("VirtualHost", "*:8080")
	vh.AddDirective("ServerName", "scratch.local")
	vh.AddDirective("DocumentRoot", "/var/www/scratch")

	dir := vh.AddBlock("Directory", "/var/www/scratch")
	dir.AddDirective("Require", "all", "granted")

	if got := len(doc.FindDirectives("Listen")); got != 1 {
		t.Fatalf("FindDirectives(Listen) = %d, want 1", got)
	}
	if got := len(doc.FindBlocks("VirtualHost")); got != 1 {
		t.Fatalf("FindBlocks(VirtualHost) = %d, want 1", got)
	}

	rendered := doc.String()
	parsed, err := ParseString(rendered)
	if err != nil {
		t.Fatalf("ParseString(rendered) error = %v", err)
	}
	if got := len(parsed.Statements); got != 2 {
		t.Fatalf("parsed top-level statements = %d, want 2", got)
	}
}
