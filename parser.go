package a2cp

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

// ParseOption configures parser behavior.
type ParseOption func(*parseOptions)

type parseOptions struct {
	includeResolution bool
	basePath          string
}

type includeState struct {
	inProgress map[string]struct{}
}

func defaultParseOptions() parseOptions {
	return parseOptions{}
}

func applyParseOptions(opts []ParseOption) parseOptions {
	cfg := defaultParseOptions()
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	return cfg
}

// WithIncludeResolution enables recursive resolution of Include and IncludeOptional directives.
//
// Relative include paths are resolved from basePath.
func WithIncludeResolution(basePath string) ParseOption {
	return func(cfg *parseOptions) {
		cfg.includeResolution = true
		cfg.basePath = basePath
	}
}

// ParseString parses Apache2 configuration text.
func ParseString(src string) (*Document, error) {
	return ParseReader(strings.NewReader(src))
}

// ParseReader parses Apache2 configuration from an io.Reader.
func ParseReader(r io.Reader, opts ...ParseOption) (*Document, error) {
	cfg := applyParseOptions(opts)
	state := &includeState{inProgress: make(map[string]struct{})}
	baseDir := "."
	if cfg.basePath != "" {
		baseDir = cfg.basePath
	}
	return parseReaderWithContext(r, baseDir, cfg, state)
}

func parseReaderWithContext(r io.Reader, baseDir string, cfg parseOptions, state *includeState) (*Document, error) {
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
			if ll.HasComment {
				appendStmt(Comment{
					Text: ll.Comment,
					Pos:  Position{Line: ll.Line, Column: ll.CommentColumn},
				})
			}
			continue
		}

		if isOpeningTag(trimmed) {
			name, args, perr := parseOpeningTag(trimmed, ll.Line)
			if perr != nil {
				return nil, perr
			}

			b := &Block{Name: name, Args: args, Pos: Position{Line: ll.Line, Column: ll.Column}}
			appendStmt(b)
			if ll.HasComment {
				appendStmt(Comment{
					Text: ll.Comment,
					Pos:  Position{Line: ll.Line, Column: ll.CommentColumn},
				})
			}
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
			if ll.HasComment {
				appendStmt(Comment{
					Text: ll.Comment,
					Pos:  Position{Line: ll.Line, Column: ll.CommentColumn},
				})
			}
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
		if ll.HasComment {
			appendStmt(Comment{
				Text: ll.Comment,
				Pos:  Position{Line: ll.Line, Column: ll.CommentColumn},
			})
		}
	}

	if len(stack) > 0 {
		top := stack[len(stack)-1]
		return nil, &ParseError{Line: top.Pos.Line, Column: top.Pos.Column, Message: fmt.Sprintf("unclosed block <%s>", top.Name)}
	}

	if cfg.includeResolution {
		resolved, err := resolveIncludeStatements(doc.Statements, baseDir, cfg, state)
		if err != nil {
			return nil, err
		}
		doc.Statements = resolved
	}

	return doc, nil
}

func resolveIncludeStatements(stmts []Statement, baseDir string, cfg parseOptions, state *includeState) ([]Statement, error) {
	out := make([]Statement, 0, len(stmts))

	for _, stmt := range stmts {
		switch s := stmt.(type) {
		case Directive:
			if strings.EqualFold(s.Name, "Include") || strings.EqualFold(s.Name, "IncludeOptional") {
				included, err := resolveIncludeDirective(s, baseDir, cfg, state)
				if err != nil {
					return nil, err
				}
				out = append(out, included...)
				continue
			}
			out = append(out, s)
		case *Block:
			children, err := resolveIncludeStatements(s.Children, baseDir, cfg, state)
			if err != nil {
				return nil, err
			}
			s.Children = children
			out = append(out, s)
		default:
			out = append(out, stmt)
		}
	}

	return out, nil
}

func resolveIncludeDirective(d Directive, baseDir string, cfg parseOptions, state *includeState) ([]Statement, error) {
	if len(d.Args) == 0 {
		return nil, fmt.Errorf("directive %s requires at least one path", d.Name)
	}

	isOptional := strings.EqualFold(d.Name, "IncludeOptional")
	out := make([]Statement, 0)

	for _, rawPattern := range d.Args {
		pattern := rawPattern
		if !filepath.IsAbs(pattern) {
			pattern = filepath.Join(baseDir, pattern)
		}

		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf("resolve include pattern %q: %w", rawPattern, err)
		}

		if len(matches) == 0 {
			if isOptional {
				continue
			}
			return nil, fmt.Errorf("include path not found for pattern %q", rawPattern)
		}

		for _, match := range matches {
			doc, err := parseIncludedFile(match, cfg, state)
			if err != nil {
				return nil, err
			}
			out = append(out, doc.Statements...)
		}
	}

	return out, nil
}

func parseIncludedFile(path string, cfg parseOptions, state *includeState) (*Document, error) {
	canonical, err := canonicalPath(path)
	if err != nil {
		return nil, err
	}

	if _, exists := state.inProgress[canonical]; exists {
		return nil, fmt.Errorf("circular include detected: %s", canonical)
	}

	state.inProgress[canonical] = struct{}{}
	defer delete(state.inProgress, canonical)

	f, err := os.Open(canonical)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	baseDir := filepath.Dir(canonical)
	if cfg.basePath != "" {
		baseDir, err = canonicalPath(cfg.basePath)
		if err != nil {
			return nil, err
		}
	}
	return parseReaderWithContext(f, baseDir, cfg, state)
}

func canonicalPath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	resolved, err := filepath.EvalSymlinks(abs)
	if err == nil {
		return resolved, nil
	}

	if os.IsNotExist(err) {
		return abs, nil
	}

	return "", err
}

type logicalLine struct {
	Text          string
	Comment       string
	HasComment    bool
	CommentColumn int
	Line          int
	Column        int
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

		inSingle, inDouble := quoteStateForContinuation(current.String())
		lineNoComment, comment, commentColumn, hasComment := splitCodeAndComment(line, inSingle, inDouble)
		lineNoComment = strings.TrimRightFunc(lineNoComment, unicode.IsSpace)

		if current.Len() == 0 {
			currentStartLine = physicalLine
		}

		if endsWithUnescapedBackslash(lineNoComment) {
			if hasComment {
				out = append(out, logicalLine{
					Text:          "",
					Comment:       comment,
					HasComment:    true,
					CommentColumn: commentColumn,
					Line:          physicalLine,
					Column:        1,
				})
			}
			fragment := strings.TrimRightFunc(strings.TrimSuffix(lineNoComment, "\\"), unicode.IsSpace)
			if current.Len() > 0 {
				fragment = strings.TrimLeftFunc(fragment, unicode.IsSpace)
			}
			current.WriteString(fragment)
			current.WriteString(" ")
			continue
		}

		fragment := lineNoComment
		if current.Len() > 0 {
			fragment = strings.TrimLeftFunc(fragment, unicode.IsSpace)
		}
		current.WriteString(fragment)
		out = append(out, logicalLine{
			Text:          current.String(),
			Comment:       comment,
			HasComment:    hasComment,
			CommentColumn: commentColumn,
			Line:          currentStartLine,
			Column:        1,
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

func quoteStateForContinuation(s string) (bool, bool) {
	inSingle := false
	inDouble := false
	escaped := false

	for _, r := range s {
		if escaped {
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
	}

	return inSingle, inDouble
}

func splitCodeAndComment(line string, inSingle, inDouble bool) (string, string, int, bool) {
	var out strings.Builder
	escaped := false
	column := 0

	for i, r := range line {
		column++

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
			comment := strings.TrimRightFunc(line[i+1:], unicode.IsSpace)
			return out.String(), comment, column, true
		}

		out.WriteRune(r)
	}

	return out.String(), "", 0, false
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
