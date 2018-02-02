# logictree

`golang` library to construct and evaluate text-based template trees.

## Install

```
go get github.com/sabhiram/logictree
```

## Components

The `tree` above is composed of nodes and leaves. 

A `Leaf` is a string expression containing a "truthy" statement (evaluates to `true` or `false` when the returned template is executed with a context).

A `Node` is the entity that combines a logical operator with its `children`.  The only supported node types are:
1. Leaf
2. And
3. Or

## How it works

When `Combine` is called at any given `Node`, it recurses down the tree and combines all sub-trees into an evaluate-able string.  Alternatively the caller may use the `Node`'s `GetTemplate` method to return a `*template.Template` version of the string which can be executed against various dynamic contexts for filtering, event monitoring and so on.


## Usage

The idea here is to build a tree that represents some arbitrary grouping of logical statements, which when executed against a context of values will evaluate to `true` or `false`.  This is useful for various if-this-then-that-esq scenarios.  Here is one such example:

Lets say we want to make a tree so that we can evaluate the following statement in a programmatic fashion:
"If the price of milk is between 4 and 6 a gallon when the price of onions is between 1 and 2 per pound, or if the price of toothpaste is more than 5; do something!"

This breaks down to the following in code (assume `m` is the price of milk, `o` the price of onions and `t` the price of toothpaste):

```
if (((m >= 4 && m <= 6) && (o >= 1 && o <= 2) || (t > 5))
   then -> doSomething();
```
    
Here is the equivalent logic tree being constructed:

```
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
    milkTree := logictree.NewNode("and",
        logictree.NewLeafNode("ge .Milk 4"),
        logictree.NewLeafNode("le .Milk 6"))

    // Now one for the onions.
    onionTree := logictree.NewNode("and",
        logictree.NewLeafNode("ge .Onions 1"),
        logictree.NewLeafNode("le .Onions 2"))

    // I think you see how this works, lets build the whole tree!
    tree := logictree.NewNode("or",
        logictree.NewNode("and", milkTree, onionTree),
        logictree.NewLeafNode("gt .Toothpaste 5"))

    // Here is the expression for the tree before it has been templateized.
    expr, err := tree.Combine()
    fatalOnError(err)
    fmt.Printf("Tree Expression: \"%s\"\n", expr)

    // Grab the template so we can execute it!
    t, err := tree.GetTemplate(nil)
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
```

Running this generates:
```
Tree Expression: "or (and (and ((ge .Milk 4)) ((le .Milk 6))) (and ((ge .Onions 1)) ((le .Onions 2)))) ((gt .Toothpaste 5))"
Result for main.Prices{Milk:5, Onions:0, Toothpaste:4} ==> false
Result for main.Prices{Milk:5, Onions:2, Toothpaste:4} ==> true
Result for main.Prices{Milk:5, Onions:0, Toothpaste:8} ==> true
```

## JSON Marshal / Unmarshal

If you would like to express your logic as JSON, the `*Node` is capable of being serialized / deserialized.

```
    //
    //  Now with JSON powers
    //

    // Write Tree -> JSON
    bs, err := json.MarshalIndent(tree, "", "  ")
    fatalOnError(err)
    fmt.Printf("Tree in JSON:\n%s\n", bs)

    // Read Tree <- JSON
    newTree := &logictree.Node{}
    err = json.Unmarshal(bs, &newTree)
    fatalOnError(err)
```

## Custom functions for your templates

This is not really a feature of `logictree`, but you can pass a `template.FuncMap` to the `*node.GetTemplate` which allows us to define custom functions.  Here is a simple example where we replace the two `and` trees using a custom `between` function.

```
    //
    //  Define your own operators for the tree by using custom
    //  `template.FuncMap`s.
    //
    mt2 := logictree.NewLeafNode("between .Milk 4 6")
    ot2 := logictree.NewLeafNode("between .Onions 1 2")
    tree2 := logictree.NewNode("or",
        logictree.NewNode("and", mt2, ot2),
        logictree.NewLeafNode("gt .Toothpaste 5"))

    // Here is the expression for the tree before it has been templateized.
    expr, err = tree2.Combine()
    fatalOnError(err)
    fmt.Printf("Tree2 Expression: \"%s\"\n", expr)

    // Since we are using a custom function `between`, teach the template
    // evaluator what it means to use this operator.
    fm := template.FuncMap{
        "between": func(v, mi, ma int) string {
            if v >= mi && v <= ma {
                return "true"
            }
            return "false"
        },
    }
    t2, err := mt2.GetTemplate(fm)
    fatalOnError(err)

    // Now we can execute a template with the `between` function!
    buf.Reset()
    t2.Execute(&buf, &p)
    fmt.Printf("Result for %#v ==> %v\n", p, buf.String())
}
```

## TODO

1. Better ways to define the `tree` from an external caller.
2. Other operators? Not all of them will apply to N leaves etc?
3. Do we want to provide any custom helper functions to aid in the template execution?
4. How does the caller specify a `FuncMap` for the template?
5. Ops should be user defined? Call register op?
