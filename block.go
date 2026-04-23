package a2cp

import (
	"fmt"
	"strings"
)

// NewBlock creates a block statement.
func NewBlock(name string, args ...string) *Block {
	return &Block{Name: name, Args: args}
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

// Walk traverses this block and all descendants depth-first in pre-order.
// Return false from fn to stop traversal immediately.
func (b *Block) Walk(fn WalkFunc) {
	if fn == nil {
		return
	}
	stopped := false
	walkStatement(b, 0, fn, &stopped)
}

// FindDirectivesRecursive returns directives with the given name in descendant statements (case-insensitive).
func (b *Block) FindDirectivesRecursive(name string) []Directive {
	out := make([]Directive, 0)
	stopped := false
	walkStatements(b.Children, 0, func(stmt Statement, depth int) bool {
		directive, ok := asDirective(stmt)
		if ok && strings.EqualFold(directive.Name, name) {
			out = append(out, directive)
		}
		return true
	}, &stopped)
	return out
}

// FindBlocksRecursive returns blocks with the given name in descendant statements (case-insensitive).
func (b *Block) FindBlocksRecursive(name string) []*Block {
	out := make([]*Block, 0)
	stopped := false
	walkStatements(b.Children, 0, func(stmt Statement, depth int) bool {
		block, ok := stmt.(*Block)
		if ok && strings.EqualFold(block.Name, name) {
			out = append(out, block)
		}
		return true
	}, &stopped)
	return out
}
