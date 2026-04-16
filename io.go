package a2cp

import (
	"io"
	"os"
)

// ParseFile parses an Apache2 .conf file from disk.
func ParseFile(path string) (*Document, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ParseReader(f)
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
