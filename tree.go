package avl

import (
	"github.com/rs/zerolog"
)

var (
	NodeNotFoundInPoolError = NewWrapError("not not found in pool")
)

type NodeTraverseFunc func(Node) (bool, error)

type Tree struct {
	*Logger
	nodePool NodePool
	root     Node
}

func NewTree(rootKey []byte, nodePool NodePool) (*Tree, error) {
	if rootKey == nil {
		return nil, NodeNotFoundInPoolError.Wrapf("empty root")
	}

	root, err := nodePool.Get(rootKey)
	if err != nil {
		return nil, err
	} else if root == nil {
		return nil, NodeNotFoundInPoolError.Wrapf("root key=%x", rootKey)
	}

	return &Tree{
		Logger: NewLogger(func(c zerolog.Context) zerolog.Context {
			return c.Str("module", "avl_tree")
		}),
		nodePool: nodePool,
		root:     root,
	}, nil
}

func (tr *Tree) NodePool() NodePool {
	return tr.nodePool
}

func (tr *Tree) Root() Node {
	return tr.root
}

func (tr *Tree) getLeaf(node Node, isLeft bool) (Node, error) {
	var key []byte
	if isLeft {
		key = node.LeftKey()
	} else {
		key = node.RightKey()
	}
	if key == nil {
		return nil, nil
	}

	return tr.nodePool.Get(key)
}

func (tr *Tree) Get(key []byte) (Node, error) {
	log_ := tr.Log().With().Bytes("key", key).Logger()

	if tr.root == nil {
		log_.Debug().Msg("empty tree")
		return nil, nil
	}

	var parent Node
	parent = tr.root

	var depth int
	for {
		c := CompareKey(key, parent.Key())
		if c == 0 {
			log_.Debug().Int("depth", depth).Msg("found node by key")
			return parent, nil
		}

		var err error
		if parent, err = tr.getLeaf(parent, c < 0); err != nil {
			return nil, err
		}
		if parent == nil {
			break
		}
		depth++
	}

	return nil, nil
}

func (tr *Tree) GetWithParents(key []byte) (Node, []Node, error) {
	log_ := tr.Log().With().Bytes("key", key).Logger()

	if tr.root == nil {
		log_.Debug().Msg("empty tree")
		return nil, nil, nil
	}

	var parents []Node
	parent := tr.root

	var depth int
	for {
		c := CompareKey(key, parent.Key())

		if c == 0 {
			log_.Debug().Int("depth", depth).Msg("found node by key")
			return parent, parents, nil
		}
		parents = append(parents, parent)

		var err error
		if parent, err = tr.getLeaf(parent, c < 0); err != nil {
			return nil, nil, err
		}
		if parent == nil {
			break
		}
		depth++
	}

	return nil, nil, nil
}

func (tr *Tree) Traverse(f NodeTraverseFunc) error {
	if tr.root == nil {
		return nil
	}

	err, _ := tr.traverse(tr.root, f)
	return err
}

func (tr *Tree) traverse(node Node, f NodeTraverseFunc) (error, bool) {
	if node == nil {
		return nil, true
	}

	if keep, err := f(node); !keep || err != nil {
		return err, false
	}

	var left, right Node
	var err error
	if left, err = tr.getLeaf(node, true); err != nil {
		return err, false
	}
	if right, err = tr.getLeaf(node, false); err != nil {
		return err, false
	}

	if left != nil {
		if err, keep := tr.traverse(left, f); err != nil || !keep {
			return err, keep
		}
	}
	if right != nil {
		if err, keep := tr.traverse(right, f); err != nil || !keep {
			return err, keep
		}
	}

	return nil, true
}

func (tr *Tree) IsValid() error {
	return NewTreeValidator(tr).IsValid()
}
