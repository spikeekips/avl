package hashstorage

type Storage interface {
	GetRoot(name []byte) ([]byte /* root key */, error)
	SetRoot(name []byte, rootKey []byte) error
	GetNode(key []byte) ([]byte, error)
	SetNode(key []byte, value []byte) error
	Batch() BatchStorage
	GetRaw(key []byte) ([]byte, error)
	SetRaw(key []byte, value []byte) error
}

type BatchStorage interface {
	SetRoot(name []byte, rootKey []byte) error
	SetNode(key []byte, value []byte) error
	Commit() error
}
