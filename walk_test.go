package a2cp

import (
	"reflect"
	"strconv"
	"testing"
)

func TestSingleLevelFindRemainsNonRecursive(t *testing.T) {
	doc := NewDocument()
	doc.AddDirective("ServerName", "root.example")
	vh := doc.AddBlock("VirtualHost", "*:80")
	vh.AddDirective("ServerName", "nested.example")

	if got := len(doc.FindDirectives("servername")); got != 1 {
		t.Fatalf("FindDirectives(servername) = %d, want 1", got)
	}
	if got := len(doc.FindDirectivesRecursive("servername")); got != 2 {
		t.Fatalf("FindDirectivesRecursive(servername) = %d, want 2", got)
	}

	outer := doc.AddBlock("Directory", "/var/www")
	outer.AddBlock("Directory", "/var/www/nested")

	if got := len(doc.FindBlocks("directory")); got != 1 {
		t.Fatalf("FindBlocks(directory) = %d, want 1", got)
	}
	if got := len(doc.FindBlocksRecursive("directory")); got != 2 {
		t.Fatalf("FindBlocksRecursive(directory) = %d, want 2", got)
	}
}

func TestRecursiveFindMultiLevel(t *testing.T) {
	doc := NewDocument()
	root := doc.AddBlock("VirtualHost", "*:443")
	root.AddDirective("ServerName", "root.example")

	level1 := root.AddBlock("Directory", "/srv/app")
	level1.AddDirective("ServerName", "level1.example")

	level2 := level1.AddBlock("Directory", "/srv/app/private")
	level2.AddDirective("servername", "level2.example")

	if got := len(doc.FindDirectivesRecursive("SERVERNAME")); got != 3 {
		t.Fatalf("FindDirectivesRecursive(SERVERNAME) = %d, want 3", got)
	}
	if got := len(doc.FindBlocksRecursive("directory")); got != 2 {
		t.Fatalf("FindBlocksRecursive(directory) = %d, want 2", got)
	}

	if got := len(root.FindDirectivesRecursive("servername")); got != 3 {
		t.Fatalf("Block FindDirectivesRecursive(servername) = %d, want 3", got)
	}
	if got := len(root.FindBlocksRecursive("directory")); got != 2 {
		t.Fatalf("Block FindBlocksRecursive(directory) = %d, want 2", got)
	}
}

func TestWalkStopsEarly(t *testing.T) {
	doc := NewDocument()
	doc.AddDirective("A")
	b := doc.AddBlock("B")
	b.AddDirective("C")
	doc.AddDirective("D")

	visited := make([]string, 0)
	doc.Walk(func(stmt Statement, depth int) bool {
		visited = append(visited, stmtLabel(stmt, depth))
		return len(visited) < 2
	})

	want := []string{"D:A@0", "B:B@0"}
	if !reflect.DeepEqual(visited, want) {
		t.Fatalf("visited = %#v, want %#v", visited, want)
	}
}

func TestWalkDepthFirstPreOrderMixedTree(t *testing.T) {
	doc := NewDocument()
	doc.AddDirective("Global")
	vh := doc.AddBlock("VirtualHost", "*:80")
	vh.AddDirective("ServerName", "example.com")
	dir := vh.AddBlock("Directory", "/var/www/html")
	dir.AddDirective("Require", "all", "granted")
	doc.AddDirective("Tail")

	visited := make([]string, 0)
	doc.Walk(func(stmt Statement, depth int) bool {
		visited = append(visited, stmtLabel(stmt, depth))
		return true
	})

	want := []string{
		"D:Global@0",
		"B:VirtualHost@0",
		"D:ServerName@1",
		"B:Directory@1",
		"D:Require@2",
		"D:Tail@0",
	}
	if !reflect.DeepEqual(visited, want) {
		t.Fatalf("visited = %#v, want %#v", visited, want)
	}

	blockVisited := make([]string, 0)
	vh.Walk(func(stmt Statement, depth int) bool {
		blockVisited = append(blockVisited, stmtLabel(stmt, depth))
		return true
	})

	blockWant := []string{
		"B:VirtualHost@0",
		"D:ServerName@1",
		"B:Directory@1",
		"D:Require@2",
	}
	if !reflect.DeepEqual(blockVisited, blockWant) {
		t.Fatalf("block visited = %#v, want %#v", blockVisited, blockWant)
	}
}

func stmtLabel(stmt Statement, depth int) string {
	switch s := stmt.(type) {
	case Directive:
		return "D:" + s.Name + "@" + strconv.Itoa(depth)
	case *Directive:
		return "D:" + s.Name + "@" + strconv.Itoa(depth)
	case *Block:
		return "B:" + s.Name + "@" + strconv.Itoa(depth)
	default:
		return "?@" + strconv.Itoa(depth)
	}
}
