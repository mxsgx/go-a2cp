package a2cp

import (
	"fmt"
	"strings"
)

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
