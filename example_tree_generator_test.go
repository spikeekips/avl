package avl_test

import (
	"fmt"

	"github.com/spikeekips/avl"
	"golang.org/x/xerrors"
)

// ExampleMutableNode implements MutableNode. ExampleMutableNode has value field
// to store custom value.
type ExampleMutableNode struct {
	key    []byte
	height int16
	left   avl.MutableNode
	right  avl.MutableNode
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

func (em *ExampleMutableNode) Left() avl.MutableNode {
	return em.left
}

func (em *ExampleMutableNode) LeftKey() []byte {
	if em.left == nil {
		return nil
	}

	return em.left.Key()
}

func (em *ExampleMutableNode) SetLeft(node avl.MutableNode) error {
	if node == nil {
		em.left = nil
		return nil
	}

	if avl.EqualKey(em.key, node.Key()) {
		return xerrors.Errorf("left is same node; key=%v", em.key)
	}

	em.left = node

	return nil
}

func (em *ExampleMutableNode) Right() avl.MutableNode {
	return em.right
}

func (em *ExampleMutableNode) RightKey() []byte {
	if em.right == nil {
		return nil
	}

	return em.right.Key()
}

func (em *ExampleMutableNode) SetRight(node avl.MutableNode) error {
	if node == nil {
		em.right = nil
		return nil
	}

	if avl.EqualKey(em.key, node.Key()) {
		return xerrors.Errorf("right is same node; key=%v", em.key)
	}

	em.right = node

	return nil
}

func (em *ExampleMutableNode) Merge(node avl.MutableNode) error {
	e, ok := node.(*ExampleMutableNode)
	if !ok {
		return xerrors.Errorf("merge node is not *ExampleMutableNode; node=%T", node)
	}

	em.value = e.value

	return nil
}

func ExampleTreeGenerator() {
	// create new TreeGenerator
	tg := avl.NewTreeGenerator()

	fmt.Println("> trying to generate new tree")

	// Generate 10 new MutableNodes and add to TreeGenerator.
	for i := 0; i < 10; i++ {
		node := &ExampleMutableNode{
			key: []byte(fmt.Sprintf("%03d", i)),
		}
		if _, err := tg.Add(node); err != nil {
			fmt.Printf("error: failed to add node: %v\n", err)
			return
		}
		fmt.Printf("> node added: key=%s\n", string(node.Key()))
	}

	// Get Tree from TreeGenerator.
	tree, err := tg.Tree()
	if err != nil {
		fmt.Printf("error: failed to get Tree from generator: %v\n", err)
		return
	}

	// Check all the nodes is added in Tree.
	var i int
	_ = tree.Traverse(func(node avl.Node) (bool, error) {
		fmt.Printf(
			"%d: node loaded: key=%s height=%d left=%s right=%s\n",
			i,
			string(node.Key()),
			node.Height(),
			string(node.LeftKey()),
			string(node.RightKey()),
		)
		i++
		return true, nil
	})

	// check Tree is valid or not
	if err := tree.IsValid(); err != nil {
		fmt.Printf("tree is invalid: %v\n", err)
		return
	}
	fmt.Println("< tree is valid")
}
