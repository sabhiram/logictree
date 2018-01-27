// Package logictree allows the building and evaluation of text template based
// truthy trees.
package logictree

////////////////////////////////////////////////////////////////////////////////

import (
	"errors"
	"fmt"
	"text/template"
)

////////////////////////////////////////////////////////////////////////////////

var (
	ErrEmptyNode = errors.New("empty node cannot be merged")
)

////////////////////////////////////////////////////////////////////////////////

// Tree is an interface that all nodes and leaves adhere to so we can
// build complicated logic trees.
type Tree interface {
	Combine() (string, error)
}

////////////////////////////////////////////////////////////////////////////////

// Operator defines how leaves are combined in a tree.
type Operator int

const (
	OperatorAnd = iota
	OperatorOr  = iota
)

func (o Operator) String() string {
	switch o {
	case OperatorAnd:
		return "and"
	case OperatorOr:
		return "or"
	default:
		panic("invalid operator type")
	}
}

// Apply combines the number of `exprs` into a evaluate-able string combining
// the expressions using the specified operator.
func (o Operator) Apply(exprs []string) string {
	switch len(exprs) {
	case 0:
		return ""
	case 1:
		return exprs[0]
	case 2:
		return fmt.Sprintf("%s (%s) (%s)", o.String(), exprs[0], exprs[1])
	}
	return fmt.Sprintf("%s (%s) (%s)", o.String(), exprs[0], o.Apply(exprs[1:]))
}

////////////////////////////////////////////////////////////////////////////////

// Leaf is a node in a tree with no children, and contains an
// evaluate-by-template string.
type Leaf struct {
	expression string
}

// NewLeaf returns an instance to a leaf wrapped with a scope operator.
// Possibly overkill.
func NewLeaf(expr string) *Leaf {
	return &Leaf{
		expression: "(" + expr + ")",
	}
}

// Combine implements the Tree interface.
func (l *Leaf) Combine() (string, error) {
	return l.expression, nil
}

////////////////////////////////////////////////////////////////////////////////

// Node is the generic node in a tree which combines a bunch of child nodes
// using it's specific operator.
type Node struct {
	op       Operator
	children []Tree
}

// NewNode returns a sub-tree which represents the combination of the `op` with
// the child sub-trees.
func NewNode(op Operator, cs ...Tree) *Node {
	return &Node{
		op:       op,
		children: cs,
	}
}

// Combine is required to satisfy the Tree interface.
func (n *Node) Combine() (string, error) {
	if len(n.children) == 0 {
		return "", ErrEmptyNode
	}

	exprs := []string{}
	for _, tm := range n.children {
		e, err := tm.Combine()
		if err != nil {
			return "", err
		}
		exprs = append(exprs, e)
	}

	return n.op.Apply(exprs), nil
}

// GetTemplate squashes the tree down from the root down into a single template
// expression.
func (n *Node) GetTemplate() (*template.Template, error) {
	e, err := n.Combine()
	if err != nil {
		return nil, err
	}

	return template.Must(template.New("tree").Parse("{{ " + e + " }}")), nil
}
