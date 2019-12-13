package hashstorage

import (
	"bytes"
	"encoding/gob"

	"github.com/rs/zerolog"
	"github.com/spikeekips/avl"
	"github.com/spikeekips/avl/hashable"
)

type DumpTreeStorage struct {
	*avl.Logger
	name    []byte
	tr      *Tree
	storage Storage
}

func NewDumpTreeStorage(name []byte, tr *Tree, storage Storage) (*DumpTreeStorage, error) {
	if e, err := storage.GetRoot(name); err != nil {
		return nil, err
	} else if e != nil {
		return nil, TreeAlreadyExistsError.Wrapf("name=%v", name)
	}

	return &DumpTreeStorage{
		Logger: avl.NewLogger(func(c zerolog.Context) zerolog.Context {
			return c.Str("module", "avl_tree_storage")
		}),
		name:    name,
		tr:      tr,
		storage: storage,
	}, nil
}

func LoadTreeFromDumpTreeStorage(name []byte, storage Storage, nodePool *StorageNodePool) (*Tree, error) {
	b, err := storage.GetRaw(name)
	if err != nil {
		return nil, err
	} else if b == nil {
		return nil, TreeNotExistsError.Wrapf("name=%v", name)
	}

	var dump dumpTreeMarshal
	if err := gob.NewDecoder(bytes.NewBuffer(b)).Decode(&dump); err != nil {
		return nil, err
	} else if !bytes.Equal(dump.Name, name) {
		return nil, InvalidTreeError.Wrapf("name does not match name=%v tree=%v", name, dump.Name)
	}

	for _, b := range dump.Nodes {
		node, err := nodePool.UnmarshalNode(b)
		if err != nil {
			return nil, err
		} else if err := nodePool.Set(node); err != nil {
			return nil, err
		}
	}

	root, err := nodePool.Get(dump.RootKey)
	if err != nil {
		return nil, err
	}

	tr := NewTree(nodePool)
	_ = tr.SetRoot(root)

	return tr, nil
}

func (ts *DumpTreeStorage) Save() error {
	if _, err := ts.tr.RootHash(); err != nil {
		return err
	}

	if err := ts.tr.IsValid(); err != nil {
		return err
	} else if ts.tr.Root() == nil {
		return nil
	}

	dump := dumpTreeMarshal{
		Name:    ts.name,
		RootKey: ts.tr.Root().Key(),
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

		dump.Nodes = append(dump.Nodes, b)

		return true, nil
	})
	if err != nil {
		return err
	}

	var b bytes.Buffer
	if err := gob.NewEncoder(&b).Encode(dump); err != nil {
		return err
	}

	if err := ts.storage.SetRaw(ts.name, b.Bytes()); err != nil {
		return err
	}

	return nil
}

type dumpTreeMarshal struct {
	Name    []byte
	RootKey []byte
	Nodes   [][]byte
}
