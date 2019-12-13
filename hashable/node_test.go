package hashable

import (
	"testing"

	"github.com/spikeekips/avl"
	"github.com/stretchr/testify/suite"
)

type testBaseHashableNode struct {
	suite.Suite
}

func (t *testBaseHashableNode) TestNew() {
	key := []byte("showme")
	n := NewBaseHashableNode(avl.NewBaseNode(key))

	t.Equal(key, n.Key())

	var hs []byte
	{ // make hash with key
		hs = n.Hash()
		t.NotNil(hs)
	}

	{ // update node, set height
		_ = n.SetHeight(n.Height() + 1)

		h := n.Hash()
		t.NotNil(h)
		t.NotEqual(hs, h)
	}
}

func (t *testBaseHashableNode) TestNewWithHashable() {
	key := []byte("showme")

	n0 := NewBaseHashableNode(avl.NewBaseNode(key))
	n2 := NewBaseHashableNode(avl.NewBaseNode([]byte("showme00")))

	t.NotEqual(n0, n2)
}

func TestBaseHashableNode(t *testing.T) {
	suite.Run(t, new(testBaseHashableNode))
}
