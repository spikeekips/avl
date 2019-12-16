package hashable

import (
	"encoding"

	"github.com/spikeekips/avl"
	"golang.org/x/xerrors"
)

var (
	InvalidProofError = avl.NewWrapError("invalid proof")
)

type Proof interface {
	encoding.BinaryMarshaler
}

type Prover interface {
	Proof(node HashableNode, parents []HashableNode) (Proof, error)
	GenerateNodeHash(HashableNode) ([]byte, error)
	Prove(Proof, rootHash []byte) error
}

type NodeHashFunc func(HashableNode) ([]byte, error)

func SetTreeNodeHash(node HashableMutableNode, hashFunc NodeHashFunc) error {
	if node.LeftKey() != nil && node.LeftHash() == nil {
		if mh, ok := node.Left().(HashableMutableNode); !ok {
			return xerrors.Errorf("not HashableMutableNode")
		} else if err := SetTreeNodeHash(mh, hashFunc); err != nil {
			return err
		}
	}

	if node.RightKey() != nil && node.RightHash() == nil {
		if mh, ok := node.Right().(HashableMutableNode); !ok {
			return xerrors.Errorf("not HashableMutableNode")
		} else if err := SetTreeNodeHash(mh, hashFunc); err != nil {
			return err
		}
	}

	h, err := hashFunc(node)
	if err != nil {
		return err
	}

	return node.SetHash(h)
}
