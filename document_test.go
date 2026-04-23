package a2cp

import (
	"strings"
	"testing"
)

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

func TestRoundTripFixtureFileWithComments(t *testing.T) {
	doc, err := ParseFile("testdata/roundtrip/comments.conf")
	if err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}

	rendered := doc.String()
	if !strings.Contains(rendered, "ServerRoot /etc/apache2 # server root") {
		t.Fatalf("rendered output missing inline directive comment:\n%s", rendered)
	}

	parsed, err := ParseString(rendered)
	if err != nil {
		t.Fatalf("ParseString(rendered) error = %v", err)
	}

	if got, want := countComments(parsed.Statements), countComments(doc.Statements); got != want {
		t.Fatalf("comment count mismatch: got %d, want %d", got, want)
	}
}

func countComments(stmts []Statement) int {
	total := 0
	for _, stmt := range stmts {
		switch s := stmt.(type) {
		case Comment, *Comment:
			total++
		case *Block:
			total += countComments(s.Children)
		}
	}
	return total
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

func TestDocumentAddCommentInlineOption(t *testing.T) {
	doc := NewDocument()
	doc.AddDirective("Listen", "8080")

	if err := doc.AddComment("app port", WithInlineComment()); err != nil {
		t.Fatalf("AddInlineComment() error = %v", err)
	}

	rendered := doc.String()
	if !strings.Contains(rendered, "Listen 8080 # app port") {
		t.Fatalf("rendered output missing inline comment:\n%s", rendered)
	}
}

func TestAddCommentWithInlineOptionWithoutStatementReturnsError(t *testing.T) {
	doc := NewDocument()
	if err := doc.AddComment("dangling", WithInlineComment()); err == nil {
		t.Fatalf("AddComment() expected error")
	}
}

func TestAddInlineCommentCompatibilityAlias(t *testing.T) {
	doc := NewDocument()
	doc.AddDirective("Listen", "8080")

	if err := doc.AddInlineComment("alias"); err != nil {
		t.Fatalf("AddInlineComment() error = %v", err)
	}

	if !strings.Contains(doc.String(), "Listen 8080 # alias") {
		t.Fatalf("rendered output missing alias inline comment")
	}
}

func TestAddCommentNormalizesDefaultText(t *testing.T) {
	doc := NewDocument()
	if err := doc.AddComment("app comment"); err != nil {
		t.Fatalf("AddComment() error = %v", err)
	}

	if got := doc.String(); got != "# app comment\n" {
		t.Fatalf("rendered output = %q, want %q", got, "# app comment\\n")
	}
}

func TestAddCommentWithRawCommentTextKeepsVerbatim(t *testing.T) {
	doc := NewDocument()
	if err := doc.AddComment("\tapp comment", WithRawCommentText()); err != nil {
		t.Fatalf("AddComment() error = %v", err)
	}

	if got := doc.String(); got != "#\tapp comment\n" {
		t.Fatalf("rendered output = %q, want %q", got, "#\\tapp comment\\n")
	}
}

func TestRenderEscapesCommentNewlines(t *testing.T) {
	doc := NewDocument()
	doc.AddDirective("Listen", "8080")

	if err := doc.AddComment("top\nline"); err != nil {
		t.Fatalf("AddComment() error = %v", err)
	}
	if err := doc.AddComment("inline\r\nline", WithInlineComment(), WithRawCommentText()); err != nil {
		t.Fatalf("AddComment() error = %v", err)
	}

	rendered := doc.String()
	if strings.Contains(rendered, "# top\nline") {
		t.Fatalf("rendered output contains unescaped newline in comment:\n%s", rendered)
	}
	if !strings.Contains(rendered, "# top\\nline") {
		t.Fatalf("rendered output missing escaped newline sequence:\n%s", rendered)
	}
	if !strings.Contains(rendered, "Listen 8080 #inline\\r\\nline") {
		t.Fatalf("rendered output missing escaped CRLF inline comment:\n%s", rendered)
	}
}

func TestAddCommentInlineOptionReplacesExistingInlineComment(t *testing.T) {
	doc := NewDocument()
	doc.AddDirective("Listen", "8080")

	if err := doc.AddComment("first", WithInlineComment()); err != nil {
		t.Fatalf("AddComment() error = %v", err)
	}
	if err := doc.AddComment("second", WithInlineComment()); err != nil {
		t.Fatalf("AddComment() error = %v", err)
	}

	rendered := doc.String()
	if !strings.Contains(rendered, "Listen 8080 # second") {
		t.Fatalf("rendered output missing replaced inline comment:\n%s", rendered)
	}
	if strings.Contains(rendered, "# first") {
		t.Fatalf("rendered output still contains old inline comment:\n%s", rendered)
	}
	if strings.Count(rendered, "#") != 1 {
		t.Fatalf("rendered output has duplicate inline comments:\n%s", rendered)
	}
}
