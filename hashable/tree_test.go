package hashable

import (
	"fmt"
	"testing"

	"github.com/spikeekips/avl"
	"github.com/stretchr/testify/suite"
)

type testTree struct {
	suite.Suite
}

func (t *testTree) newNode(i int) *ExampleHashableMutableNode {
	return &ExampleHashableMutableNode{
		key: []byte(fmt.Sprintf("%03d", i)),
	}
}

func (t *testTree) TestHashNodes() {
	tg := avl.NewTreeGenerator()

	nodes := map[int]HashableMutableNode{}

	for i := 0; i < 10; i++ {
		n := t.newNode(i)
		_, err := tg.Add(n)
		t.NoError(err)

		nodes[i] = n
	}

	tr, err := tg.Tree()
	t.NoError(err)

	tr.Traverse(func(node avl.Node) (bool, error) {
		t.Nil(node.(*ExampleHashableMutableNode).hash)
		return true, nil
	})

	rootHash := tg.Root().(HashableNode).Hash()
	t.NotNil(rootHash)

	t.Equal(nodes[3].Hash(), rootHash)

	// hash was generated
	tr.Traverse(func(node avl.Node) (bool, error) {
		t.NotNil(node.(*ExampleHashableMutableNode).hash)
		return true, nil
	})
}

func (t *testTree) TestHashNodeProof() {
	tg := avl.NewTreeGenerator()

	nodes := map[int]HashableMutableNode{}

	for i := 0; i < 21; i++ {
		n := t.newNode(i)
		_, err := tg.Add(n)
		t.NoError(err)

		nodes[i] = n
	}

	tr, _ := tg.Tree()

	{
		// find parents of 0 height node, 14
		zeroKey := nodes[14].Key()
		var zero avl.Node
		var parents []HashableNode
		tr.Traverse(func(node avl.Node) (bool, error) {
			if zero != nil {
				return false, nil
			}

			if avl.EqualKey(zeroKey, node.Key()) {
				zero = node
				return false, nil
			}

			var p avl.Node
			if len(parents) < 1 {
				parents = append(parents, node.(HashableNode))
				return true, nil
			}

			p = parents[len(parents)-1]

			isLeft := avl.CompareKey(zeroKey, p.Key()) < 0
			if isLeft == (avl.CompareKey(p.Key(), node.Key()) < 0) {
				return false, nil
			}

			parents = append(parents, node.(HashableNode))

			return true, nil
		})

		{ // test with *Tree.GetWithParents()
			_, parents0, _ := tr.GetWithParents(zero.Key())
			t.Equal(len(parents), len(parents0))
			for i, p := range parents {
				t.Equal(p.Key(), parents0[i].Key())
			}
		}
	}

	prover := ExampleProver{}

	zeroKey := nodes[14].Key()
	pr, err := prover.Proof(tr, zeroKey)
	t.NoError(err)

	{ // prove proof
		err := prover.Prove(pr, tr.Root().(HashableNode).Hash())
		t.NoError(err)
	}
}

func TestTree(t *testing.T) {
	suite.Run(t, new(testTree))
}
