package hashable

import (
	"encoding"

	"github.com/spikeekips/avl"
)

var (
	InvalidProofError = avl.NewWrapError("invalid proof")
)

type Proof interface {
	encoding.BinaryMarshaler
}

type Prover interface {
	Proof(tr *avl.Tree, key []byte) (Proof, error)
	GenerateNodeHash(HashableNode) ([]byte, error)
	Prove(Proof, rootHash []byte) error
}
