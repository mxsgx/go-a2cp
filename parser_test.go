package a2cp

import "testing"

func TestParseDirectivesAndComments(t *testing.T) {
	src := `
# global comment
ServerRoot "/etc/apache2"
Listen 80   # inline comment
`
	doc, err := ParseString(src)
	if err != nil {
		t.Fatalf("ParseString() error = %v", err)
	}

	if got := len(doc.Statements); got != 2 {
		t.Fatalf("statements = %d, want 2", got)
	}

	d0, ok := doc.Statements[0].(Directive)
	if !ok {
		t.Fatalf("statement[0] not Directive")
	}
	if d0.Name != "ServerRoot" {
		t.Fatalf("directive name = %q, want %q", d0.Name, "ServerRoot")
	}
	if len(d0.Args) != 1 || d0.Args[0] != "/etc/apache2" {
		t.Fatalf("directive args = %#v", d0.Args)
	}

	d1, ok := doc.Statements[1].(Directive)
	if !ok {
		t.Fatalf("statement[1] not Directive")
	}
	if d1.Name != "Listen" || len(d1.Args) != 1 || d1.Args[0] != "80" {
		t.Fatalf("directive mismatch: %#v", d1)
	}
}

func TestParseNestedBlocks(t *testing.T) {
	src := `
<VirtualHost *:80>
    ServerName example.com
    <Directory "/var/www/html">
        Require all granted
    </Directory>
</VirtualHost>
`
	doc, err := ParseString(src)
	if err != nil {
		t.Fatalf("ParseString() error = %v", err)
	}

	if len(doc.Statements) != 1 {
		t.Fatalf("top-level statements = %d, want 1", len(doc.Statements))
	}

	vh, ok := doc.Statements[0].(*Block)
	if !ok {
		t.Fatalf("statement[0] not *Block")
	}
	if vh.Name != "VirtualHost" {
		t.Fatalf("block name = %q, want VirtualHost", vh.Name)
	}
	if len(vh.Args) != 1 || vh.Args[0] != "*:80" {
		t.Fatalf("block args = %#v", vh.Args)
	}
	if len(vh.Children) != 2 {
		t.Fatalf("block children = %d, want 2", len(vh.Children))
	}

	dir, ok := vh.Children[1].(*Block)
	if !ok {
		t.Fatalf("VirtualHost child[1] not *Block")
	}
	if dir.Name != "Directory" {
		t.Fatalf("nested block name = %q, want Directory", dir.Name)
	}
}

func TestParseLineContinuation(t *testing.T) {
	src := "LogFormat \"%h %l \\\n%u %t\" common\n"
	doc, err := ParseString(src)
	if err != nil {
		t.Fatalf("ParseString() error = %v", err)
	}

	if len(doc.Statements) != 1 {
		t.Fatalf("statements = %d, want 1", len(doc.Statements))
	}
	d, ok := doc.Statements[0].(Directive)
	if !ok {
		t.Fatalf("statement[0] not Directive")
	}
	if d.Name != "LogFormat" {
		t.Fatalf("directive name = %q, want LogFormat", d.Name)
	}
	if len(d.Args) != 2 || d.Args[1] != "common" {
		t.Fatalf("directive args = %#v", d.Args)
	}
}

func TestParseErrors(t *testing.T) {
	cases := []string{
		"</Directory>",
		"<Directory /tmp>\n</Files>",
		"<VirtualHost *:80>",
		"ServerName \"example.com",
	}

	for _, tc := range cases {
		if _, err := ParseString(tc); err == nil {
			t.Fatalf("ParseString(%q) expected error", tc)
		}
	}
}

func TestParseFileFixture(t *testing.T) {
	doc, err := ParseFile("testdata/parser/basic.conf")
	if err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}

	if got := len(doc.Statements); got != 3 {
		t.Fatalf("statements = %d, want 3", got)
	}

	root, ok := doc.Statements[0].(Directive)
	if !ok || root.Name != "ServerRoot" {
		t.Fatalf("statement[0] mismatch: %#v", doc.Statements[0])
	}

	directory, ok := doc.Statements[2].(*Block)
	if !ok || directory.Name != "Directory" {
		t.Fatalf("statement[2] mismatch: %#v", doc.Statements[2])
	}
	if len(directory.Children) != 1 {
		t.Fatalf("Directory children = %d, want 1", len(directory.Children))
	}
}
