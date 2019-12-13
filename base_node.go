package avl

import (
	"fmt"
)

type BaseNode struct {
	key    []byte
	left   []byte
	right  []byte
	height int16
}

func NewBaseNode(key []byte) *BaseNode {
	return &BaseNode{
		key:    key,
		height: 0,
	}
}

func (bd *BaseNode) String() string {
	var left, right []byte
	if bd.left != nil {
		left = bd.left
	}
	if bd.right != nil {
		right = bd.right
	}
	return fmt.Sprintf(
		"<BaseNode: key=%s height=%d left=%v right=%v>",
		bd.key,
		bd.height,
		left,
		right,
	)
}

func (bd BaseNode) Key() []byte {
	return bd.key
}

func (bd BaseNode) Height() int16 {
	return bd.height
}

func (bd *BaseNode) SetHeight(h int16) error {
	if h < 0 {
		return FailedToUpdateNodeError.Wrapf("height must be greater than 0 or equal; %v", h)
	} else if bd.height == h {
		return nil
	}

	bd.height = h

	return nil
}

func (bd BaseNode) LeftKey() []byte {
	if bd.left == nil {
		return nil
	}

	return bd.left
}

func (bd BaseNode) RightKey() []byte {
	if bd.right == nil {
		return nil
	}

	return bd.right
}

func (bd *BaseNode) SetLeftKey(key []byte) error {
	bd.left = key
	return nil
}

func (bd *BaseNode) SetRightKey(key []byte) error {
	bd.right = key
	return nil
}
