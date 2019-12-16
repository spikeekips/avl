package hashable

import (
	"github.com/spikeekips/avl"
)

type baseHashableNode interface {
	Hash() []byte
	LeftHash() []byte
	RightHash() []byte
	ValueHash() []byte
}

type HashableNode interface {
	avl.Node
	baseHashableNode
}

type HashableMutableNode interface {
	avl.MutableNode
	baseHashableNode
	SetHash([]byte) error
	ResetHash()
}
