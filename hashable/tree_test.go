package hashable

import (
	"fmt"
	"testing"

	"github.com/spikeekips/avl"
	"github.com/stretchr/testify/suite"
)

type testHashableTree struct {
	suite.Suite
}

func (t *testHashableTree) makeNodeKey(i int) []byte {
	return []byte(fmt.Sprintf("%03d", i))
}

func (t *testHashableTree) newNode(i int) HashableNode {
	return NewBaseHashableNode(
		avl.NewBaseNode(t.makeNodeKey(i)),
	)
}

func (t *testHashableTree) newTree() *Tree {
	tr := NewTree(avl.NewMapNodePool())
	_ = tr.SetLogger(log)

	return tr
}

func (t *testHashableTree) TestNew() {
	tr := t.newTree()

	_, err := tr.Add(t.newNode(10))
	t.NoError(err)

	var root avl.Node
	var rootHash []byte
	{
		root = tr.Root()
		hn, ok := root.(*BaseHashableNode)
		t.True(ok)

		rootHash = hn.Hash()
		t.NotEmpty(rootHash)
	}

	{
		_, _ = tr.Add(t.newNode(5))
		_, _ = tr.Add(t.newNode(15))

		t.Equal(root.Key(), tr.Root().Key())
		t.NotEqual(rootHash, tr.Root().(*BaseHashableNode).Hash())
	}
}

func (t *testHashableTree) TestResetHashSimple() {
	tr := t.newTree()

	_, _ = tr.Add(t.newNode(100))
	_, _ = tr.Add(t.newNode(50))
	_, _ = tr.Add(t.newNode(150))
	_, _ = tr.Add(t.newNode(30))
	_, _ = tr.Add(t.newNode(10))

	_ = tr.Traverse(func(node avl.Node) (bool, error) {
		hn := node.(*BaseHashableNode)
		t.Nil(hn.RawHash())
		return true, nil
	})

	rootHash, err := tr.NodeHash(tr.Root().Key())
	t.NoError(err)
	t.NotNil(rootHash)

	_ = tr.Traverse(func(node avl.Node) (bool, error) {
		hn := node.(*BaseHashableNode)
		t.NotNil(hn.RawHash())
		return true, nil
	})

	_, _ = tr.Add(t.newNode(5))

	getNode := func(i int) *BaseHashableNode {
		node, err := tr.Get(t.makeNodeKey(i))
		t.NoError(err)
		return node.(*BaseHashableNode)
	}

	for _, p := range []int{100, 30, 10} {
		t.Nil(getNode(p).RawHash(), "hash of node, %d  should be nil, but not", p)
	}
	for _, p := range []int{50, 150} {
		t.NotNil(getNode(p).RawHash(), "hash of node, %d  should be not nil, but not", p)
	}
}

func (t *testHashableTree) TestResetHashLeftLeftRotation() {
	tr := t.newTree()

	_, _ = tr.Add(t.newNode(100))
	_, _ = tr.Add(t.newNode(50))
	_, _ = tr.Add(t.newNode(150))
	_, _ = tr.Add(t.newNode(30))
	_, _ = tr.Add(t.newNode(10))

	_ = tr.Traverse(func(node avl.Node) (bool, error) {
		hn := node.(*BaseHashableNode)
		t.Nil(hn.RawHash())
		return true, nil
	})

	rootHash, err := tr.NodeHash(tr.Root().Key())
	t.NoError(err)
	t.NotNil(rootHash)

	_ = tr.Traverse(func(node avl.Node) (bool, error) {
		hn := node.(*BaseHashableNode)
		t.NotNil(hn.RawHash())
		return true, nil
	})

	_, _ = tr.Add(t.newNode(5))

	getNode := func(i int) *BaseHashableNode {
		node, err := tr.Get(t.makeNodeKey(i))
		t.NoError(err)
		return node.(*BaseHashableNode)
	}

	for _, p := range []int{100, 30, 10} {
		t.Nil(getNode(p).RawHash(), "hash of node, %d  should be nil, but not", p)
	}
	for _, p := range []int{50, 150} {
		t.NotNil(getNode(p).RawHash(), "hash of node, %d  should be not nil, but not", p)
	}

	rootHash, err = tr.RootHash()
	t.NoError(err)
	t.NotNil(rootHash)

	_ = tr.Traverse(func(node avl.Node) (bool, error) {
		hn := node.(*BaseHashableNode)
		t.NotNil(hn.RawHash())
		return true, nil
	})
}

func TestHashableTree(t *testing.T) {
	suite.Run(t, new(testHashableTree))
}
