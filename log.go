package avl

import (
	"os"

	"github.com/rs/zerolog"
)

var (
	nilLog zerolog.Logger = zerolog.Nop()
	log    zerolog.Logger = SetDefaultLog()
)

func SetDefaultLog() zerolog.Logger {
	if os.Getenv("AVL_DEBUG") != "1" {
		return zerolog.Nop()
	}

	return zerolog.
		New(os.Stderr).
		With().
		Timestamp().
		Caller().
		Stack().
		Logger().
		Level(zerolog.DebugLevel)
}

type Logger struct {
	root        zerolog.Logger
	l           *zerolog.Logger
	contextFunc []func(zerolog.Context) zerolog.Context
}

func NewLogger(cf func(zerolog.Context) zerolog.Context) *Logger {
	zl := &Logger{}
	if cf != nil {
		zl.contextFunc = append(zl.contextFunc, cf)
	}

	return zl
}

func (zl *Logger) SetLogger(l zerolog.Logger) *Logger {
	zl.root = l
	if len(zl.contextFunc) > 0 {
		for _, cf := range zl.contextFunc {
			l = cf(l.With()).Logger()
		}
		zl.l = &l
	} else {
		zl.l = &l
	}

	return zl
}

func (zl *Logger) RootLog() zerolog.Logger {
	return zl.root
}

func (zl *Logger) Log() *zerolog.Logger {
	if zl.l == nil {
		return &nilLog
	}

	return zl.l
}
