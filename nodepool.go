package avl

var (
	StorageInternalError = NewWrapError("storage internal error occurred")
)

type NodePool interface {
	Get(key []byte) (Node, error)
	Set(node Node) error
}
