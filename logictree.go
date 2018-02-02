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

// Operator defines how leaves are combined in a tree.
type Operator string

const (
	OperatorLeaf = "leaf"
	OperatorAnd  = "and"
	OperatorOr   = "or"
)

func (o Operator) String() string {
	switch o {
	case OperatorLeaf, OperatorAnd, OperatorOr:
		return string(o)
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

// Node is the generic node in a tree which combines a bunch of child nodes
// using it's specific operator.
type Node struct {
	Op    Operator `json:"Op"`
	Nodes []*Node  `json:"Nodes,omitempty"`
	Leaf  string   `json:"Leaf,omitempty"`
}

// NewNode returns a sub-tree which represents the combination of the `op` with
// the child sub-trees.
func NewNode(op Operator, cs ...*Node) *Node {
	return &Node{
		Op:    op,
		Nodes: cs,
		Leaf:  "",
	}
}

// NewLeafNode returns a new leaf node.
func NewLeafNode(expr string) *Node {
	return &Node{
		Op:    OperatorLeaf,
		Nodes: nil,
		Leaf:  "(" + expr + ")",
	}
}

// Combine merges this node with any of its children (evaluated).
func (n *Node) Combine() (string, error) {
	// If we are a leaf node, we just return our expression.
	if n.Op == OperatorLeaf {
		return n.Leaf, nil
	}

	if len(n.Nodes) == 0 {
		return "", ErrEmptyNode
	}

	exprs := []string{}
	for _, tm := range n.Nodes {
		e, err := tm.Combine()
		if err != nil {
			return "", err
		}
		exprs = append(exprs, e)
	}

	return n.Op.Apply(exprs), nil
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
