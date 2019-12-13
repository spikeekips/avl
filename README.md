# avl

[![CircleCI](https://circleci.com/gh/spikeekips/avl/tree/master.svg?style=svg)](https://circleci.com/gh/spikeekips/avl/tree/master)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/spikeekips/avl)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fspikeekips%2Favl.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fspikeekips%2Favl?ref=badge_shield)
[![Go Report Card](https://goreportcard.com/badge/github.com/spikeekips/avl)](https://goreportcard.com/report/github.com/spikeekips/avl)
[![](https://tokei.rs/b1/github/spikeekips/avl?category=lines)](https://github.com/spikeekips/avl)

avl is simple [AVL Tree](https://en.wikipedia.org/wiki/AVL_tree) from scratch.

## `avl-print`

You can simply generate avl tree in command line.

```sh
$ go get github.com/spikeekips/avl/cmd/avl-print
```

```sh
$ avl-print -h
avl-print tree

Usage:
  avl-print [<key> ...] [flags]

Flags:
      --cpuprofile string       write cpu profile to file
  -h, --help                    help for avl-print
      --log string              log output directory
      --log-format log-format   log format: {json terminal} (default terminal)
      --log-level log-level     log level: {debug error warn info crit} (default error)
      --memprofile string       write memory profile to file
      --quiet                   no output
      --trace string            write trace to file
```

You can make the avl tree, which has 10 nodes.
```sh
$ seq 3 | avl-print 
graph graphname {
  2 [label="2 (1)"];
  2 -- 1;
  1 [label="1 (0)"];
  10 [label=" " style="filled" color="white" bgcolor="white"];
  1 -- 10 [style="solid" color="white" bgcolor="white"];
  11 [label=" " style="filled" color="white" bgcolor="white"];
  1 -- 11 [style="solid" color="white" bgcolor="white"];
  2 -- 3;
  3 [label="3 (0)"];
  30 [label=" " style="filled" color="white" bgcolor="white"];
  3 -- 30 [style="solid" color="white" bgcolor="white"];
  31 [label=" " style="filled" color="white" bgcolor="white"];
  3 -- 31 [style="solid" color="white" bgcolor="white"];
}
```

It generates the [dotgraph](https://www.google.com/search?q=dotgraph). You can draw it by external
tool, like `dot` of [graphviz](https://www.graphviz.org).

```
$ seq 100 | avl-print | dot -Tpng -o/tmp/100-node-avl-tree.png
```
![100-node-avl-tree.png](https://user-images.githubusercontent.com/174565/70210262-2977be00-172a-11ea-848c-ee00f330f1ee.png)


## UnitTest

```sh
$ go clean -testcache; go test -race -v ./ -run ..
```

> You can set the environment variable, `AVL_DEBUG=1` for enabling debug mode.


## TODO

* [ ] Performance test
* [ ] Code Documentation
* [ ] Examples
* [ ] Proof support
