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

	_ = tr.Traverse(func(node avl.Node) (bool, error) {
		t.Nil(node.(*ExampleHashableMutableNode).Hash())
		return true, nil
	})

	prover := ExampleProver{}

	err = SetTreeNodeHash(tr.Root().(HashableMutableNode), prover.GenerateNodeHash)
	t.NoError(err)

	rootHash := tg.Root().(HashableNode).Hash()
	t.NotNil(rootHash)

	t.Equal(nodes[3].Hash(), rootHash)

	// hash was generated
	_ = tr.Traverse(func(node avl.Node) (bool, error) {
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

	var zero HashableNode
	var parents []HashableNode
	{
		// find parents of 0 height node, 14
		zeroKey := nodes[14].Key()
		var ps []avl.Node
		_ = tr.Traverse(func(node avl.Node) (bool, error) {
			if zero != nil {
				return false, nil
			}

			if avl.EqualKey(zeroKey, node.Key()) {
				zero = node.(HashableNode)
				return false, nil
			}

			var p avl.Node
			if len(ps) < 1 {
				ps = append(ps, node.(HashableNode))
				return true, nil
			}

			p = ps[len(ps)-1]

			isLeft := avl.CompareKey(zeroKey, p.Key()) < 0
			if isLeft == (avl.CompareKey(p.Key(), node.Key()) < 0) {
				return false, nil
			}

			ps = append(ps, node.(HashableNode))

			return true, nil
		})

		// test with *Tree.GetWithParents()
		_, pst, _ := tr.GetWithParents(zero.Key())
		t.Equal(len(ps), len(pst))
		for i, p := range ps {
			t.Equal(p.Key(), pst[i].Key())

			parents = append(parents, p.(HashableNode))
		}
	}

	prover := ExampleProver{}

	// all nodes must has nil hash
	_ = tr.Traverse(func(node avl.Node) (bool, error) {
		t.Nil(node.(HashableNode).Hash())

		return true, nil
	})

	err := SetTreeNodeHash(
		tr.Root().(HashableMutableNode),
		prover.GenerateNodeHash,
	)
	t.NoError(err)

	// check hash is empty or not
	_ = tr.Traverse(func(node avl.Node) (bool, error) {
		t.NotNil(node.(HashableNode).Hash())

		return true, nil
	})

	pr, err := prover.Proof(zero, parents)
	t.NoError(err)

	{ // prove proof
		err := prover.Prove(pr, tr.Root().(HashableNode).Hash())
		t.NoError(err)
	}
}

func TestTree(t *testing.T) {
	suite.Run(t, new(testTree))
}
