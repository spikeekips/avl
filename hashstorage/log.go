package hashstorage

import (
	"github.com/rs/zerolog"

	"github.com/spikeekips/avl"
)

var (
	log zerolog.Logger = avl.SetDefaultLog()
)
