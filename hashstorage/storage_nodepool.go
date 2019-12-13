package hashstorage

import (
	"github.com/rs/zerolog"
	"github.com/spikeekips/avl"
	"github.com/spikeekips/avl/hashable"
)

type StorageNodePool struct {
	*avl.Logger
	*avl.MapNodePool
	storage Storage
}

func NewStorageNodePool(storage Storage) *StorageNodePool {
	return &StorageNodePool{
		Logger: avl.NewLogger(func(c zerolog.Context) zerolog.Context {
			return c.Str("module", "avl_storage_nodepool")
		}),
		MapNodePool: avl.NewMapNodePool(),
		storage:     storage,
	}
}

func (sn *StorageNodePool) Get(key []byte) (avl.Node, error) {
	log_ := sn.Log().With().Bytes("key", key).Logger()

	if node, err := sn.MapNodePool.Get(key); err != nil {
		log_.Error().Err(err).Msg("failed to Get")
		return nil, err
	} else if node != nil {
		hn, ok := node.(*hashable.BaseHashableNode)
		if !ok {
			err := avl.InvalidNodeError.Wrapf("not HashableNode type: %T", node)
			log_.Error().Err(err).Send()
			return nil, err
		}
		log_.Debug().Msg("Get from map")
		return hn, nil
	}

	v, err := sn.storage.GetNode(key)
	if err != nil {
		err := avl.StorageInternalError.Wrap(err)
		log_.Error().Err(err).Send()
		return nil, err
	}

	var bh hashable.BaseHashableNode
	if err := bh.UnmarshalBinary(v); err != nil {
		log_.Error().Err(err).Send()
		return nil, err
	}

	log_.Debug().Msg("Get from storage")
	return &bh, nil
}

func (sn *StorageNodePool) Set(node avl.Node) error {
	log_ := sn.Log().With().Bytes("key", node.Key()).Logger()

	err := sn.MapNodePool.Set(node)
	if err != nil {
		log_.Error().Err(err).Msg("failed to Set")
		return err
	}

	log_.Debug().Msg("Set")
	return nil
}

func (sn *StorageNodePool) UnmarshalNode(data []byte) (*hashable.BaseHashableNode, error) {
	var bh hashable.BaseHashableNode
	if err := bh.UnmarshalBinary(data); err != nil {
		return nil, err
	}

	return &bh, nil
}
