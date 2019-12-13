package avl

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.org/x/xerrors"
)

type testBaseNode struct {
	suite.Suite
}

func (t *testBaseNode) newNode(i int) *BaseNode {
	return NewBaseNode(
		[]byte(fmt.Sprintf("%03d", i)),
	)
}

func (t *testBaseNode) TestNew() {
	key := []byte("showme")
	n := NewBaseNode(key)

	t.Equal(key, n.Key())
}

func (t *testBaseNode) TestLeftByKey() {
	n := t.newNode(10)
	l := t.newNode(5)
	err := n.SetLeftKey(l.Key())
	t.NoError(err)
	t.Equal(n.left, l.Key())
}

func (t *testBaseNode) TestRightByKey() {
	n := t.newNode(10)
	l := t.newNode(20)
	err := n.SetRightKey(l.Key())
	t.NoError(err)
	t.Equal(n.right, l.Key())
}

func (t *testBaseNode) TestLeftByKeySameKey() {
	{ // same with parent
		n := t.newNode(10)
		err := n.SetLeftKey(n.Key())
		t.NoError(err)
	}

	{ // same with existing left
		n := t.newNode(10)
		l := t.newNode(5)
		err := n.SetLeftKey(l.Key())
		t.NoError(err)

		err = n.SetLeftKey(l.Key())
		t.NoError(err)
	}
}

func (t *testBaseNode) TestRightByKeySameKey() {
	{ // same with parent
		n := t.newNode(10)
		err := n.SetRightKey(n.Key())
		t.NoError(err)
	}

	{ // same with existing right
		n := t.newNode(10)
		r := t.newNode(20)
		err := n.SetRightKey(r.Key())
		t.NoError(err)

		err = n.SetRightKey(r.Key())
		t.NoError(err)
	}
}

func (t *testBaseNode) TestInvalid() {
	{ // ok
		n := t.newNode(10)

		err := IsValidNode(n, nil, nil)
		t.NoError(err)
	}

	{ // key is nil
		n := t.newNode(10)
		n.key = nil

		err := IsValidNode(n, nil, nil)
		t.NotNil(err)
		t.True(xerrors.Is(err, InvalidNodeError))
	}

	{ // key is empty
		n := t.newNode(10)
		n.key = []byte{}

		err := IsValidNode(n, nil, nil)
		t.NotNil(err)
		t.True(xerrors.Is(err, InvalidNodeError))
	}

	{ // wrong height
		n := t.newNode(10)
		l := t.newNode(5)
		_ = n.SetLeftKey(l.Key())
		r := t.newNode(15)
		_ = r.SetHeight(5)
		_ = n.SetRightKey(r.Key())

		_ = n.SetHeight(r.Height() + 1)

		err := IsValidNode(n, l, r)
		t.NoError(err)

		_ = n.SetHeight(3)

		err = IsValidNode(n, l, r)
		t.NotNil(err)
		t.True(xerrors.Is(err, InvalidNodeError))
	}
}

func TestBaseNode(t *testing.T) {
	suite.Run(t, new(testBaseNode))
}
