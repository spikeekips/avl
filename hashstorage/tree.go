package hashstorage

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/spikeekips/avl"
	"github.com/spikeekips/avl/hashable"
)

type Tree struct {
	*hashable.Tree
}

func NewTree(nodePool avl.NodePool) *Tree {
	return &Tree{Tree: hashable.NewTree(nodePool)}
}

func (tr *Tree) Adds(nodes []hashable.HashableNode) ([]avl.Node, error) {
	var log_ zerolog.Logger
	if tr.Log().GetLevel() == zerolog.DebugLevel {
		var keys [][]byte
		for _, n := range nodes {
			keys = append(keys, n.Key())
		}

		log_ = tr.Log().With().Interface("keys", keys).Logger()
	}

	parents := map[string]avl.Node{}
	for _, n := range nodes {
		ps, err := tr.Tree.Tree.Add(n)
		if err != nil {
			return nil, err
		}
		for _, p := range ps {
			if _, found := parents[string(p.Key())]; found {
				continue
			}
			parents[string(p.Key())] = p
		}
	}

	if tr.Log().GetLevel() == zerolog.DebugLevel {
		if len(parents) > 0 {
			le := log_.Debug()
			for i, p := range parents {
				le.Bytes(fmt.Sprintf("parent_%s", i), p.Key())
			}
			le.Msg("reset nodes hash")
		}
	}

	var collectedParents []avl.Node
	for _, p := range parents {
		p.(hashable.HashableNode).ResetHash()
		collectedParents = append(collectedParents, p)
	}

	return collectedParents, nil
}
