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
