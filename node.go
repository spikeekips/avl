package avl

import (
	"bytes"
)

var (
	FailedToUpdateNodeError = NewWrapError("failed to update node")
	InvalidNodeError        = NewWrapError("invalid node found")
)

type Node interface {
	Key() []byte
	Height() int16
	SetHeight(int16) error
	LeftKey() []byte
	SetLeftKey(key []byte) error
	RightKey() []byte
	SetRightKey([]byte) error
}

func IsValidNode(node, left, right Node) error {
	// check empty key
	if node.Key() == nil || len(node.Key()) < 1 {
		return InvalidNodeError.Wrapf("key is empty")
	}

	// key of children correctness
	if left != nil && CompareNode(left, node) >= 0 {
		return InvalidNodeError.Wrapf(
			"left is greater: left=%v > node=%v",
			left.Key(), node.Key(),
		)
	}
	if right != nil && CompareNode(right, node) <= 0 {
		return InvalidNodeError.Wrapf(
			"right is lesser: right=%v > node=%v",
			right.Key(), node.Key(),
		)
	}

	// check height
	if left == nil && right == nil && node.Height() != 0 {
		return InvalidNodeError.Wrapf("height must be 0 without children; height=%d", node.Height())
	}

	var baseHeight int16 = -1
	if left != nil {
		baseHeight = left.Height()
	}

	if right != nil && right.Height() > baseHeight {
		baseHeight = right.Height()
	}

	if node.Height() != baseHeight+1 {
		return InvalidNodeError.Wrapf(
			"height must be +1 by children; left_or_right=%d height=%d",
			baseHeight, node.Height(),
		)
	}

	return nil
}

func CompareNodeKey(a, b []byte) int {
	return bytes.Compare(a, b)
}

func CompareNode(a, b Node) int {
	return CompareNodeKey(a.Key(), b.Key())
}
