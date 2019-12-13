package hashable

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/spikeekips/avl"
)

type Tree struct {
	*avl.Logger
	*avl.Tree
}

func NewTree(nodePool avl.NodePool) *Tree {
	return &Tree{
		Logger: avl.NewLogger(func(c zerolog.Context) zerolog.Context {
			return c.Str("module", "avl_hashable_tree")
		}),
		Tree: avl.NewTree(nodePool),
	}
}

func (tr *Tree) SetLogger(l zerolog.Logger) *avl.Logger {
	_ = tr.Tree.SetLogger(l)
	return tr.Logger.SetLogger(l)
}

func (tr *Tree) Add(node HashableNode) ([]avl.Node /* updated Node */, error) {
	log_ := tr.Log().With().Bytes("key", node.Key()).Logger()

	parents, err := tr.Tree.Add(node)
	if err != nil {
		return nil, err
	}

	if tr.Log().GetLevel() == zerolog.DebugLevel {
		if len(parents) > 0 {
			le := log_.Debug()
			for i, p := range parents {
				le.Bytes(fmt.Sprintf("parent_%d", i), p.Key())
			}
			le.Msg("reset nodes hash")
		}
	}

	for _, p := range parents {
		p.(HashableNode).ResetHash()
	}

	return parents, nil
}

func (tr *Tree) nodeLeafHash(key []byte) ([]byte, error) {
	if key == nil {
		return nil, nil
	}

	node, err := tr.NodePool().Get(key)
	if err != nil {
		return nil, err
	} else if node == nil {
		return nil, avl.NodeNotFound.Wrapf("leaf key=%v", key)
	}

	hn := node.(HashableNode)
	if hn.RawHash() != nil {
		return hn.RawHash(), nil
	}

	return tr.NodeHash(node.Key())
}

func (tr *Tree) NodeHash(key []byte) ([]byte, error) {
	n, err := tr.NodePool().Get(key)
	if err != nil {
		return nil, err
	}

	node := n.(HashableNode)
	if node.RawHash() != nil {
		return node.RawHash(), nil
	}

	if h, err := tr.nodeLeafHash(node.LeftKey()); err != nil {
		return nil, err
	} else if h != nil {
		node.SetLeftHash(h)
	}

	if h, err := tr.nodeLeafHash(node.RightKey()); err != nil {
		return nil, err
	} else if h != nil {
		node.SetRightHash(h)
	}

	return node.Hash(), nil
}

func (tr *Tree) RootHash() ([]byte, error) {
	if tr.Root() == nil {
		return nil, nil
	}

	return tr.NodeHash(tr.Root().Key())
}
