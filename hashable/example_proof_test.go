package hashable

import "encoding/json"

type ExampleParentProof struct {
	key       []byte
	height    int16
	valueHash []byte
	leftHash  []byte
	rightHash []byte
	hash      []byte
}

func (pp ExampleParentProof) MarshalBinary() ([]byte, error) {
	return json.Marshal(pp)
}

func (pp ExampleParentProof) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"key":        pp.key,
		"height":     pp.height,
		"hash":       pp.hash,
		"left_hash":  pp.leftHash,
		"right_hash": pp.rightHash,
		"value_hash": pp.valueHash,
	})
}

type ExampleProof struct {
	node    ExampleParentProof
	parents []ExampleParentProof
}

func (pr ExampleProof) String() string {
	b, _ := json.Marshal(pr)

	return string(b)
}

func (pr ExampleProof) MarshalBinary() ([]byte, error) {
	return json.Marshal(pr)
}

func (pr ExampleProof) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"node":    pr.node,
		"parents": pr.parents,
	})
}
