# avl

[![CircleCI](https://circleci.com/gh/spikeekips/avl/tree/master.svg?style=svg)](https://circleci.com/gh/spikeekips/avl/tree/master)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/spikeekips/avl)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fspikeekips%2Favl.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fspikeekips%2Favl?ref=badge_shield)
[![Go Report Card](https://goreportcard.com/badge/github.com/spikeekips/avl)](https://goreportcard.com/report/github.com/spikeekips/avl)
[![](https://tokei.rs/b1/github/spikeekips/avl?category=lines)](https://github.com/spikeekips/avl)

avl is simple [AVL Tree](https://en.wikipedia.org/wiki/AVL_tree) from scratch.

## UnitTest

```sh
$ go clean -testcache; go test -race -v ./ -run ..
```

> You can set the environment variable, `AVL_DEBUG=1` for enabling debug mode.

## `avl2dot`

`avl2dot` is simple utility to generate dot graph from `TreeGenerator`. You can install `avl2dot`.

```sh
$ go get github.com/spikeekips/avl/cmd/avl2dot
```

You can pass the nodes keys to avl2dot, it will print dot graph.
```sh
$ avl2dot my name is spike ekips i am developer
```

![avl2dot-example-08bb.png](https://user-images.githubusercontent.com/174565/70696864-6e689b00-1cbc-11ea-90da-4d030a8ba2ad.png)

> The nodes key can be passed by standard input.

## Performance

### Elapsed Time By Number Of Nodes

![TreeGenerator-elapsedtime Plot](https://user-images.githubusercontent.com/174565/70695519-cfdb3a80-1cb9-11ea-9885-59f2b3012a72.png)

> Estimated by `avl2dot` at 08bbcafe2359356bf0da02a1177635855f66de8d .

### golang Benchmark

```sh
$ go clean -testcache; go test -race -v -run _ -bench BenchmarkTreeGenerator ./
goos: darwin
goarch: amd64
pkg: github.com/spikeekips/avl
BenchmarkTree10-8      	   14068	     83553 ns/op
BenchmarkTree100-8     	     919	   1284131 ns/op
BenchmarkTree200-8     	     423	   2768402 ns/op
BenchmarkTree300-8     	     271	   4393870 ns/op
BenchmarkTree400-8     	     195	   6032931 ns/op
BenchmarkTree500-8     	     152	   7734220 ns/op
BenchmarkTree600-8     	     126	   9456218 ns/op
BenchmarkTree700-8     	     105	  11239279 ns/op
BenchmarkTree800-8     	      80	  13124616 ns/op
BenchmarkTree900-8     	      72	  15055493 ns/op
BenchmarkTree1000-8    	      61	  16941144 ns/op
BenchmarkTree1100-8    	      56	  19040152 ns/op
BenchmarkTree1200-8    	      51	  20899512 ns/op
BenchmarkTree1300-8    	      45	  22839384 ns/op
BenchmarkTree1400-8    	      43	  24848736 ns/op
BenchmarkTree1500-8    	      40	  26989702 ns/op
BenchmarkTree1600-8    	      37	  28986443 ns/op
BenchmarkTree1700-8    	      34	  31502533 ns/op
BenchmarkTree1800-8    	      32	  33902428 ns/op
BenchmarkTree1900-8    	      30	  36840383 ns/op
BenchmarkTree10000-8   	       5	 235527047 ns/op
```

## TODO

* [ ] Code Documentation with examples
* [ ] Proof support for hashable node
