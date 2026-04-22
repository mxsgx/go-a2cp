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
	for i := 0; i < len(statements); i++ {
		stmt := statements[i]
		if _, ok := asComment(stmt); ok {
			renderStatement(b, stmt, depth, "", false)
			continue
		}

		inlineText := ""
		hasInline := false
		if i+1 < len(statements) {
			nextComment, ok := asComment(statements[i+1])
			if ok {
				stmtPos, hasPos := statementPos(stmt)
				if hasPos && stmtPos.Line > 0 && nextComment.Pos.Line == stmtPos.Line {
					inlineText = nextComment.Text
					hasInline = true
					i++
				}
			}
		}

		renderStatement(b, stmt, depth, inlineText, hasInline)
	}
}

func renderStatement(b *strings.Builder, stmt Statement, depth int, inlineText string, hasInline bool) {
	indent := strings.Repeat("    ", depth)

	switch s := stmt.(type) {
	case Directive:
		b.WriteString(indent)
		b.WriteString(s.Name)
		for _, arg := range s.Args {
			b.WriteString(" ")
			b.WriteString(renderArg(arg))
		}
		if hasInline {
			b.WriteString(" #")
			b.WriteString(inlineText)
		}
		b.WriteString("\n")
	case *Directive:
		b.WriteString(indent)
		b.WriteString(s.Name)
		for _, arg := range s.Args {
			b.WriteString(" ")
			b.WriteString(renderArg(arg))
		}
		if hasInline {
			b.WriteString(" #")
			b.WriteString(inlineText)
		}
		b.WriteString("\n")
	case Comment:
		b.WriteString(indent)
		b.WriteString("#")
		b.WriteString(s.Text)
		b.WriteString("\n")
	case *Comment:
		b.WriteString(indent)
		b.WriteString("#")
		b.WriteString(s.Text)
		b.WriteString("\n")
	case *Block:
		b.WriteString(indent)
		b.WriteString("<")
		b.WriteString(s.Name)
		for _, arg := range s.Args {
			b.WriteString(" ")
			b.WriteString(renderArg(arg))
		}
		if hasInline {
			b.WriteString("> #")
			b.WriteString(inlineText)
		} else {
			b.WriteString(">")
		}
		b.WriteString("\n")

		renderStatements(b, s.Children, depth+1)

		b.WriteString(indent)
		b.WriteString("</")
		b.WriteString(s.Name)
		b.WriteString(">\n")
	}
}

func asComment(stmt Statement) (Comment, bool) {
	switch c := stmt.(type) {
	case Comment:
		return c, true
	case *Comment:
		return *c, true
	default:
		return Comment{}, false
	}
}

func statementPos(stmt Statement) (Position, bool) {
	switch s := stmt.(type) {
	case Directive:
		return s.Pos, true
	case *Directive:
		return s.Pos, true
	case *Block:
		return s.Pos, true
	case Comment:
		return s.Pos, true
	case *Comment:
		return s.Pos, true
	default:
		return Position{}, false
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
