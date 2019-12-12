package avl

import "github.com/rs/zerolog"

var (
	FailedToAddNodeInTreeError = NewWrapError("failed to add node to TreeGenerator")
)

type TreeGenerator struct {
	*Logger
	root  MutableNode
	nodes map[string]MutableNode
}

func NewTreeGenerator() *TreeGenerator {
	return &TreeGenerator{
		Logger: NewLogger(func(c zerolog.Context) zerolog.Context {
			return c.Str("module", "avl_tree_generator")
		}),
		nodes: map[string]MutableNode{},
	}
}

func (tg *TreeGenerator) Root() MutableNode {
	return tg.root
}

func (tg *TreeGenerator) Nodes() map[string]MutableNode {
	return tg.nodes
}

func (tg *TreeGenerator) Tree() (*Tree, error) {
	return NewTree(tg.root.Key(), NodePoolFromMutableNodeMap(tg.nodes))
}

func (tg *TreeGenerator) Add(node MutableNode) ([]MutableNode /* parents node */, error) {
	log_ := tg.Log().With().Bytes("key", node.Key()).Logger()

	_ = node.SetHeight(0)
	_ = node.SetLeft(nil)
	_ = node.SetRight(nil)

	if err := IsValidNode(node, nil, nil); err != nil {
		log_.Error().Err(err).Msg("invalid node found")
		return nil, err
	}

	if tg.root == nil {
		log_.Debug().Msg("root is empty; new node will be root")
		tg.root = node

		tg.nodes[string(node.Key())] = node

		return nil, nil
	} else if EqualKey(tg.root.Key(), node.Key()) {
		log_.Debug().Msg("same with root; root overrided")

		if err := tg.root.Merge(node); err != nil {
			return nil, err
		}

		return nil, nil
	}

	parents, err := tg.add(node)
	if err != nil {
		return nil, err
	}

	if _, found := tg.nodes[string(node.Key())]; !found {
		tg.nodes[string(node.Key())] = node
	}

	return parents, nil
}

func (tg *TreeGenerator) add(node MutableNode) ([]MutableNode /* parents node */, error) {
	log_ := tg.Log().With().Bytes("key", node.Key()).Logger()

	var parents []MutableNode
	var parent MutableNode = tg.root

	for {
		newParent, cmp, err := tg.findNode(node, parent)
		if err != nil {
			return nil, err
		}

		if cmp == 0 {
			if err := parent.Merge(node); err != nil {
				return nil, err
			}

			return parents, nil
		}

		parents = append(parents, parent)
		if newParent == nil {
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
		if tg.checkSingleViolation(p1, p2) {
			var head MutableNode
			if len(parents) > 2 {
				head = parents[len(parents)-3]
			}

			if err := tg.singleRotation(head, p2, p1, node); err != nil {
				return nil, err
			}

			return parents, nil
		}
	}

	var head, violated MutableNode
	var violatedLeft bool
	for i := len(parents) - 2; i > -1; i-- {
		// NOTE last parent already has correct height
		p := parents[i]

		if isLeft, v := isSiblingNodesViolated(p.Left(), p.Right()); v {
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
			after_height, _ := tg.resetNodeHeight(p, true)
			log_.Debug().
				Bytes("parent_key", p.Key()).
				Int16("before_height", p.Height()).
				Int16("after_height", after_height).
				Msg("reset height of parent")
		}
		if _, err := tg.resetNodeHeight(p, false); err != nil {
			return nil, err
		}
	}

	if violated == nil {
		return parents, nil
	}

	leaf := tg.getLeaf(violated, violatedLeft)
	if leaf == nil {
		return nil, FailedToAddNodeInTreeError.Wrapf(
			"leaf of violated must not be empty: violated=%v isLeft=%v",
			violated, violatedLeft,
		)
	}

	if violatedLeft == (CompareKey(node.Key(), leaf.Key()) < 0) {
		// same side rotation(left-left or right-right)
		return parents, tg.leftLeftRotation(head, violated, node, violatedLeft)
	}

	// different side(left-right or right-left)
	return parents, tg.leftRightRotation(head, violated, node, violatedLeft)
}

func (tg *TreeGenerator) getLeaf(node MutableNode, isLeft bool) MutableNode {
	if isLeft {
		return node.Left()
	}
	return node.Right()

}

func (tg *TreeGenerator) setLeaf(parent, node MutableNode, isLeft bool) error {
	if isLeft {
		return parent.SetLeft(node)
	}

	return parent.SetRight(node)
}

func (tg *TreeGenerator) resetNodeHeight(node MutableNode, dryRun bool) (int16, error) {
	var baseHeight int16
	if node.Left() != nil {
		baseHeight = node.Left().Height() + 1
	}

	if node.Right() != nil && node.Right().Height() >= baseHeight {
		baseHeight = node.Right().Height() + 1
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

func (tg *TreeGenerator) findNode(node, parent MutableNode) (
	MutableNode, /* parent */
	int, /* bytes.Compare */
	error,
) {
	log_ := tg.Log().With().Bytes("key", node.Key()).Logger()

	c := CompareKey(node.Key(), parent.Key())
	if c == 0 {
		log_.Debug().Bytes("parent_key", parent.Key()).Msg("node has same key with parent")
		return nil, c, nil
	}

	if c < 0 { // left
		if parent.Left() != nil {
			log_.Debug().Bytes("parent_key", parent.Key()).Msg("next left parent")

			return tg.getLeaf(parent, true), c, nil
		}

		if err := parent.SetLeft(node); err != nil {
			return nil, c, err
		}
		if _, err := tg.resetNodeHeight(parent, false); err != nil {
			return nil, c, err
		}

		log_.Debug().Bytes("parent_key", parent.Key()).Msg("node set at left")
		return nil, c, nil
	}

	// right
	if parent.Right() != nil {
		log_.Debug().Bytes("parent_key", parent.Key()).Msg("next right parent")

		return tg.getLeaf(parent, false), c, nil
	}

	if err := parent.SetRight(node); err != nil {
		return nil, c, err
	}
	if _, err := tg.resetNodeHeight(parent, false); err != nil {
		return nil, c, err
	}
	log_.Debug().Bytes("parent_key", parent.Key()).Msg("node set at right")

	return nil, c, nil
}

func (tg *TreeGenerator) checkSingleViolation(p1, p2 MutableNode) bool {
	var notFound bool
	for _, p := range []MutableNode{p1, p2} {
		if p.Left() != nil && p.Right() != nil {
			notFound = true
			break
		}
	}

	return !notFound
}

func (tg *TreeGenerator) singleRotation(head, p2, p1, node MutableNode) error {
	log_ := tg.Log().With().
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

	isLeft := CompareKey(p1.Key(), p2.Key()) < 0

	var top MutableNode
	if isLeft == (CompareKey(node.Key(), p1.Key()) < 0) {
		if isLeft {
			if err := p2.SetLeft(nil); err != nil {
				return err
			}
		} else {
			if err := p2.SetRight(nil); err != nil {
				return err
			}
		}
		if _, err := tg.resetNodeHeight(p2, false); err != nil {
			return err
		}

		if err := tg.setLeaf(p1, p2, !isLeft); err != nil {
			return err
		}

		top = p1
	} else {
		if err := tg.setLeaf(p2, nil, isLeft); err != nil {
			return err
		}
		if _, err := tg.resetNodeHeight(p2, false); err != nil {
			return err
		}
		if err := tg.setLeaf(p1, nil, !isLeft); err != nil {
			return err
		}
		if _, err := tg.resetNodeHeight(p1, false); err != nil {
			return err
		}
		if err := tg.setLeaf(node, p1, isLeft); err != nil {
			return err
		}
		if err := tg.setLeaf(node, p2, !isLeft); err != nil {
			return err
		}
		if _, err := tg.resetNodeHeight(node, false); err != nil {
			return err
		}

		top = node
	}

	if head == nil {
		tg.root = top
		return nil
	}

	return tg.setLeaf(head, top, CompareKey(top.Key(), head.Key()) < 0)
}

func (tg *TreeGenerator) leftLeftRotation(head, violated, node MutableNode, isLeft bool) error {
	log_ := tg.Log().With().
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

	p2 := tg.getLeaf(violated, isLeft)
	if p2 == nil {
		return FailedToAddNodeInTreeError.Wrapf(
			"leaf of violated, p2 must not be empty: violated=%v isLeft=%v",
			violated, isLeft,
		)
	}

	p3 := violated
	p2r := tg.getLeaf(p2, !isLeft)
	if p2 == nil {
		return FailedToAddNodeInTreeError.Wrapf(
			"leaf of p2, p2r must not be empty: p2=%v isLeft=%v",
			p2, !isLeft,
		)
	}

	if err := tg.setLeaf(p3, p2r, isLeft); err != nil {
		return err
	}
	if err := tg.setLeaf(p2, p3, !isLeft); err != nil {
		return err
	}
	if _, err := tg.resetNodeHeight(p3, false); err != nil {
		return err
	}

	if head == nil {
		tg.root = p2

		return nil
	}

	return tg.setLeaf(head, p2, CompareKey(p3.Key(), head.Key()) < 0)
}

func (tg *TreeGenerator) leftRightRotation(head, violated, node MutableNode, isLeft bool) error {
	log_ := tg.Log().With().
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
	p2 := tg.getLeaf(violated, isLeft)
	if p2 == nil {
		return FailedToAddNodeInTreeError.Wrapf(
			"leaf of violated must not be empty: violated=%v isLeft=%v",
			violated, isLeft,
		)
	}

	p1 := tg.getLeaf(p2, !isLeft)
	if p1 == nil {
		return FailedToAddNodeInTreeError.Wrapf(
			"leaf of p2 must not be empty: p2=%v isLeft=%v",
			p2, !isLeft,
		)
	}

	leafLeft := CompareKey(node.Key(), p1.Key()) < 0
	n0 := tg.getLeaf(p1, leafLeft)
	if n0 == nil {
		return FailedToAddNodeInTreeError.Wrapf(
			"n0, leaf of p1 must not be empty: p1=%v isLeft=%v",
			p1, leafLeft,
		)
	}
	n1 := tg.getLeaf(p1, !leafLeft)

	if isLeft != leafLeft {
		if err := tg.setLeaf(p2, n1, !isLeft); err != nil {
			return err
		}
		if _, err := tg.resetNodeHeight(p2, false); err != nil {
			return err
		}
		if err := tg.setLeaf(p3, n0, isLeft); err != nil {
			return err
		}
		if _, err := tg.resetNodeHeight(p3, false); err != nil {
			return err
		}
	} else {
		if err := tg.setLeaf(p3, n1, isLeft); err != nil {
			return err
		}
		if _, err := tg.resetNodeHeight(p3, false); err != nil {
			return err
		}
		if err := tg.setLeaf(p2, n0, !isLeft); err != nil {
			return err
		}
		if _, err := tg.resetNodeHeight(p2, false); err != nil {
			return err
		}
	}

	if err := tg.setLeaf(p1, p2, isLeft); err != nil {
		return err
	}
	if err := tg.setLeaf(p1, p3, !isLeft); err != nil {
		return err
	}
	if _, err := tg.resetNodeHeight(p1, false); err != nil {
		return err
	}

	if head == nil {
		tg.root = p1
		return nil
	}

	return tg.setLeaf(head, p1, CompareKey(p3.Key(), head.Key()) < 0)
}
