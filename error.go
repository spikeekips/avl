package avl

import (
	"fmt"

	"golang.org/x/xerrors"
)

type WrapError struct {
	S     string
	Err   error
	Frame xerrors.Frame
}

func NewWrapError(s string, a ...interface{}) WrapError {
	return WrapError{S: fmt.Sprintf(s, a...)}
}

func (we WrapError) Wrap(err error) error {
	return WrapError{
		S:     we.S,
		Err:   err,
		Frame: xerrors.Caller(1),
	}
}

func (we WrapError) Wrapf(s string, a ...interface{}) error {
	return WrapError{
		S:     we.S,
		Err:   xerrors.Errorf(s, a...),
		Frame: xerrors.Caller(1),
	}
}

func (we WrapError) Is(err error) bool {
	if err == nil {
		return false
	}

	e, ok := err.(WrapError)
	if !ok {
		return false
	}

	return e.S == we.S
}

func (we WrapError) Unwrap() error {
	return we.Err
}

func (we WrapError) FormatError(p xerrors.Printer) error {
	we.Frame.Format(p)
	return we.Unwrap()
}

func (we WrapError) Error() string {
	if we.Err == nil {
		return we.S
	}

	return fmt.Sprintf("%s; %v", we.S, we.Err)
}
