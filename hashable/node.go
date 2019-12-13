package hashable

import (
	"bytes"
	"crypto/sha256"
	"encoding"
	"encoding/gob"
	"fmt"

	"github.com/spikeekips/avl"
)

var (
	NotHashableNodeError = avl.NewWrapError("not HashableNode")
)

type HashableNode interface {
	encoding.BinaryMarshaler
	avl.Node

	RawHash() []byte
	Hash() []byte
	LeftHash() []byte
	RightHash() []byte
	SetLeftHash([]byte)
	SetRightHash([]byte)
	ResetHash()
}

type BaseHashableNode struct {
	*avl.BaseNode
	hash      []byte
	leftHash  []byte
	rightHash []byte
}

func NewBaseHashableNode(node *avl.BaseNode) *BaseHashableNode {
	return &BaseHashableNode{BaseNode: node}
}

type baseHashableNodeMarshal struct {
	Hash      []byte
	Key       []byte
	LeftKey   []byte
	RightKey  []byte
	Height    int16
	LeftHash  []byte
	RightHash []byte
}

func (hn *BaseHashableNode) MarshalBinary() ([]byte, error) {
	hm := baseHashableNodeMarshal{
		Key:       hn.Key(),
		LeftKey:   hn.LeftKey(),
		RightKey:  hn.RightKey(),
		Height:    hn.Height(),
		LeftHash:  hn.LeftHash(),
		RightHash: hn.RightHash(),
	}

	var b bytes.Buffer
	if err := gob.NewEncoder(&b).Encode(hm); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (hn *BaseHashableNode) UnmarshalBinary(data []byte) error {
	var hm baseHashableNodeMarshal
	if err := gob.NewDecoder(bytes.NewBuffer(data)).Decode(&hm); err != nil {
		return err
	}
	if hm.Key == nil || len(hm.Key) < 1 {
		return avl.InvalidNodeError.Wrapf("empty key")
	}
	if hm.Height < 0 {
		return avl.InvalidNodeError.Wrapf("invalid height, %v", hm.Height)
	}
	if hm.LeftKey != nil && (hm.LeftHash == nil || len(hm.LeftHash) < 1) {
		return avl.InvalidNodeError.Wrapf("left key exists, but left hash empty, %v", hm.LeftKey)
	}
	if hm.RightKey != nil && (hm.RightHash == nil || len(hm.RightHash) < 1) {
		return avl.InvalidNodeError.Wrapf("right key exists, but right hash empty, %v", hm.RightKey)
	}

	hn.BaseNode = avl.NewBaseNode(hm.Key)
	if hm.Hash != nil {
		hn.hash = hm.Hash
	}
	_ = hn.BaseNode.SetLeftKey(hm.LeftKey)
	_ = hn.BaseNode.SetRightKey(hm.RightKey)
	if err := hn.BaseNode.SetHeight(hm.Height); err != nil {
		return err
	}
	hn.leftHash = hm.LeftHash
	hn.rightHash = hm.RightHash

	return nil
}

func (hn *BaseHashableNode) Hash() []byte {
	if hn.hash != nil {
		return hn.hash
	}

	b, err := hn.MarshalBinary()
	if err != nil {
		return nil
	}

	h := sha256.New()
	_, _ = h.Write(b)
	hn.hash = h.Sum(nil)

	return hn.hash
}

func (hn *BaseHashableNode) String() string {
	var left_key, right_key []byte
	var left_hash, right_hash []byte
	if hn.LeftKey() != nil {
		left_key = hn.LeftKey()
		left_hash = hn.LeftHash()
	}
	if hn.RightKey() != nil {
		right_key = hn.RightKey()
		right_hash = hn.RightHash()
	}
	return fmt.Sprintf(
		"<Node(hashable): key=%s heihgt=%d left_key=%v right_key=%v left_hash=%v right_hash=%v>",
		hn.Key(),
		hn.Height(),
		left_key,
		right_key,
		left_hash,
		right_hash,
	)
}

func (hn *BaseHashableNode) RawHash() []byte {
	return hn.hash
}

func (hn *BaseHashableNode) ResetHash() {
	hn.hash = nil
}

func (hn *BaseHashableNode) SetHeight(h int16) error {
	if hn.Height() != h {
		hn.ResetHash()
	}

	return hn.BaseNode.SetHeight(h)
}

func (hn *BaseHashableNode) SetLeftKey(key []byte) error {
	if !bytes.Equal(hn.LeftKey(), key) {
		hn.ResetHash()
	}

	return hn.BaseNode.SetLeftKey(key)
}

func (hn *BaseHashableNode) SetRightKey(key []byte) error {
	if !bytes.Equal(hn.RightKey(), key) {
		hn.ResetHash()
	}

	return hn.BaseNode.SetRightKey(key)
}

func (hn *BaseHashableNode) LeftHash() []byte {
	return hn.leftHash
}

func (hn *BaseHashableNode) RightHash() []byte {
	return hn.rightHash
}

func (hn *BaseHashableNode) SetLeftHash(h []byte) {
	if !bytes.Equal(hn.leftHash, h) {
		hn.ResetHash()
	}

	hn.leftHash = h
}

func (hn *BaseHashableNode) SetRightHash(h []byte) {
	if !bytes.Equal(hn.rightHash, h) {
		hn.ResetHash()
	}

	hn.rightHash = h
}
