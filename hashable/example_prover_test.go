package hashable

import (
	"bytes"
	"crypto/sha256"
	"sort"
)

type ExampleProver struct {
}

func (ep ExampleProver) GenerateNodeHash(node HashableNode) ([]byte, error) {
	return ep.generateNodeHash(
		node.Key(),
		node.Height(),
		node.ValueHash(),
		node.LeftHash(),
		node.RightHash(),
	)
}

func (ep ExampleProver) generateNodeHash(
	key []byte, height int16, valueHash, leftHash, rightHash []byte,
) ([]byte, error) {
	h := sha256.New()
	_, _ = h.Write(key)
	_, _ = h.Write(int16ToBytes(height))
	if valueHash != nil {
		_, _ = h.Write(valueHash)
	}
	if leftHash != nil {
		_, _ = h.Write(leftHash)
	}
	if rightHash != nil {
		_, _ = h.Write(rightHash)
	}

	return h.Sum(nil), nil
}

func (ep ExampleProver) Proof(node HashableNode, parents []HashableNode) (Proof, error) {
	// NOTE sort by height; lower height will be first item
	sort.Slice(parents, func(i, j int) bool { return parents[i].Height() < parents[j].Height() })

	var parentProofs []ExampleParentProof
	for _, node := range parents {
		p := node.(HashableNode)
		parentProofs = append(parentProofs, ExampleParentProof{
			key:       p.Key(),
			height:    p.Height(),
			hash:      p.Hash(),
			leftHash:  p.LeftHash(),
			rightHash: p.RightHash(),
			valueHash: p.ValueHash(),
		})
	}

	n := node.(HashableNode)
	return ExampleProof{
		node: ExampleParentProof{
			key:       n.Key(),
			height:    n.Height(),
			hash:      n.Hash(),
			leftHash:  n.LeftHash(),
			rightHash: n.RightHash(),
			valueHash: n.ValueHash(),
		},
		parents: parentProofs,
	}, nil
}

func (ep ExampleProver) proveProofNode(pr ExampleParentProof) error {
	nodeHash, err := ep.generateNodeHash(
		pr.key,
		pr.height,
		pr.valueHash,
		pr.leftHash,
		pr.rightHash,
	)
	if err != nil {
		return err
	}

	if !bytes.Equal(pr.hash, nodeHash) {
		return InvalidProofError.Wrapf(
			"node hash not match: proof.hash=%v != generated=%v",
			pr.hash,
			nodeHash,
		)
	}

	return nil
}

func (ep ExampleProver) Prove(proof Proof, rootHash []byte) error {
	pr := proof.(ExampleProof)

	// check node itself: node hash
	if err := ep.proveProofNode(pr.node); err != nil {
		return err
	}

	var prRootHash []byte
	if len(pr.parents) < 1 {
		prRootHash = pr.node.hash
	} else {
		prRootHash = pr.parents[len(pr.parents)-1].hash
	}

	// check with root hash with node hash
	if !bytes.Equal(prRootHash, rootHash) {
		return InvalidProofError.Wrapf(
			"top node hash not match with root hash: topnode.hash=%v != root.hash=%v",
			prRootHash,
			rootHash,
		)
	}

	if len(pr.parents) < 1 {
		return nil
	}

	// check parents hashes
	var leaf ExampleParentProof = pr.node
	for _, p := range pr.parents {
		if !bytes.Equal(leaf.hash, p.leftHash) && !bytes.Equal(leaf.hash, p.rightHash) {
			return InvalidProofError.Wrapf(
				"node hash not match with leafs of 1st parents: leaf.hash=%v != (parents.leftHash=%v, parent.rightHash=%v)",
				pr.node.hash,
				p.leftHash,
				p.rightHash,
			)
		}

		if err := ep.proveProofNode(p); err != nil {
			return err
		}

		leaf = p
	}

	return nil
}
