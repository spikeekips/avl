package cmd

import (
	"golang.org/x/xerrors"

	"github.com/spikeekips/avl"
)

type MutableNode struct {
	key    []byte
	height int16
	left   avl.MutableNode
	right  avl.MutableNode
	value  int
}

func NewMutableNode(key []byte) *MutableNode {
	return &MutableNode{key: key}
}

func (em *MutableNode) Key() []byte {
	return em.key
}

func (em *MutableNode) Height() int16 {
	return em.height
}

func (em *MutableNode) SetHeight(height int16) error {
	if height < 0 {
		return xerrors.Errorf("height must be greater than zero; height=%d", height)
	}

	em.height = height

	return nil
}

func (em *MutableNode) Left() avl.MutableNode {
	return em.left
}

func (em *MutableNode) LeftKey() []byte {
	if em.left == nil {
		return nil
	}

	return em.left.Key()
}

func (em *MutableNode) SetLeft(node avl.MutableNode) error {
	if node == nil {
		em.left = nil
		return nil
	}

	if avl.EqualKey(em.key, node.Key()) {
		return xerrors.Errorf("left is same node; key=%v", em.key)
	}

	em.left = node.(*MutableNode)

	return nil
}

func (em *MutableNode) Right() avl.MutableNode {
	return em.right
}

func (em *MutableNode) RightKey() []byte {
	if em.right == nil {
		return nil
	}

	return em.right.Key()
}

func (em *MutableNode) SetRight(node avl.MutableNode) error {
	if node == nil {
		em.right = nil
		return nil
	}

	if avl.EqualKey(em.key, node.Key()) {
		return xerrors.Errorf("right is same node; key=%v", em.key)
	}

	em.right = node.(*MutableNode)

	return nil
}

func (em *MutableNode) Merge(node avl.MutableNode) error {
	e, ok := node.(*MutableNode)
	if !ok {
		return xerrors.Errorf("merge node is not *MutableNode; node=%T", node)
	}

	em.value = e.value

	return nil
}
