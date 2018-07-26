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

// AddLeaves ...
func (p *ParseTree) AddLeaves(c ...*ParseTree) {
	p.children = append(p.children, c...)
}

// Leaves ...
func (p *ParseTree) Leaves() []*ParseTree {
	return p.children
}

// HasLeaves ...
func (p *ParseTree) HasLeaves() bool {
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
	for _, leaf := range p.Leaves() {
		leaf.TraversePreorder(fn)
	}
}
