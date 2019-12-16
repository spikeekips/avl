package avl

import (
	"bytes"
)

var (
	InvalidNodeError = NewWrapError("invalid node")
)

// Node defines the basic node. By comparing with MutableNode, Node stands for
// immutable node.
type Node interface {
	Key() []byte
	Height() int16
	LeftKey() []byte
	RightKey() []byte
}

// MutableNode is, the name said, is mutable node.Mainly it is used for TreeGenerator.
type MutableNode interface {
	Node
	// SetHeight set height.
	SetHeight(int16) error
	// Left() returns the MutableNode of left leaf.
	Left() MutableNode
	// Right() returns the MutableNode of right leaf.
	Right() MutableNode
	// SetLeft() replaces left leaf.
	SetLeft(MutableNode) error
	// SetRight() replaces right leaf.
	SetRight(MutableNode) error
	// Merge() merges source node. The basic properties(key, height, left right
	// key) will not be merged.
	Merge(source MutableNode) error
}

// EqualKey checks node keys are same. it acts like bytes.Equal()
func EqualKey(a, b []byte) bool {
	return bytes.Equal(a, b)
}

// CompareKey compares node keys. it acts like bytes.Compare()
func CompareKey(a, b []byte) int {
	return bytes.Compare(a, b)
}

// IsValidNode checks node is valid and well defined.
func IsValidNode(node, left, right Node) error {
	// check empty key
	if node.Key() == nil || len(node.Key()) < 1 {
		return InvalidNodeError.Wrapf("key is empty")
	}

	// key of leaf correctness
	if left != nil && CompareKey(left.Key(), node.Key()) >= 0 {
		return InvalidNodeError.Wrapf(
			"left must be lesser: left=%v > node=%v",
			left.Key(), node.Key(),
		)
	}
	if right != nil && CompareKey(right.Key(), node.Key()) <= 0 {
		return InvalidNodeError.Wrapf(
			"right must be greater: right=%v > node=%v",
			right.Key(), node.Key(),
		)
	}

	// check height
	if left == nil && right == nil {
		if node.Height() != 0 {
			return InvalidNodeError.Wrapf("height must be 0 without leaf; height=%d", node.Height())
		}
	} else if isLeft, violated := isSiblingNodesViolated(left, right); violated {
		return InvalidNodeError.Wrapf("left or right leaf is violated; isLeft=%v", isLeft)
	} else {
		var baseHeight int16 = -1
		if left != nil {
			baseHeight = left.Height()
		}

		if right != nil && right.Height() > baseHeight {
			baseHeight = right.Height()
		}

		if node.Height() != baseHeight+1 {
			return InvalidNodeError.Wrapf(
				"height must be +1 by leaf; left_or_right=%d height=%d",
				baseHeight, node.Height(),
			)
		}
	}

	return nil
}

// isSiblingNodesViolated checks the AVL violation of node. It checks the height
// of leaves and it's own height.
func isSiblingNodesViolated(a, b Node) (
	bool, /* left(true), right(false) */
	bool, /* violated */
) {
	if a == nil && b == nil {
		return false, false
	} else if a == nil || b == nil {
		if a == nil {
			if b.Height() > 1 {
				return false, true
			}
		} else {
			if a.Height() > 1 {
				return true, true
			}
		}

		return false, false
	}

	d := a.Height() - b.Height()
	if d < 2 && d > -2 {
		return false, false
	}

	return d > 1, true
}
