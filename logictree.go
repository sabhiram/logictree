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

// TreeMerger is an interface that all nodes and leaves adhere to so we can
// build complicated logic trees.
type TreeMerger interface {
	Merge() (string, error)
}

////////////////////////////////////////////////////////////////////////////////

// A leaf is the evaluate-able expression of which there are 1 or more in a
// tree.
type Leaf struct {
	expression string
}

func NewLeaf(expr string) *Leaf {
	return &Leaf{
		expression: "(" + expr + ")",
	}
}

// Merge implements the TreeMerger interface.
func (l *Leaf) Merge() (string, error) {
	return l.expression, nil
}

////////////////////////////////////////////////////////////////////////////////

// Operator defines how leaves are combined in a tree
type Operator int

const (
	cOperatorAnd = iota
	cOperatorOr  = iota
)

func (o Operator) String() string {
	switch o {
	case cOperatorAnd:
		return "and"
	case cOperatorOr:
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

// Node defines the evaluate-able tree.
type Node struct {
	leaves []TreeMerger
	op     Operator
}

func (n *Node) Merge() (string, error) {
	if len(n.leaves) == 0 {
		return "", ErrEmptyNode
	}

	exprs := []string{}
	for _, tm := range n.leaves {
		e, err := tm.Merge()
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
	e, err := n.Merge()
	if err != nil {
		return nil, err
	}

	return template.Must(template.New("tree").Parse("{{ " + e + " }}")), nil
}

////////////////////////////////////////////////////////////////////////////////
