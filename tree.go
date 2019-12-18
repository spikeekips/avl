package avl

import (
	"github.com/rs/zerolog"
)

var (
	NodeNotFoundInPoolError = NewWrapError("node not found in pool")
)

// NodeTraverseFunc is used for Tree.Traverse(). If keep is false, traversing
// will be stopped and error also stops traversing.
type NodeTraverseFunc func(Node) (keep bool, err error)

// Tree is AVL tree. Mainly Tree is used for loading the existing nodes.
type Tree struct {
	*Logger
	nodePool NodePool
	root     Node
}

// NewTree loads tree from NodePool.
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

// NodePool returns NodePool of this tree.
func (tr *Tree) NodePool() NodePool {
	return tr.nodePool
}

// Root returns root node.
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

// Get finds and returns node by key. It traverse the entire tree.  Unlike
// NodePool.Get() the only organized(not orphan) node will be returned.  For
// performance, NodePool.Get() will be better.
func (tr *Tree) Get(key []byte) (Node, error) {
	logs := tr.Log().With().Bytes("key", key).Logger()

	if tr.root == nil {
		logs.Debug().Msg("empty tree")
		return nil, nil
	}

	var parent Node
	parent = tr.root

	var depth int
	for {
		c := CompareKey(key, parent.Key())
		if c == 0 {
			logs.Debug().Int("depth", depth).Msg("found node by key")
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

// GetWithParents returns node with it's parents node.
func (tr *Tree) GetWithParents(key []byte) (Node, []Node, error) {
	logs := tr.Log().With().Bytes("key", key).Logger()

	if tr.root == nil {
		logs.Debug().Msg("empty tree")
		return nil, nil, nil
	}

	var parents []Node
	parent := tr.root

	var depth int
	for {
		c := CompareKey(key, parent.Key())

		if c == 0 {
			logs.Debug().Int("depth", depth).Msg("found node by key")
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

// Traverse traverses the entire tree. The error of NodeTraverseFunc mainly
// error is from the external storage or other system.
func (tr *Tree) Traverse(f NodeTraverseFunc) error {
	if tr.root == nil {
		return nil
	}

	_, err := tr.traverse(tr.root, f)
	return err
}

func (tr *Tree) traverse(node Node, f NodeTraverseFunc) (bool, error) {
	if node == nil {
		return true, nil
	}

	if keep, err := f(node); err != nil {
		return false, err
	} else if !keep {
		return true, nil
	}

	var left, right Node
	var err error
	if left, err = tr.getLeaf(node, true); err != nil {
		return false, err
	}
	if right, err = tr.getLeaf(node, false); err != nil {
		return false, err
	}

	if left != nil {
		if keep, err := tr.traverse(left, f); err != nil {
			return false, err
		} else if !keep {
			return true, nil
		}
	}
	if right != nil {
		if keep, err := tr.traverse(right, f); err != nil {
			return false, err
		} else if !keep {
			return true, nil
		}
	}

	return true, nil
}

// IsValid checks whether Tree is valid or not.
func (tr *Tree) IsValid() error {
	return NewTreeValidator(tr).IsValid()
}
