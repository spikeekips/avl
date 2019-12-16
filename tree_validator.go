package avl

import "github.com/rs/zerolog"

var (
	InvalidTreeError = NewWrapError("invalid tree")
)

// TreeValidator will validate Tree is formed properly.
type TreeValidator struct {
	*Logger
	tr *Tree
}

// NewTreeValidator returns new TreeValidator.
func NewTreeValidator(tr *Tree) TreeValidator {
	return TreeValidator{
		Logger: NewLogger(func(c zerolog.Context) zerolog.Context {
			return c.Str("module", "avl_tree_validator")
		}),
		tr: tr,
	}
}

// IsValid checks whether tree is valid or not.
func (tv TreeValidator) IsValid() error {
	if tv.tr.Root() == nil {
		return nil
	}

	if err := tv.validate(tv.tr.Root(), nil); err != nil {
		return err
	}

	// check orphans
	if found, err := tv.hasOrphans(); err != nil {
		return err
	} else if found {
		return InvalidTreeError.Wrapf("orphan(s) found")
	}

	return nil
}

func (tv TreeValidator) hasOrphans() (bool, error) {
	var found bool
	err := tv.tr.NodePool().Traverse(func(node Node) (bool, error) {
		n, err := tv.tr.Get(node.Key())
		if err != nil {
			return false, err
		} else if n == nil {
			log.Debug().Bytes("key", node.Key()).Msg("orphan found")
			found = true
			return false, nil
		}

		return true, nil
	})

	return found, err
}

func (tv TreeValidator) validate(node Node, parents []Node) error {
	log_ := tv.Log().With().Int("parents", len(parents)).Bytes("key", node.Key()).Logger()

	np := tv.tr.NodePool()

	var left, right Node
	var err error
	if left, err = np.Get(node.LeftKey()); err != nil {
		return err
	}
	if right, err = np.Get(node.RightKey()); err != nil {
		return err
	}

	if err := IsValidNode(node, left, right); err != nil {
		log_.Error().Err(err).Msg("invalid node found")
		return err
	}

	if left != nil {
		if err := tv.validate(left, append(parents, node)); err != nil {
			return err
		}
	}
	if right != nil {
		if err := tv.validate(right, append(parents, node)); err != nil {
			return err
		}
	}

	log_.Debug().Msg("validated")

	return nil
}
