package share

import (
	"errors"
)

type Tree struct {
	parent   *Tree
	children map[string]*Tree
	clocks   Clocks
	value    string
}

type Leaf struct {
	Clocks string
	Value  string
}

type Leaves map[string]Leaf

func (p *Tree) GetLeaves(prefix string) (leaves Leaves) {
	leaves = make(Leaves)
	p.leaves(prefix, leaves)
	return
}

func (p *Tree) leaves(prefix string, result Leaves) {
	if len(p.clocks) != 0 {
		result[prefix] = Leaf{p.clocks.ToString(), p.value}
	}
	for k, v := range p.children {
		path := k
		if len(prefix) != 0 {
			path = prefix + "/" + k
		}
		v.leaves(path, result)
	}
}

func NewTree() *Tree {
	return &Tree{
		children: make(map[string]*Tree),
		clocks:   NewClocks(),
	}
}

func (t *Tree) Get(name string) *Tree {
	child, ok := t.children[name]
	if !ok {
		child = &Tree{
			parent:   t,
			children: make(map[string]*Tree),
			clocks:   NewClocks(),
		}
		t.children[name] = child
	}
	return child
}

func (p *Tree) Put(clocks Clocks, value string) error {
	status := clocks.Compare(p.clocks)
	if status != Greater && status != Equal {
		return errors.New(CompareResultString(status))
	}
	p.value = value
	clocks.Absorb(p.clocks)
	p.clocks = clocks
	return nil
}

func (t *Tree) Value() string {
	return t.value
}

func (t *Tree) Clocks() Clocks {
	return t.clocks
}
