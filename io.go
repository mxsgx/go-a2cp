package a2cp

import (
	"io"
	"os"
)

// ParseFile parses an Apache2 .conf file from disk.
func ParseFile(path string, opts ...ParseOption) (*Document, error) {
	cfg := applyParseOptions(opts)
	state := &includeState{inProgress: make(map[string]struct{})}
	return parseIncludedFile(path, cfg, state)
}

// WriteTo writes the document as Apache2 config to w.
func (d *Document) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, d.String())
	return int64(n), err
}

// Save writes the document to path with file mode 0644.
func (d *Document) Save(path string) error {
	return d.SaveWithMode(path, 0o644)
}

// SaveWithMode writes the document to path using the provided file mode.
func (d *Document) SaveWithMode(path string, mode os.FileMode) error {
	return os.WriteFile(path, []byte(d.String()), mode)
}
