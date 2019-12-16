package avl

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.org/x/xerrors"
)

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

type testNodeFuncs struct {
	suite.Suite
}

func (t *testNodeFuncs) TestCompareKey() {
	type b []byte
	cases := []struct {
		name   string
		a      []byte
		b      []byte
		result int
	}{
		{"1 char#0", b("a"), b("b"), -1},
		{"1 char#1", b("a"), b("a"), 0},
		{"1 char with space", b("a"), b("a "), -1},
		{"int#0", b("0"), b("1"), -1},
		{"int#1", b("1"), b("0"), 1},
		{"int#2", b("0"), b("01"), -1},
	}

	for _, c := range cases {
		c := c
		t.Run(
			c.name,
			func() {
				t.Equal(CompareKey(c.a, c.b), c.result, "%s:  result not match", c.name)
			},
		)
	}
}

func (t *testNodeFuncs) TestEqualKey() {
	type b []byte
	cases := []struct {
		name   string
		a      []byte
		b      []byte
		result bool
	}{
		{"1 char#0", b("a"), b("b"), false},
		{"1 char#1", b("a"), b("a"), true},
		{"1 char with space", b("a"), b("a "), false},
		{"int#0", b("0"), b("1"), false},
		{"int#1", b("1"), b("0"), false},
		{"int#2", b("0"), b("01"), false},
	}

	for _, c := range cases {
		c := c
		t.Run(
			c.name,
			func() {
				t.Equal(EqualKey(c.a, c.b), c.result, "%s:  result not match", c.name)
			},
		)
	}
}

func (t *testNodeFuncs) TestIsValidNode() {
	type b []byte
	cases := []struct {
		name   string
		node   Node
		left   Node
		right  Node
		result string
	}{
		{
			name: "just with key",
			node: &ExampleMutableNode{key: b("showme")},
		},
		{
			name:   "1 height without leafs",
			node:   &ExampleMutableNode{key: b("showme"), height: 1},
			result: "height must be 0",
		},
		{
			name:   "0 height with left",
			node:   &ExampleMutableNode{key: b("5"), height: 0},
			left:   &ExampleMutableNode{key: b("3"), height: 0},
			result: "height must be +1",
		},
		{
			name:   "0 height with right",
			node:   &ExampleMutableNode{key: b("5"), height: 0},
			right:  &ExampleMutableNode{key: b("8"), height: 0},
			result: "height must be +1",
		},
		{
			name:   "bad height with left",
			node:   &ExampleMutableNode{key: b("5"), height: 3},
			left:   &ExampleMutableNode{key: b("3"), height: 0},
			result: "height must be +1",
		},
		{
			name:   "bad height with right",
			node:   &ExampleMutableNode{key: b("5"), height: 3},
			right:  &ExampleMutableNode{key: b("8"), height: 0},
			result: "height must be +1",
		},
		{
			name:   "bad left",
			node:   &ExampleMutableNode{key: b("5"), height: 1},
			left:   &ExampleMutableNode{key: b("8"), height: 0},
			result: "left must be lesser",
		},
		{
			name:   "bad right",
			node:   &ExampleMutableNode{key: b("5"), height: 1},
			left:   &ExampleMutableNode{key: b("3"), height: 0},
			result: "right must be greater",
		},
		{
			name:  "height with left is greater",
			node:  &ExampleMutableNode{key: b("5"), height: 2},
			left:  &ExampleMutableNode{key: b("3"), height: 1},
			right: &ExampleMutableNode{key: b("8"), height: 0},
		},
		{
			name:   "bad height with left is greater",
			node:   &ExampleMutableNode{key: b("5"), height: 1},
			left:   &ExampleMutableNode{key: b("3"), height: 1},
			right:  &ExampleMutableNode{key: b("8"), height: 0},
			result: "height must be +1",
		},
		{
			name:   "height violation; left > right+1",
			node:   &ExampleMutableNode{key: b("5"), height: 3},
			left:   &ExampleMutableNode{key: b("3"), height: 2},
			right:  &ExampleMutableNode{key: b("8"), height: 0},
			result: "leaf is violated",
		},
		{
			name:   "height violation; right > left+1",
			node:   &ExampleMutableNode{key: b("5"), height: 3},
			left:   &ExampleMutableNode{key: b("3"), height: 0},
			right:  &ExampleMutableNode{key: b("8"), height: 2},
			result: "leaf is violated",
		},
	}

	for _, c := range cases {
		c := c
		t.Run(
			c.name,
			func() {
				node := c.node.(*ExampleMutableNode)
				if c.left != nil {
					node.left = c.left.(*ExampleMutableNode)
				}
				if c.right != nil {
					node.right = c.right.(*ExampleMutableNode)
				}

				err := IsValidNode(node, c.left, c.right)
				if len(c.result) < 1 {
					t.Nil(err, "%s: result must be nil, but %v", err)
				} else {
					if err == nil {
						t.Nil(err, "%s: err must not be nil; expected=%v", c.result)
					} else {
						t.Contains(err.Error(), c.result, "%s: err does not match", c.name)
					}
				}
			},
		)
	}
}

func TestNodeFuncs(t *testing.T) {
	suite.Run(t, new(testNodeFuncs))
}
