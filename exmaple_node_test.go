package avl

import "golang.org/x/xerrors"

type ExampleMutableNode struct {
	key    []byte
	height int16
	left   MutableNode
	right  MutableNode
	value  int
}

func (em *ExampleMutableNode) Key() []byte {
	return em.key
}

func (em *ExampleMutableNode) Height() int16 {
	return em.height
}

func (em *ExampleMutableNode) SetHeight(height int16) error {
	if height < 0 {
		return xerrors.Errorf("height must be greater than zero; height=%d", height)
	}

	em.height = height

	return nil
}

func (em *ExampleMutableNode) Left() MutableNode {
	return em.left
}

func (em *ExampleMutableNode) LeftKey() []byte {
	if em.left == nil {
		return nil
	}

	return em.left.Key()
}

func (em *ExampleMutableNode) SetLeft(node MutableNode) error {
	if node == nil {
		em.left = nil
		return nil
	}

	if EqualKey(em.key, node.Key()) {
		return xerrors.Errorf("left is same node; key=%v", em.key)
	}

	em.left = node

	return nil
}

func (em *ExampleMutableNode) Right() MutableNode {
	return em.right
}

func (em *ExampleMutableNode) RightKey() []byte {
	if em.right == nil {
		return nil
	}

	return em.right.Key()
}

func (em *ExampleMutableNode) SetRight(node MutableNode) error {
	if node == nil {
		em.right = nil
		return nil
	}

	if EqualKey(em.key, node.Key()) {
		return xerrors.Errorf("right is same node; key=%v", em.key)
	}

	em.right = node

	return nil
}

func (em *ExampleMutableNode) Merge(node MutableNode) error {
	e, ok := node.(*ExampleMutableNode)
	if !ok {
		return xerrors.Errorf("merge node is not *ExampleMutableNode; node=%T", node)
	}

	em.value = e.value

	return nil
}

type ExampleNode struct {
	key    []byte
	height int16
	left   []byte
	right  []byte
}

func (em *ExampleNode) Key() []byte {
	return em.key
}

func (em *ExampleNode) Height() int16 {
	return em.height
}

func (em *ExampleNode) LeftKey() []byte {
	if em.left == nil || len(em.left) < 1 {
		return nil
	}

	return em.left
}

func (em *ExampleNode) RightKey() []byte {
	if em.right == nil || len(em.right) < 1 {
		return nil
	}

	return em.right
}
