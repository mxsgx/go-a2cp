package a2cp

import "testing"

func TestNewDirective(t *testing.T) {
	d := NewDirective("Listen", "80")

	if d.Name != "Listen" {
		t.Fatalf("Name = %q, want %q", d.Name, "Listen")
	}
	if len(d.Args) != 1 || d.Args[0] != "80" {
		t.Fatalf("Args = %#v, want %#v", d.Args, []string{"80"})
	}
}

func TestAsDirective(t *testing.T) {
	dv, ok := asDirective(Directive{Name: "ServerName", Args: []string{"example.com"}})
	if !ok {
		t.Fatalf("asDirective(Directive) ok = false, want true")
	}
	if dv.Name != "ServerName" {
		t.Fatalf("asDirective(Directive).Name = %q, want %q", dv.Name, "ServerName")
	}

	dp, ok := asDirective(&Directive{Name: "DocumentRoot", Args: []string{"/var/www"}})
	if !ok {
		t.Fatalf("asDirective(*Directive) ok = false, want true")
	}
	if dp.Name != "DocumentRoot" {
		t.Fatalf("asDirective(*Directive).Name = %q, want %q", dp.Name, "DocumentRoot")
	}

	if _, ok := asDirective(&Block{Name: "VirtualHost"}); ok {
		t.Fatalf("asDirective(*Block) ok = true, want false")
	}
}
