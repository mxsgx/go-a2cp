package a2cp

import (
	"strconv"
	"strings"
)

// String renders the document back to Apache2 configuration text.
func (d *Document) String() string {
	var b strings.Builder
	renderStatements(&b, d.Statements, 0)
	return b.String()
}

func renderStatements(b *strings.Builder, statements []Statement, depth int) {
	for _, stmt := range statements {
		renderStatement(b, stmt, depth)
	}
}

func renderStatement(b *strings.Builder, stmt Statement, depth int) {
	indent := strings.Repeat("    ", depth)

	switch s := stmt.(type) {
	case Directive:
		b.WriteString(indent)
		b.WriteString(s.Name)
		for _, arg := range s.Args {
			b.WriteString(" ")
			b.WriteString(renderArg(arg))
		}
		b.WriteString("\n")
	case *Directive:
		b.WriteString(indent)
		b.WriteString(s.Name)
		for _, arg := range s.Args {
			b.WriteString(" ")
			b.WriteString(renderArg(arg))
		}
		b.WriteString("\n")
	case *Block:
		b.WriteString(indent)
		b.WriteString("<")
		b.WriteString(s.Name)
		for _, arg := range s.Args {
			b.WriteString(" ")
			b.WriteString(renderArg(arg))
		}
		b.WriteString(">\n")

		renderStatements(b, s.Children, depth+1)

		b.WriteString(indent)
		b.WriteString("</")
		b.WriteString(s.Name)
		b.WriteString(">\n")
	}
}

func renderArg(arg string) string {
	if arg == "" {
		return `""`
	}
	if !needsQuote(arg) {
		return arg
	}
	return strconv.Quote(arg)
}

func needsQuote(arg string) bool {
	for _, r := range arg {
		if r == '#' || r == '"' || r == '\'' || r == '\\' {
			return true
		}
		if r == '<' || r == '>' || r == '=' || r == ';' {
			return true
		}
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			return true
		}
	}
	return false
}
