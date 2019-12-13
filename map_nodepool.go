package avl

import "sync"

type SyncMapNodePool struct {
	m sync.Map
}

func NewSyncMapNodePool() *SyncMapNodePool {
	return &SyncMapNodePool{m: sync.Map{}}
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

type MapNodePool struct {
	m map[string]Node
}

func NewMapNodePool() *MapNodePool {
	return &MapNodePool{m: map[string]Node{}}
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
