package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/sabhiram/logictree"
)

func fatalOnError(err error) {
	if err != nil {
		fmt.Printf("Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}

func main() {
	// Lets say we want to make a tree so that we can evaluate the following
	// statement in a programmatic fashion:
	// "If the price of milk is between 4 and 6 a gallon when the price of
	// onions is between 1 and 2 per pound, or if the price of toothpaste is
	// more than 5; do something!"

	// This breaks down to the following in code (assume `m` is the price of
	// milk, `o` the price of onions and `t` the price of toothpaste):
	// if (((m >= 4 && m <= 6) && (o >= 1 && o <= 2) || (t > 5))
	//    then -> doSomething()

	// Ok so we need prices for things, and some place to put them.
	type Prices struct {
		Milk       int
		Onions     int
		Toothpaste int
	}

	// Lets build a sub-tree for the `milk` portion of our statement.
	milkTree := logictree.NewNode(logictree.OperatorAnd,
		logictree.NewLeaf("ge .Milk 4"),
		logictree.NewLeaf("le .Milk 6"))

	// Now one for the onions.
	onionTree := logictree.NewNode(logictree.OperatorAnd,
		logictree.NewLeaf("ge .Onions 1"),
		logictree.NewLeaf("le .Onions 2"))

	// I think you see how this works, lets build the whole tree!
	tree := logictree.NewNode(logictree.OperatorOr,
		logictree.NewNode(logictree.OperatorAnd, milkTree, onionTree),
		logictree.NewLeaf("gt .Toothpaste 5"))

	// Here is the expression for the tree before it has been templateized.
	expr, err := tree.Combine()
	fatalOnError(err)

	fmt.Printf("Tree Expression: \"%s\"\n", expr)

	// Grab the template so we can execute it!
	t, err := tree.GetTemplate()
	fatalOnError(err)

	// Setup some prices - this is expected to evaluate to false.
	p := Prices{
		Milk:       5,
		Onions:     0,
		Toothpaste: 4,
	}

	var buf bytes.Buffer
	t.Execute(&buf, &p)
	fmt.Printf("Result for %#v ==> %v\n", p, buf.String())

	// Lets make the statement true.
	p.Onions = 2

	buf.Reset()
	t.Execute(&buf, &p)
	fmt.Printf("Result for %#v ==> %v\n", p, buf.String())

	// Lets try it the other way.
	p.Onions = 0
	p.Toothpaste = 8

	buf.Reset()
	t.Execute(&buf, &p)
	fmt.Printf("Result for %#v ==> %v\n", p, buf.String())
}
