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

// Container represents a node that owns statement slices.
type Container interface {
	statements() []Statement
	setStatements([]Statement)
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
	Name       string
	Args       []string
	Children   []Statement
	Parent     Container
	Pos        Position
	EndPos     Position
	EndComment string
}

func (Block) isStatement() {}

// Document is the parsed representation of a .conf file.
type Document struct {
	Statements []Statement
}

func setParents(doc *Document) {
	if doc == nil {
		return
	}
	setParentOnStatements(doc, doc.Statements)
}

func setParentOnStatements(parent Container, stmts []Statement) {
	for _, stmt := range stmts {
		block, ok := stmt.(*Block)
		if !ok {
			continue
		}
		block.Parent = parent
		setParentOnStatements(block, block.Children)
	}
}

func attachParent(parent Container, stmt Statement) {
	block, ok := stmt.(*Block)
	if !ok {
		return
	}
	block.Parent = parent
}
