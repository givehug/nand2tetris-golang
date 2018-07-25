package parsetree

// ParseTree ...
type ParseTree struct {
	t        string       // leaf type string constant
	value    string       // leaf value
	children []*ParseTree // leaf children
}

// New constructs new ParseTree
func New(t string, v string) *ParseTree {
	var c []*ParseTree
	return &ParseTree{t, v, c}
}

// AddChildren ...
func (p *ParseTree) AddChildren(c ...*ParseTree) {
	p.children = append(p.children, c...)
}

// Children ...
func (p *ParseTree) Children() []*ParseTree {
	return p.children
}

// HasChildren ...
func (p *ParseTree) HasChildren() bool {
	return len(p.children) > 0
}

// SetValue ...
func (p *ParseTree) SetValue(s string) {
	p.value = s
}

// Value ...
func (p *ParseTree) Value() string {
	return p.value
}

// SetType ...
func (p *ParseTree) SetType(s string) {
	p.t = s
}

// Type ...
func (p *ParseTree) Type() string {
	return p.t
}

// TraversePreorder performs depth first ParseTree traversal
func (p *ParseTree) TraversePreorder(fn func(l *ParseTree)) {
	fn(p)
	for _, leaf := range p.Children() {
		leaf.TraversePreorder(fn)
	}
}
