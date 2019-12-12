package avl

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.org/x/xerrors"
)

type testWrapError struct {
	suite.Suite
}

func (t *testWrapError) TestNew() {
	error0 := NewWrapError("0 error")
	t.True(xerrors.Is(error0, error0))

	error01 := error0.Wrapf("findme")
	t.True(xerrors.Is(error01, error0))
}

func (t *testWrapError) TestAs0() {
	error0 := NewWrapError("0 error")

	var error01 WrapError
	t.True(xerrors.As(error0, &error01))
}

func (t *testWrapError) TestAs1() {
	error0 := NewWrapError("0 error")
	error1 := error0.Wrap(os.ErrClosed)

	var error2 error
	t.True(xerrors.As(error1, &error2))
}

func (t *testWrapError) TestIs0() {
	error0 := NewWrapError("0 error")
	error1 := error0.Wrap(os.ErrClosed)

	t.True(xerrors.Is(error1, error0))
	t.True(xerrors.Is(error1, os.ErrClosed))
	t.False(xerrors.Is(error1, os.ErrNotExist))
}

func TestWrapError(t *testing.T) {
	suite.Run(t, new(testWrapError))
}
