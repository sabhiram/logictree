# logictree

`golang` library to construct and evaluate text-based template trees.

## Install

```
go get github.com/sabhiram/logictree
```

## Usage

```
    tree := &logictree.Node{
        op: logictree.OperatorAnd,
        leaves: []logictree.TreeMerger{
            logictree.NewLeaf("gt 1 0"),
            logictree.NewLeaf("gt 2 0"),
            logictree.NewLeaf("gt 3 0"),
            logictree.NewLeaf("gt 4 2"),
            &logictree.Node{
                op: logictree.OperatorOr,
                leaves: []logictree.TreeMerger{
                    logictree.NewLeaf("gt 1 10"),
                    logictree.NewLeaf("gt 2 10"),
                    logictree.NewLeaf("gt 3 10"),
                    logictree.NewLeaf("gt 40 2"),
                },
            },
        },
    }

    tmpl, err := tree.GetTemplate()
    if err != nil {
        t.Errorf("GetTemplate() failed with error: %s\n", err.Error())
    }

    tmpl.Execute(os.Stdout, nil)
```

## Components

The `tree` above is composed of `Node`s and `Leaf`s.

A `Node` contains an operator which merges the `leaves` contained within the said `Node`.  Each entry in the `leaves` array can either be a terminating node (`Leaf`), or a plain `Node`.

A `Leaf` is a terminating node of the tree which only contains an expression (which will later be templatized and possibly evaluated).

## How it works

When `Merge` is called at any given `Node`, it recursively invokes the `Merge` method for all `leaves` in the `Node` and then combines them using the specified operator.  All `Leaf` and `Node` instances implement a `Merge` method making them compatible with the `TreeMerger` interface.

## TODO

1. Better ways to define the `tree` from an external caller.
2. Other operators? Not all of them will apply to N leaves etc?
3. Do we want to provide any custom helper functions to aid in the template execution?
4. How does the caller specify a `FuncMap` for the template?
