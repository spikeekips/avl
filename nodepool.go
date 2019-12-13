package avl

import (
	"sync"

	"golang.org/x/xerrors"
)

type NodePool interface {
	Get(key []byte) (Node, error)
	Set(Node) error
	Traverse(NodeTraverseFunc) error
}

type SyncMapNodePool struct {
	m *sync.Map
}

func NewSyncMapNodePool(m *sync.Map) *SyncMapNodePool {
	if m == nil {
		m = &sync.Map{}
	}

	return &SyncMapNodePool{m: m}
}

func (mn *SyncMapNodePool) Get(key []byte) (Node, error) {
	v, found := mn.m.Load(string(key))
	if !found {
		return nil, nil
	} else if node, ok := v.(Node); !ok {
		return nil, InvalidNodeError.Wrapf("not Node type: %T", v)
	} else {
		return node, nil
	}
}

func (mn *SyncMapNodePool) Set(node Node) error {
	mn.m.Store(string(node.Key()), node)
	return nil
}

func (mn *SyncMapNodePool) Traverse(f NodeTraverseFunc) error {
	var err error
	mn.m.Range(func(_, value interface{}) bool {
		var keep bool
		node, ok := value.(Node)
		if !ok {
			err = xerrors.Errorf("invalid type found in nodepool; value=%T", value)
			return false
		}
		if keep, err = f(node); err != nil {
			return false
		}

		return keep
	})

	return err
}

type MapNodePool struct {
	m map[string]Node
}

func NewMapNodePool(m map[string]Node) *MapNodePool {
	if m == nil {
		m = map[string]Node{}
	}

	return &MapNodePool{m: m}
}

func (mn *MapNodePool) Get(key []byte) (Node, error) {
	node, found := mn.m[string(key)]
	if !found {
		return nil, nil
	}

	return node, nil
}

func (mn *MapNodePool) Set(node Node) error {
	mn.m[string(node.Key())] = node

	return nil
}

func (mn *MapNodePool) Traverse(f NodeTraverseFunc) error {
	for _, node := range mn.m {
		if keep, err := f(node); err != nil {
			return err
		} else if !keep {
			break
		}
	}

	return nil
}

type MapMutableNodePool struct {
	m map[string]MutableNode
}

func NewMapMutableNodePool(m map[string]MutableNode) *MapMutableNodePool {
	if m == nil {
		m = map[string]MutableNode{}
	}

	return &MapMutableNodePool{m: m}
}

func (mn *MapMutableNodePool) Get(key []byte) (Node, error) {
	node, found := mn.m[string(key)]
	if !found {
		return nil, nil
	}

	return node, nil
}

func (mn *MapMutableNodePool) Set(node Node) error {
	n, ok := node.(MutableNode)
	if !ok {
		return xerrors.Errorf("not MutableNode; %T", node)
	}

	mn.m[string(node.Key())] = n

	return nil
}

func (mn *MapMutableNodePool) Traverse(f NodeTraverseFunc) error {
	for _, node := range mn.m {
		if keep, err := f(node); err != nil {
			return err
		} else if !keep {
			break
		}
	}

	return nil
}
