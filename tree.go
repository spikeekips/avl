package avl

import (
	"bytes"

	"github.com/rs/zerolog"
)

var (
	NodeNotFound = NewWrapError("node not found")
)

type Tree struct {
	*Logger
	root     Node
	nodePool NodePool
}

func NewTree(nodePool NodePool) *Tree {
	return &Tree{
		Logger: NewLogger(func(c zerolog.Context) zerolog.Context {
			return c.Str("module", "avl_tree")
		}),
		nodePool: nodePool,
	}
}

func (tr *Tree) Root() Node {
	return tr.root
}

func (tr *Tree) SetRoot(root Node) *Tree {
	tr.root = root

	return tr
}

func (tr *Tree) NodePool() NodePool {
	return tr.nodePool
}

func (tr *Tree) Get(key []byte) (Node, error) {
	log_ := tr.Log().With().Bytes("key", key).Logger()

	if tr.root == nil {
		log_.Debug().Msg("empty tree")
		return nil, nil
	}

	var parent Node
	parent = tr.root

	var depth int
	for {
		c := bytes.Compare(key, parent.Key())

		if c == 0 {
			log_.Debug().Int("depth", depth).Msg("found node by key")
			return parent, nil
		}

		var err error
		if parent, err = tr.getLeaf(parent, c < 0); err != nil {
			return nil, err
		}
		if parent == nil {
			break
		}
		depth++
	}

	return nil, nil
}

func (tr *Tree) GetWithParents(key []byte) (Node, []Node, error) {
	log_ := tr.Log().With().Bytes("key", key).Logger()

	if tr.root == nil {
		log_.Debug().Msg("empty tree")
		return nil, nil, nil
	}

	var parents []Node
	parent := tr.root

	var depth int
	for {
		c := bytes.Compare(key, parent.Key())

		if c == 0 {
			log_.Debug().Int("depth", depth).Msg("found node by key")
			return parent, parents, nil
		}
		parents = append(parents, parent)

		var err error
		if parent, err = tr.getLeaf(parent, c < 0); err != nil {
			return nil, nil, err
		}
		if parent == nil {
			break
		}
		depth++
	}

	return nil, nil, nil
}

type TraverseFunc func(Node) (bool, error)

func (tr *Tree) Traverse(f TraverseFunc) error {
	if tr.root == nil {
		return nil
	}

	err, _ := tr.traverse(tr.root, f)

	return err
}

func (tr *Tree) traverse(node Node, f TraverseFunc) (error, bool) {
	if keep, err := f(node); !keep || err != nil {
		return err, false
	}

	if node.LeftKey() != nil {
		left, err := tr.getLeafEnsure(node, true)
		if err != nil {
			return err, false
		}
		if err, keep := tr.traverse(left, f); err != nil || !keep {
			return err, keep
		}
	}
	if node.RightKey() != nil {
		right, err := tr.getLeafEnsure(node, false)
		if err != nil {
			return err, false
		}
		if err, keep := tr.traverse(right, f); err != nil || !keep {
			return err, keep
		}
	}

	return nil, true
}

func (tr *Tree) Add(node Node) ([]Node /* updated Node */, error) {
	log_ := tr.Log().With().Bytes("key", node.Key()).Logger()

	_ = node.SetHeight(0)
	_ = node.SetLeftKey(nil)
	_ = node.SetRightKey(nil)
	if err := IsValidNode(node, nil, nil); err != nil {
		log_.Error().Err(err).Msg("invalid node found")
		return nil, err
	}
	if err := tr.nodePool.Set(node); err != nil {
		return nil, err
	}

	if tr.root == nil {
		log_.Debug().Msg("root is empty; new node will be root")
		tr.root = node
		return nil, nil
	}

	var parent Node = tr.root
	var parents []Node

	for {
		newParent, err := tr.findNode(node, parent)
		if CompareNode(node, parent) != 0 {
			parents = append(parents, parent)
		}
		if err != nil {
			return nil, err
		} else if newParent == nil {
			break
		}

		parent = newParent
	}

	log_.Debug().Int("parents", len(parents)).Msg("found parents")

	if len(parents) < 2 {
		log_.Debug().Msg("not enough parents for rotation; done")
		return parents, nil
	}

	// check single rotation
	{
		p1 := parents[len(parents)-1]
		p2 := parents[len(parents)-2]
		if tr.checkSingleViolation(p1, p2) {
			var head Node
			if len(parents) > 2 {
				head = parents[len(parents)-3]
			}

			if err := tr.singleRotation(head, p2, p1, node); err != nil {
				return nil, err
			}

			return parents, nil
		}
	}

	var head, violated Node
	var violatedLeft bool
	for i := len(parents) - 2; i > -1; i-- {
		// NOTE last parent already has correct height
		p := parents[i]

		left, right, err := tr.getLeafs(p)
		if err != nil {
			return nil, err
		}

		if isLeft, v := isViolated(left, right); v {
			violated = p
			violatedLeft = isLeft

			var head_key []byte
			if i > 0 {
				head = parents[i-1]
				head_key = head.Key()
			}

			log_.Debug().
				Bytes("head_key", head_key).
				Bytes("violated_key", violated.Key()).
				Msg("violated found")

			break
		}

		if log_.GetLevel() == zerolog.DebugLevel {
			after_height, _ := tr.resetNodeHeight(p, true)
			log_.Debug().
				Bytes("parent_key", p.Key()).
				Int16("before_height", p.Height()).
				Int16("after_height", after_height).
				Msg("reset height of parent")
		}
		if _, err = tr.resetNodeHeight(p, false); err != nil {
			return nil, err
		}
	}

	if violated == nil {
		return parents, nil
	}

	leaf, err := tr.getLeafEnsure(violated, violatedLeft)
	if err != nil {
		return nil, err
	}

	if violatedLeft == (CompareNode(node, leaf) < 0) {
		// same side rotation(left-left or right-right)
		return parents, tr.leftLeftRotation(head, violated, node, violatedLeft)
	}

	// different side(left-right or right-left)
	return parents, tr.leftRightRotation(head, violated, node, violatedLeft)
}

func (tr *Tree) findNode(node, parent Node) (Node /* parent */, error) {
	log_ := tr.Log().With().Bytes("key", node.Key()).Logger()

	c := CompareNode(node, parent)
	if c == 0 {
		log_.Debug().Bytes("parent_key", parent.Key()).Msg("node has same key with parent")
		return nil, nil
	}

	if c < 0 { // left
		if parent.LeftKey() != nil {
			log_.Debug().Bytes("parent_key", parent.Key()).Msg("next left parent")

			var err error
			newParent, err := tr.getLeaf(parent, true)
			if err != nil {
				return nil, err
			}

			return newParent, nil
		}

		if err := parent.SetLeftKey(node.Key()); err != nil {
			return nil, err
		}
		if _, err := tr.resetNodeHeight(parent, false); err != nil {
			return nil, err
		}

		log_.Debug().Bytes("parent_key", parent.Key()).Msg("node set at left")
		return nil, nil
	}

	// right
	if parent.RightKey() != nil {
		log_.Debug().Bytes("parent_key", parent.Key()).Msg("next right parent")

		newParent, err := tr.getLeaf(parent, false)
		if err != nil {
			return nil, err
		}

		return newParent, nil
	}

	if err := parent.SetRightKey(node.Key()); err != nil {
		return nil, err
	}
	if _, err := tr.resetNodeHeight(parent, false); err != nil {
		return nil, err
	}
	log_.Debug().Bytes("parent_key", parent.Key()).Msg("node set at right")

	return nil, nil
}

func (tr *Tree) checkSingleViolation(p1, p2 Node) bool {
	var notFound bool
	for _, p := range []Node{p1, p2} {
		if p.LeftKey() != nil && p.RightKey() != nil {
			notFound = true
			break
		}
	}

	return !notFound
}

func (tr *Tree) getLeaf(parent Node, isLeft bool) (Node, error) {
	var key []byte
	if isLeft {
		key = parent.LeftKey()
	} else {
		key = parent.RightKey()
	}
	if key == nil {
		return nil, nil
	}

	return tr.nodePool.Get(key)
}

func (tr *Tree) getLeafEnsure(parent Node, isLeft bool) (Node, error) {
	node, err := tr.getLeaf(parent, isLeft)
	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, NodeNotFound.Wrapf(
			"leaf node not found in parent, %v; isLeft=%v",
			parent.Key(),
			isLeft,
		)
	}

	return node, nil
}

func (tr *Tree) getLeafs(node Node) (Node, Node, error) {
	var err error
	var left, right Node
	if node.LeftKey() != nil {
		if left, err = tr.getLeaf(node, true); err != nil {
			return nil, nil, err
		}
	}
	if node.RightKey() != nil {
		if right, err = tr.getLeaf(node, false); err != nil {
			return nil, nil, err
		}
	}

	return left, right, nil
}

func (tr *Tree) setLeaf(parent, node Node, isLeft bool) error {
	var key []byte
	if node == nil {
		key = nil
	} else {
		key = node.Key()
	}

	if isLeft {
		return parent.SetLeftKey(key)
	}

	return parent.SetRightKey(key)
}

func (tr *Tree) singleRotation(head, p2, p1, node Node) error {
	log_ := tr.Log().With().
		Bytes("key", node.Key()).
		Logger()

	if log_.GetLevel() == zerolog.DebugLevel {
		var head_key []byte
		if head != nil {
			head_key = head.Key()
		}

		log_.Debug().
			Bytes("head_key", head_key).
			Bytes("p2_key", p2.Key()).
			Msg("found single rotation")
	}

	isLeft := CompareNode(p1, p2) < 0

	var top Node
	if isLeft == (CompareNode(node, p1) < 0) {
		if isLeft {
			if err := p2.SetLeftKey(nil); err != nil {
				return err
			}
		} else {
			if err := p2.SetRightKey(nil); err != nil {
				return err
			}
		}
		if _, err := tr.resetNodeHeight(p2, false); err != nil {
			return err
		}

		if err := tr.setLeaf(p1, p2, !isLeft); err != nil {
			return err
		}

		top = p1
	} else {
		if err := tr.setLeaf(p2, nil, isLeft); err != nil {
			return err
		}
		if _, err := tr.resetNodeHeight(p2, false); err != nil {
			return err
		}
		if err := tr.setLeaf(p1, nil, !isLeft); err != nil {
			return err
		}
		if _, err := tr.resetNodeHeight(p1, false); err != nil {
			return err
		}
		if err := tr.setLeaf(node, p1, isLeft); err != nil {
			return err
		}
		if err := tr.setLeaf(node, p2, !isLeft); err != nil {
			return err
		}
		if _, err := tr.resetNodeHeight(node, false); err != nil {
			return err
		}

		top = node
	}

	if head == nil {
		tr.root = top
	} else {
		if err := tr.setLeaf(head, top, CompareNode(top, head) < 0); err != nil {
			return err
		}
	}

	return nil
}

func (tr *Tree) leftLeftRotation(head, violated, node Node, isLeft bool) error {
	log_ := tr.Log().With().
		Bytes("key", node.Key()).
		Logger()

	if log_.GetLevel() == zerolog.DebugLevel {
		var head_key []byte
		if head != nil {
			head_key = head.Key()
		}

		log__ := log_.With().
			Bytes("head_key", head_key).
			Bytes("violated_key", violated.Key()).
			Logger()

		if isLeft {
			log__.Debug().Msg("found left-left rotation")
		} else {
			log__.Debug().Msg("found right-right rotation")
		}
	}

	p2, err := tr.getLeafEnsure(violated, isLeft)
	if err != nil {
		return err
	}

	p3 := violated
	p2r, err := tr.getLeafEnsure(p2, !isLeft)
	if err != nil {
		return err
	}

	if err := tr.setLeaf(p3, p2r, isLeft); err != nil {
		return err
	}
	if err := tr.setLeaf(p2, p3, !isLeft); err != nil {
		return err
	}
	if _, err := tr.resetNodeHeight(p3, false); err != nil {
		return err
	}

	if head == nil {
		tr.root = p2
	} else {
		if err := tr.setLeaf(head, p2, CompareNode(p3, head) < 0); err != nil {
			return err
		}
	}

	return nil
}

func (tr *Tree) leftRightRotation(head, violated, node Node, isLeft bool) error {
	log_ := tr.Log().With().
		Bytes("key", node.Key()).
		Logger()

	if log_.GetLevel() == zerolog.DebugLevel {
		var head_key []byte
		if head != nil {
			head_key = head.Key()
		}

		log__ := log_.With().
			Bytes("head_key", head_key).
			Bytes("violated_key", violated.Key()).
			Logger()
		if isLeft {
			log__.Debug().Msg("found left-right rotation")
		} else {
			log__.Debug().Msg("found right-left rotation")
		}
	}

	p3 := violated
	p2, err := tr.getLeafEnsure(violated, isLeft)
	if err != nil {
		return err
	}
	p1, err := tr.getLeafEnsure(p2, !isLeft)
	if err != nil {
		return err
	}

	leafLeft := CompareNode(node, p1) < 0
	n0, err := tr.getLeafEnsure(p1, leafLeft)
	if err != nil {
		return err
	}
	n1, err := tr.getLeaf(p1, !leafLeft)
	if err != nil {
		return err
	}

	if isLeft != leafLeft {
		if err := tr.setLeaf(p2, n1, !isLeft); err != nil {
			return err
		}
		if _, err := tr.resetNodeHeight(p2, false); err != nil {
			return err
		}
		if err := tr.setLeaf(p3, n0, isLeft); err != nil {
			return err
		}
		if _, err := tr.resetNodeHeight(p3, false); err != nil {
			return err
		}
	} else {
		if err := tr.setLeaf(p3, n1, isLeft); err != nil {
			return err
		}
		if _, err := tr.resetNodeHeight(p3, false); err != nil {
			return err
		}
		if err := tr.setLeaf(p2, n0, !isLeft); err != nil {
			return err
		}
		if _, err := tr.resetNodeHeight(p2, false); err != nil {
			return err
		}
	}

	if err := tr.setLeaf(p1, p2, isLeft); err != nil {
		return err
	}
	if err := tr.setLeaf(p1, p3, !isLeft); err != nil {
		return err
	}
	if _, err := tr.resetNodeHeight(p1, false); err != nil {
		return err
	}

	if head == nil {
		tr.root = p1
		return nil
	}

	return tr.setLeaf(head, p1, CompareNode(p3, head) < 0)
}

func (tr *Tree) IsValid() error {
	if tr.root == nil {
		return nil
	}

	return tr.validate(tr.root, nil)
}

func (tr *Tree) validate(node Node, parents []Node) error {
	log_ := tr.Log().With().Int("parents", len(parents)).Bytes("key", node.Key()).Logger()

	left, right, err := tr.getLeafs(node)
	if err != nil {
		return err
	}

	if err := IsValidNode(node, left, right); err != nil {
		log_.Error().Err(err).Msg("failed to validate node")
		return err
	}

	if left != nil {
		if err := tr.validate(left, append(parents, node)); err != nil {
			return err
		}
	}
	if right != nil {
		if err := tr.validate(right, append(parents, node)); err != nil {
			return err
		}
	}

	log_.Debug().Msg("validated")

	return nil
}

func (tr *Tree) resetNodeHeight(node Node, dryRun bool) (int16, error) {
	left, right, err := tr.getLeafs(node)
	if err != nil {
		return 0, err
	}

	var baseHeight int16
	if left != nil {
		baseHeight = left.Height() + 1
	}

	if right != nil && right.Height() >= baseHeight {
		baseHeight = right.Height() + 1
	}
	if baseHeight == node.Height() {
		return baseHeight, nil
	}

	if dryRun {
		return baseHeight, nil
	}

	if err := node.SetHeight(baseHeight); err != nil {
		return 0, err
	}

	return baseHeight, nil
}

func isViolated(a, b Node) (bool /* left(true) or right(false) violated */, bool /* violated */) {
	if a == nil || b == nil {
		return false, false
	}

	d := a.Height() - b.Height()
	if d < 2 && d > -2 {
		return false, false
	}

	return d > 1, true
}
