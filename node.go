package avl

import "bytes"

var (
	InvalidNodeError = NewWrapError("invalid node")
)

type Node interface {
	Key() []byte
	Height() int16
	LeftKey() []byte
	RightKey() []byte
}

type MutableNode interface {
	Node
	SetHeight(int16) error
	Left() MutableNode
	Right() MutableNode
	SetLeft(MutableNode) error
	SetRight(MutableNode) error
	Merge(source MutableNode) error
}

func IsEqualKey(a, b []byte) bool {
	return bytes.Equal(a, b)
}

func CompareKey(a, b []byte) int {
	return bytes.Compare(a, b)
}

func IsValidNode(node, left, right Node) error {
	// check empty key
	if node.Key() == nil || len(node.Key()) < 1 {
		return InvalidNodeError.Wrapf("key is empty")
	}

	// key of children correctness
	if left != nil && CompareKey(left.Key(), node.Key()) >= 0 {
		return InvalidNodeError.Wrapf(
			"left is greater: left=%v > node=%v",
			left.Key(), node.Key(),
		)
	}
	if right != nil && CompareKey(right.Key(), node.Key()) <= 0 {
		return InvalidNodeError.Wrapf(
			"right is lesser: right=%v > node=%v",
			right.Key(), node.Key(),
		)
	}

	// check height
	if left == nil && right == nil {
		if node.Height() != 0 {
			return InvalidNodeError.Wrapf("height must be 0 without children; height=%d", node.Height())
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
				"height must be +1 by children; left_or_right=%d height=%d",
				baseHeight, node.Height(),
			)
		}
	}

	return nil
}

func isSiblingNodesViolated(a, b Node) (bool /* left(true) or right(false) violated */, bool /* violated */) {
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
