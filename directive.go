package a2cp

// NewDirective creates a directive statement.
func NewDirective(name string, args ...string) Directive {
	return Directive{Name: name, Args: args}
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
