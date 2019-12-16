package avl_test

import (
	"fmt"

	"github.com/spikeekips/avl"
)

// ExampleNode implements Node.
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

// ExampleTree shows to load the existing nodes to Tree.
func ExampleTree_withNode() {
	shape := map[int]struct {
		height int16
		left   int
		right  int
	}{
		100: {height: 3, left: 50, right: 150},
		50:  {height: 1, left: 30, right: 70},
		150: {height: 2, left: 130, right: 180},
		30:  {height: 0},
		70:  {height: 0},
		130: {height: 0},
		170: {height: 0},
		180: {height: 1, left: 170, right: 200},
		200: {height: 0},
	}

	i2k := func(i int) []byte {
		return []byte(fmt.Sprintf("%03d", i))
	}

	// NodePool will contain all the nodes.
	nodePool := avl.NewMapNodePool(nil)
	for key, properties := range shape {
		node := &ExampleNode{
			key: i2k(key),
		}
		node.height = properties.height

		if properties.left > 0 {
			node.left = i2k(properties.left)
		}
		if properties.right > 0 {
			node.right = i2k(properties.right)
		}

		_ = nodePool.Set(node)
	}

	// Trying to load nodes from NodePool
	tree, err := avl.NewTree(i2k(100), nodePool)
	if err != nil {
		fmt.Printf("error: failed to load nodes; %v", err)
		return
	}
	fmt.Println("< tree is loaded")

	fmt.Printf("tree root is %s\n", string(tree.Root().Key()))

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

	// Check Tree is valid or not
	if err := tree.IsValid(); err != nil {
		fmt.Printf("error: failed to validate; %v", err)
		return
	}
	fmt.Println("< tree is valid")
}
