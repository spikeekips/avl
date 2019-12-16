package avl

import (
	"os"

	"github.com/rs/zerolog"
)

var (
	nilLog zerolog.Logger = zerolog.Nop()
	log    zerolog.Logger = SetDefaultLog()
)

// SetDefaultLog returns the predefined(default) zerolog.Logger.
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

// Logger is basic log for this package.
type Logger struct {
	root        zerolog.Logger
	l           *zerolog.Logger
	contextFunc []func(zerolog.Context) zerolog.Context
}

// NewLogger returns new Logger. With argument function you can pass the context
// to Logger. For example,
//
//	logger := NewLogger(func(c zerolog.Context) zerolog.Context {
//		return c.Str("module", "avl_tree_validator")
//	})
func NewLogger(cf func(zerolog.Context) zerolog.Context) *Logger {
	zl := &Logger{}
	if cf != nil {
		zl.contextFunc = append(zl.contextFunc, cf)
	}

	return zl
}

// SetLogger set the new zerolog.Logger.
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

func (zl *Logger) Log() *zerolog.Logger {
	if zl.l == nil {
		return &nilLog
	}

	return zl.l
}
