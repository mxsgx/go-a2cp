package a2cp

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"
)

// ParseError contains source location and parsing detail.
type ParseError struct {
	Line    int
	Column  int
	Message string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error at %d:%d: %s", e.Line, e.Column, e.Message)
}

// ParseString parses Apache2 configuration text.
func ParseString(src string) (*Document, error) {
	return ParseReader(strings.NewReader(src))
}

// ParseReader parses Apache2 configuration from an io.Reader.
func ParseReader(r io.Reader) (*Document, error) {
	logicalLines, err := readLogicalLines(r)
	if err != nil {
		return nil, err
	}

	doc := &Document{}
	stack := []*Block{}

	appendStmt := func(stmt Statement) {
		if len(stack) == 0 {
			doc.Statements = append(doc.Statements, stmt)
			return
		}
		top := stack[len(stack)-1]
		top.Children = append(top.Children, stmt)
	}

	for _, ll := range logicalLines {
		trimmed := strings.TrimSpace(ll.Text)
		if trimmed == "" {
			continue
		}

		if isOpeningTag(trimmed) {
			name, args, perr := parseOpeningTag(trimmed, ll.Line)
			if perr != nil {
				return nil, perr
			}

			b := &Block{Name: name, Args: args, Pos: Position{Line: ll.Line, Column: ll.Column}}
			appendStmt(b)
			stack = append(stack, b)
			continue
		}

		if isClosingTag(trimmed) {
			name, perr := parseClosingTag(trimmed, ll.Line)
			if perr != nil {
				return nil, perr
			}
			if len(stack) == 0 {
				return nil, &ParseError{Line: ll.Line, Column: ll.Column, Message: fmt.Sprintf("unexpected closing tag </%s>", name)}
			}

			top := stack[len(stack)-1]
			if !strings.EqualFold(top.Name, name) {
				return nil, &ParseError{Line: ll.Line, Column: ll.Column, Message: fmt.Sprintf("mismatched closing tag </%s>, expected </%s>", name, top.Name)}
			}
			top.EndPos = Position{Line: ll.Line, Column: ll.Column}
			stack = stack[:len(stack)-1]
			continue
		}

		fields, perr := splitFieldsApache(trimmed, ll.Line)
		if perr != nil {
			return nil, perr
		}
		if len(fields) == 0 {
			continue
		}
		appendStmt(Directive{
			Name: fields[0],
			Args: fields[1:],
			Pos:  Position{Line: ll.Line, Column: ll.Column},
		})
	}

	if len(stack) > 0 {
		top := stack[len(stack)-1]
		return nil, &ParseError{Line: top.Pos.Line, Column: top.Pos.Column, Message: fmt.Sprintf("unclosed block <%s>", top.Name)}
	}

	return doc, nil
}

type logicalLine struct {
	Text   string
	Line   int
	Column int
}

func readLogicalLines(r io.Reader) ([]logicalLine, error) {
	scanner := bufio.NewScanner(r)
	physicalLine := 0
	out := make([]logicalLine, 0)

	var current strings.Builder
	currentStartLine := 0

	for scanner.Scan() {
		physicalLine++
		line := scanner.Text()

		lineNoComment := stripComments(line)
		lineNoComment = strings.TrimRightFunc(lineNoComment, unicode.IsSpace)

		if current.Len() == 0 {
			currentStartLine = physicalLine
		}

		if endsWithUnescapedBackslash(lineNoComment) {
			current.WriteString(strings.TrimSuffix(lineNoComment, "\\"))
			current.WriteString(" ")
			continue
		}

		current.WriteString(lineNoComment)
		out = append(out, logicalLine{
			Text:   current.String(),
			Line:   currentStartLine,
			Column: 1,
		})
		current.Reset()
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if current.Len() > 0 {
		out = append(out, logicalLine{
			Text:   current.String(),
			Line:   currentStartLine,
			Column: 1,
		})
	}

	return out, nil
}

func stripComments(line string) string {
	var out strings.Builder
	inSingle := false
	inDouble := false
	escaped := false

	for _, r := range line {
		if escaped {
			out.WriteRune(r)
			escaped = false
			continue
		}

		if r == '\\' {
			escaped = true
			out.WriteRune(r)
			continue
		}

		if r == '\'' && !inDouble {
			inSingle = !inSingle
			out.WriteRune(r)
			continue
		}

		if r == '"' && !inSingle {
			inDouble = !inDouble
			out.WriteRune(r)
			continue
		}

		if r == '#' && !inSingle && !inDouble {
			break
		}

		out.WriteRune(r)
	}

	return out.String()
}

func endsWithUnescapedBackslash(s string) bool {
	if s == "" {
		return false
	}

	backslashes := 0
	for i := len(s) - 1; i >= 0 && s[i] == '\\'; i-- {
		backslashes++
	}
	return backslashes%2 == 1
}

func isOpeningTag(s string) bool {
	return strings.HasPrefix(s, "<") && strings.HasSuffix(s, ">") && !strings.HasPrefix(s, "</")
}

func isClosingTag(s string) bool {
	return strings.HasPrefix(s, "</") && strings.HasSuffix(s, ">")
}

func parseOpeningTag(s string, line int) (string, []string, error) {
	inner := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(s, "<"), ">"))
	if inner == "" {
		return "", nil, &ParseError{Line: line, Column: 1, Message: "empty opening tag"}
	}

	fields, err := splitFieldsApache(inner, line)
	if err != nil {
		return "", nil, err
	}
	if len(fields) == 0 {
		return "", nil, &ParseError{Line: line, Column: 1, Message: "empty opening tag"}
	}

	return fields[0], fields[1:], nil
}

func parseClosingTag(s string, line int) (string, error) {
	inner := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(s, "</"), ">"))
	if inner == "" {
		return "", &ParseError{Line: line, Column: 1, Message: "empty closing tag"}
	}

	fields, err := splitFieldsApache(inner, line)
	if err != nil {
		return "", err
	}
	if len(fields) != 1 {
		return "", &ParseError{Line: line, Column: 1, Message: "invalid closing tag"}
	}

	return fields[0], nil
}

func splitFieldsApache(s string, line int) ([]string, error) {
	fields := make([]string, 0, 8)
	var cur strings.Builder
	inSingle := false
	inDouble := false
	escaped := false

	flush := func() {
		if cur.Len() == 0 {
			return
		}
		fields = append(fields, cur.String())
		cur.Reset()
	}

	for _, r := range s {
		if escaped {
			cur.WriteRune(r)
			escaped = false
			continue
		}

		if r == '\\' {
			escaped = true
			continue
		}

		if r == '\'' && !inDouble {
			inSingle = !inSingle
			continue
		}

		if r == '"' && !inSingle {
			inDouble = !inDouble
			continue
		}

		if unicode.IsSpace(r) && !inSingle && !inDouble {
			flush()
			continue
		}

		cur.WriteRune(r)
	}

	if escaped || inSingle || inDouble {
		return nil, &ParseError{Line: line, Column: 1, Message: "unterminated escape or quote"}
	}

	flush()
	return fields, nil
}
