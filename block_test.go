package a2cp

import (
	"strings"
	"testing"
)

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

func TestBlockAddCommentInlineOption(t *testing.T) {
	doc := NewDocument()
	vh := doc.AddBlock("VirtualHost", "*:8080")
	vh.AddDirective("ServerName", "scratch.local")

	if err := vh.AddComment("hostname", WithInlineComment()); err != nil {
		t.Fatalf("AddComment() error = %v", err)
	}

	rendered := doc.String()
	if !strings.Contains(rendered, "    ServerName scratch.local # hostname") {
		t.Fatalf("rendered output missing inline block comment:\n%s", rendered)
	}
}
