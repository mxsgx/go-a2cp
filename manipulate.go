package a2cp

import (
	"fmt"
	"strings"
)

// CommentOption configures how a comment is added to a document or block.
type CommentOption func(*commentOptions)

type commentOptions struct {
	inline  bool
	rawText bool
}

// WithInlineComment marks the comment as inline with the previous non-comment statement.
func WithInlineComment() CommentOption {
	return func(cfg *commentOptions) {
		cfg.inline = true
	}
}

// WithRawCommentText preserves comment text verbatim without inserting a leading space.
// By default, AddComment normalizes non-empty text so rendered comments use `# ` prefix.
func WithRawCommentText() CommentOption {
	return func(cfg *commentOptions) {
		cfg.rawText = true
	}
}

// NewDocument creates an empty configuration document.
func NewDocument() *Document {
	return &Document{}
}

// NewDirective creates a directive statement.
func NewDirective(name string, args ...string) Directive {
	return Directive{Name: name, Args: args}
}

// NewBlock creates a block statement.
func NewBlock(name string, args ...string) *Block {
	return &Block{Name: name, Args: args}
}

// Append adds a statement to the document root.
func (d *Document) Append(stmt Statement) {
	d.Statements = append(d.Statements, stmt)
}

// AddDirective appends a directive and returns the document for chaining.
func (d *Document) AddDirective(name string, args ...string) *Document {
	d.Append(NewDirective(name, args...))
	return d
}

// AddComment appends a comment to the document root.
// By default, non-empty text is normalized to include exactly one leading space after `#`
// when rendered (for example, text `"app port"` renders as `# app port`).
// Use WithRawCommentText to preserve text verbatim, and WithInlineComment to render on
// the same line as the previous statement.
func (d *Document) AddComment(text string, opts ...CommentOption) error {
	return addComment(&d.Statements, text, opts...)
}

// AddInlineComment appends an inline comment for the last non-comment root statement.
// Deprecated: use AddComment(text, WithInlineComment()) instead.
func (d *Document) AddInlineComment(text string) error {
	return d.AddComment(text, WithInlineComment())
}

// AddBlock appends a block and returns it for nested chaining.
func (d *Document) AddBlock(name string, args ...string) *Block {
	b := NewBlock(name, args...)
	d.Append(b)
	return b
}

// Insert inserts a statement at index in the document root.
func (d *Document) Insert(index int, stmt Statement) error {
	if index < 0 || index > len(d.Statements) {
		return fmt.Errorf("insert index out of range: %d", index)
	}
	d.Statements = append(d.Statements[:index], append([]Statement{stmt}, d.Statements[index:]...)...)
	return nil
}

// Remove removes and returns the statement at index in the document root.
func (d *Document) Remove(index int) (Statement, error) {
	if index < 0 || index >= len(d.Statements) {
		return nil, fmt.Errorf("remove index out of range: %d", index)
	}
	removed := d.Statements[index]
	d.Statements = append(d.Statements[:index], d.Statements[index+1:]...)
	return removed, nil
}

// FindDirectives returns root directives with the given name (case-insensitive).
func (d *Document) FindDirectives(name string) []Directive {
	out := make([]Directive, 0)
	for _, stmt := range d.Statements {
		directive, ok := asDirective(stmt)
		if !ok {
			continue
		}
		if strings.EqualFold(directive.Name, name) {
			out = append(out, directive)
		}
	}
	return out
}

// FindBlocks returns root blocks with the given name (case-insensitive).
func (d *Document) FindBlocks(name string) []*Block {
	out := make([]*Block, 0)
	for _, stmt := range d.Statements {
		block, ok := stmt.(*Block)
		if !ok {
			continue
		}
		if strings.EqualFold(block.Name, name) {
			out = append(out, block)
		}
	}
	return out
}

// Append adds a child statement to the block.
func (b *Block) Append(stmt Statement) {
	b.Children = append(b.Children, stmt)
}

// AddDirective appends a directive child and returns the block for chaining.
func (b *Block) AddDirective(name string, args ...string) *Block {
	b.Append(NewDirective(name, args...))
	return b
}

// AddComment appends a comment to the block.
// By default, non-empty text is normalized to include exactly one leading space after `#`
// when rendered (for example, text `"app port"` renders as `# app port`).
// Use WithRawCommentText to preserve text verbatim, and WithInlineComment to render on
// the same line as the previous statement.
func (b *Block) AddComment(text string, opts ...CommentOption) error {
	return addComment(&b.Children, text, opts...)
}

// AddInlineComment appends an inline comment for the last non-comment child statement.
// Deprecated: use AddComment(text, WithInlineComment()) instead.
func (b *Block) AddInlineComment(text string) error {
	return b.AddComment(text, WithInlineComment())
}

// AddBlock appends a nested block child and returns it.
func (b *Block) AddBlock(name string, args ...string) *Block {
	child := NewBlock(name, args...)
	b.Append(child)
	return child
}

// Insert inserts a child statement at index in the block.
func (b *Block) Insert(index int, stmt Statement) error {
	if index < 0 || index > len(b.Children) {
		return fmt.Errorf("insert index out of range: %d", index)
	}
	b.Children = append(b.Children[:index], append([]Statement{stmt}, b.Children[index:]...)...)
	return nil
}

// Remove removes and returns the child statement at index in the block.
func (b *Block) Remove(index int) (Statement, error) {
	if index < 0 || index >= len(b.Children) {
		return nil, fmt.Errorf("remove index out of range: %d", index)
	}
	removed := b.Children[index]
	b.Children = append(b.Children[:index], b.Children[index+1:]...)
	return removed, nil
}

// FindDirectives returns block child directives with the given name (case-insensitive).
func (b *Block) FindDirectives(name string) []Directive {
	out := make([]Directive, 0)
	for _, stmt := range b.Children {
		directive, ok := asDirective(stmt)
		if !ok {
			continue
		}
		if strings.EqualFold(directive.Name, name) {
			out = append(out, directive)
		}
	}
	return out
}

// FindBlocks returns block child blocks with the given name (case-insensitive).
func (b *Block) FindBlocks(name string) []*Block {
	out := make([]*Block, 0)
	for _, stmt := range b.Children {
		block, ok := stmt.(*Block)
		if !ok {
			continue
		}
		if strings.EqualFold(block.Name, name) {
			out = append(out, block)
		}
	}
	return out
}

func asDirective(stmt Statement) (Directive, bool) {
	switch d := stmt.(type) {
	case Directive:
		return d, true
	case *Directive:
		return *d, true
	default:
		return Directive{}, false
	}
}

func addComment(stmts *[]Statement, text string, opts ...CommentOption) error {
	cfg := commentOptions{}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}

	text = normalizeCommentText(text, cfg.rawText)

	if !cfg.inline {
		*stmts = append(*stmts, Comment{Text: text})
		return nil
	}

	idx := lastNonCommentIndex(*stmts)
	if idx < 0 {
		return fmt.Errorf("inline comment requires at least one non-comment statement")
	}

	line := statementLine((*stmts)[idx])
	if line <= 0 {
		line = nextSyntheticLine(*stmts)
		if err := setStatementLine(stmts, idx, line); err != nil {
			return err
		}
	}

	comment := Comment{Text: text, Pos: Position{Line: line, Column: 1}}
	insertAt := idx + 1
	for insertAt < len(*stmts) {
		switch s := (*stmts)[insertAt].(type) {
		case Comment:
			if s.Pos.Line != line {
				goto insert
			}
		case *Comment:
			if s.Pos.Line != line {
				goto insert
			}
		default:
			goto insert
		}
		insertAt++
	}

insert:
	*stmts = append((*stmts)[:insertAt], append([]Statement{comment}, (*stmts)[insertAt:]...)...)
	return nil
}

func normalizeCommentText(text string, raw bool) string {
	if raw || text == "" {
		return text
	}
	return " " + strings.TrimLeft(text, " \t")
}

func lastNonCommentIndex(stmts []Statement) int {
	for i := len(stmts) - 1; i >= 0; i-- {
		switch stmts[i].(type) {
		case Comment, *Comment:
			continue
		default:
			return i
		}
	}
	return -1
}

func statementLine(stmt Statement) int {
	switch s := stmt.(type) {
	case Directive:
		return s.Pos.Line
	case *Directive:
		return s.Pos.Line
	case *Block:
		return s.Pos.Line
	default:
		return 0
	}
}

func setStatementLine(stmts *[]Statement, index int, line int) error {
	switch s := (*stmts)[index].(type) {
	case Directive:
		s.Pos.Line = line
		(*stmts)[index] = s
		return nil
	case *Directive:
		s.Pos.Line = line
		return nil
	case *Block:
		s.Pos.Line = line
		return nil
	default:
		return fmt.Errorf("last statement does not support inline comments")
	}
}

func nextSyntheticLine(stmts []Statement) int {
	maxLine := 0
	for _, stmt := range stmts {
		line := statementLine(stmt)
		if line > maxLine {
			maxLine = line
		}
	}
	if maxLine == 0 {
		return 1
	}
	return maxLine + 1
}
