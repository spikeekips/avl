package hashable

import (
	"encoding/binary"

	"golang.org/x/xerrors"

	"github.com/spikeekips/avl"
)

type ExampleHashableMutableNode struct {
	key    []byte
	height int16
	left   HashableMutableNode
	right  HashableMutableNode
	value  int
	hash   []byte
}

func (eh *ExampleHashableMutableNode) Key() []byte {
	return eh.key
}

func (eh *ExampleHashableMutableNode) Height() int16 {
	return eh.height
}

func (eh *ExampleHashableMutableNode) SetHeight(height int16) error {
	if height < 0 {
		return xerrors.Errorf("height must be greater than zero; height=%d", height)
	}

	eh.height = height

	return nil
}

func (eh *ExampleHashableMutableNode) Left() avl.MutableNode {
	return eh.left
}

func (eh *ExampleHashableMutableNode) LeftKey() []byte {
	if eh.left == nil {
		return nil
	}

	return eh.left.Key()
}

func (eh *ExampleHashableMutableNode) SetLeft(node avl.MutableNode) error {
	if node == nil {
		eh.left = nil
		return nil
	}

	m, ok := node.(HashableMutableNode)
	if !ok {
		return xerrors.Errorf("not HashableMutableNode; %T", node)
	}

	if avl.EqualKey(eh.key, node.Key()) {
		return xerrors.Errorf("left is same node; key=%v", eh.key)
	}

	eh.left = m

	return nil
}

func (eh *ExampleHashableMutableNode) Right() avl.MutableNode {
	return eh.right
}

func (eh *ExampleHashableMutableNode) RightKey() []byte {
	if eh.right == nil {
		return nil
	}

	return eh.right.Key()
}

func (eh *ExampleHashableMutableNode) SetRight(node avl.MutableNode) error {
	if node == nil {
		eh.right = nil
		return nil
	}

	m, ok := node.(HashableMutableNode)
	if !ok {
		return xerrors.Errorf("not HashableMutableNode; %T", node)
	}

	if avl.EqualKey(eh.key, node.Key()) {
		return xerrors.Errorf("right is same node; key=%v", eh.key)
	}

	eh.right = m

	return nil
}

func (eh *ExampleHashableMutableNode) Merge(node avl.MutableNode) error {
	e, ok := node.(*ExampleHashableMutableNode)
	if !ok {
		return xerrors.Errorf("merge node is not HashableMutableNode; node=%T", node)
	}

	eh.value = e.value

	return nil
}

func (eh *ExampleHashableMutableNode) Hash() []byte {
	if eh.hash != nil {
		return eh.hash
	}

	eh.hash = ExampleProver{}.GenerateNodeHash(eh)

	return eh.hash
}

func (eh *ExampleHashableMutableNode) LeftHash() []byte {
	if eh.left == nil {
		return nil
	}

	return eh.left.Hash()
}

func (eh *ExampleHashableMutableNode) RightHash() []byte {
	if eh.right == nil {
		return nil
	}

	return eh.right.Hash()
}

func (eh *ExampleHashableMutableNode) ValueHash() []byte {
	return int64ToBytes(int64(eh.value))
}

type ExampleHashableNode struct {
	key       []byte
	height    int16
	leftKey   []byte
	rightKey  []byte
	value     int
	hash      []byte
	leftHash  []byte
	rightHash []byte
	valueHash []byte
}

func (eh ExampleHashableNode) Key() []byte {
	return eh.key
}

func (eh ExampleHashableNode) Height() int16 {
	return eh.height
}

func (eh ExampleHashableNode) LeftKey() []byte {
	return eh.leftKey
}

func (eh ExampleHashableNode) RightKey() []byte {
	return eh.rightKey
}

func (eh ExampleHashableNode) Hash() []byte {
	return eh.hash
}

func (eh ExampleHashableNode) LeftHash() []byte {
	return eh.leftHash
}

func (eh ExampleHashableNode) RightHash() []byte {
	return eh.rightHash
}

func (eh ExampleHashableNode) ValueHash() []byte {
	return eh.valueHash
}

func int64ToBytes(i int64) []byte {
	bs := make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(bs, i)

	return bs
}

func int16ToBytes(i int16) []byte {
	bs := make([]byte, binary.MaxVarintLen16)
	binary.PutVarint(bs, int64(i))

	return bs
}
