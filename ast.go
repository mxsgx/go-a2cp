package a2cp

// Position identifies a location in the source config.
type Position struct {
	Line   int
	Column int
}

// Statement is implemented by all Apache config AST nodes.
type Statement interface {
	isStatement()
}

// Directive represents a single config directive, e.g. `Listen 80`.
type Directive struct {
	Name string
	Args []string
	Pos  Position
}

func (Directive) isStatement() {}

// Comment represents a config comment, without the leading #.
type Comment struct {
	Text string
	Pos  Position
}

func (Comment) isStatement() {}

// Block represents a container directive, e.g. `<Directory /var/www> ... </Directory>`.
type Block struct {
	Name     string
	Args     []string
	Children []Statement
	Pos      Position
	EndPos   Position
}

func (Block) isStatement() {}

// Document is the parsed representation of a .conf file.
type Document struct {
	Statements []Statement
}
