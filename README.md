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


## TODO

* [ ] Performance test
* [ ] Code Documentation
* [ ] Examples
* [ ] Proof support for hashable node
