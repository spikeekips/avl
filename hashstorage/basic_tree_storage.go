package hashstorage

import (
	"github.com/rs/zerolog"
	"github.com/spikeekips/avl"
	"github.com/spikeekips/avl/hashable"
)

var (
	TreeAlreadyExistsError = avl.NewWrapError("tree already exists in storage")
	TreeNotExistsError     = avl.NewWrapError("tree not exist in storage")
	InvalidTreeError       = avl.NewWrapError("invalid tree storage")
)

type BasicTreeStorage struct {
	*avl.Logger
	name    []byte
	tr      *Tree
	storage Storage
}

func NewBasicTreeStorage(name []byte, tr *Tree, storage Storage) (*BasicTreeStorage, error) {
	if e, err := storage.GetRoot(name); err != nil {
		return nil, err
	} else if e != nil {
		return nil, TreeAlreadyExistsError.Wrapf("name=%v", name)
	}

	return &BasicTreeStorage{
		Logger: avl.NewLogger(func(c zerolog.Context) zerolog.Context {
			return c.Str("module", "avl_tree_storage")
		}),
		name:    name,
		tr:      tr,
		storage: storage,
	}, nil
}

func LoadTreeFromBasicTreeStorage(name []byte, storage Storage, nodePool avl.NodePool) (*Tree, error) {
	rootKey, err := storage.GetRoot(name)
	if err != nil {
		return nil, err
	} else if rootKey == nil {
		return nil, TreeNotExistsError.Wrapf("name=%v", name)
	}

	root, err := nodePool.Get(rootKey)
	if err != nil {
		return nil, InvalidTreeError.Wrap(err)
	}

	tr := NewTree(nodePool)
	_ = tr.SetRoot(root)

	return tr, nil
}

func (ts *BasicTreeStorage) Save() error {
	if _, err := ts.tr.RootHash(); err != nil {
		return err
	}

	if err := ts.tr.IsValid(); err != nil {
		return err
	} else if ts.tr.Root() == nil {
		return nil
	}

	batch := ts.storage.Batch()

	if err := batch.SetRoot(ts.name, ts.tr.Root().Key()); err != nil {
		return err
	}

	err := ts.tr.Traverse(func(node avl.Node) (bool, error) {
		hn, ok := node.(hashable.HashableNode)
		if !ok {
			return false, hashable.NotHashableNodeError.Wrapf("key=%v", node.Key())
		}

		b, err := hn.MarshalBinary()
		if err != nil {
			return false, err
		}

		if err := batch.SetNode(hn.Key(), b); err != nil {
			return false, err
		}

		return true, nil
	})
	if err != nil {
		return err
	}

	if err := batch.Commit(); err != nil {
		return err
	}

	return nil
}
