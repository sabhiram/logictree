package logictree

////////////////////////////////////////////////////////////////////////////////

import (
	"fmt"
	"os"
	"testing"
)

////////////////////////////////////////////////////////////////////////////////

func TestLeafMerge(t *testing.T) {
	for _, tc := range []struct {
		expr     string
		expected string
	}{
		{"1", "(1)"},
		{"a and b", "(a and b)"},
	} {
		l := NewLeaf(tc.expr)
		e, err := l.Merge()
		if err != nil {
			t.Errorf("Leaf::Merge() error: %s\n", err.Error())
		}

		if e != tc.expected {
			t.Errorf("Leaf::Merge() expected=%s actual=%s\n", tc.expected, e)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////

func TestTreeConstruction(t *testing.T) {
	tree := &Node{
		op: cOperatorAnd,
		leaves: []TreeMerger{
			NewLeaf("gt 1 0"),
			NewLeaf("gt 2 0"),
			NewLeaf("gt 3 0"),
			NewLeaf("gt 4 2"),
			&Node{
				op: cOperatorOr,
				leaves: []TreeMerger{
					NewLeaf("gt 1 10"),
					NewLeaf("gt 2 10"),
					NewLeaf("gt 3 10"),
					NewLeaf("gt 40 2"),
				},
			},
		},
	}

	s, err := tree.Merge()
	if err != nil {
		t.Errorf("Merge() failed with error: %s\n", err.Error())
	}
	fmt.Printf("MERGE: %s\n", s)

	tmpl, err := tree.GetTemplate()
	if err != nil {
		t.Errorf("GetTemplate() failed with error: %s\n", err.Error())
	}

	fmt.Printf("Template: %#v\n", tmpl)
	tmpl.Execute(os.Stdout, nil)
}

////////////////////////////////////////////////////////////////////////////////
