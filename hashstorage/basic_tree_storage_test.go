package hashstorage

import (
	"fmt"
	"testing"

	"github.com/spikeekips/avl"
	"github.com/spikeekips/avl/hashable"
	"github.com/stretchr/testify/suite"
	"golang.org/x/xerrors"
)

type testBasicTreeStorage struct {
	suite.Suite
}

func newNode(i int) *hashable.BaseHashableNode {
	return hashable.NewBaseHashableNode(
		avl.NewBaseNode([]byte(fmt.Sprintf("%03d", i))),
	)
}

func (t *testBasicTreeStorage) TestSimpleSave() {
	ls, err := NewLevelDBStorage("memory://")
	_ = ls.SetLogger(log)
	t.NoError(err)
	nodePool := NewStorageNodePool(ls)

	tr := NewTree(nodePool)
	_ = tr.SetLogger(log)
	_, _ = tr.Add(newNode(100))

	treeName := []byte("new ts")
	ts, err := NewBasicTreeStorage(treeName, tr, ls)
	t.NoError(err)
	_ = ts.SetLogger(log)

	err = ts.Save()
	t.NoError(err)

	rootHash, err := ls.GetRoot(treeName)
	t.NoError(err)
	t.NotEmpty(rootHash)
}

func (t *testBasicTreeStorage) TestSimpleLoad() {
	treeName := []byte("new ts")

	ls, err := NewLevelDBStorage("memory://")
	_ = ls.SetLogger(log)
	t.NoError(err)

	nodePool := NewStorageNodePool(ls)
	_ = nodePool.SetLogger(log)
	t.NoError(err)

	{ // save
		tr := NewTree(nodePool)
		_ = tr.SetLogger(log)
		_, _ = tr.Add(newNode(100))

		ts, err := NewBasicTreeStorage(treeName, tr, ls)
		t.NoError(err)
		_ = ts.SetLogger(log)

		err = ts.Save()
		t.NoError(err)
	}

	{ // load
		tr, err := LoadTreeFromBasicTreeStorage(treeName, ls, nodePool)
		t.NoError(err)
		t.NotNil(tr)

		err = tr.IsValid()
		t.NoError(err)
	}
}

func (t *testBasicTreeStorage) TestLoad100() {
	treeName := []byte("new ts")

	ls, err := NewLevelDBStorage("memory://")
	_ = ls.SetLogger(log)
	t.NoError(err)

	var nodes []avl.Node
	nodePool := NewStorageNodePool(ls)
	_ = nodePool.SetLogger(log)

	{ // save
		tr := NewTree(nodePool)
		_ = tr.SetLogger(log)
		for i := 0; i < 100; i++ {
			node := newNode(i)
			_, _ = tr.Add(node)
			nodes = append(nodes, node)
		}

		ts, err := NewBasicTreeStorage(treeName, tr, ls)
		t.NoError(err)
		_ = ts.SetLogger(log)

		err = ts.Save()
		t.NoError(err)
	}

	{ // load
		newNodePool := NewStorageNodePool(ls)
		_ = newNodePool.SetLogger(log)

		tr, err := LoadTreeFromBasicTreeStorage(treeName, ls, newNodePool)
		t.NoError(err)
		t.NotNil(tr)

		err = tr.IsValid()
		t.NoError(err)

		for _, n := range nodes {
			node, err := tr.Get(n.Key())
			t.NoError(err)
			t.NotNil(node)
		}
	}
}

func testBasicTreeStorageSaveN(n int, name []byte, storage Storage) error {
	var nodes []hashable.HashableNode
	nodePool := NewStorageNodePool(storage)
	_ = nodePool.SetLogger(log)

	for i := 0; i < n; i++ {
		nodes = append(nodes, newNode(i))
	}

	tr := NewTree(nodePool)
	_ = tr.SetLogger(log)

	_, _ = tr.Adds(nodes)

	ts, err := NewBasicTreeStorage(name, tr, storage)
	if err != nil {
		return err
	}
	_ = ts.SetLogger(log)

	if err := ts.Save(); err != nil {
		return err
	}

	return nil
}

func testBasicTreeStorageLoadN(name []byte, storage Storage) error {
	nodePool := NewStorageNodePool(storage)
	_ = nodePool.SetLogger(log)

	tr, err := LoadTreeFromBasicTreeStorage(name, storage, nodePool)
	if err != nil {
		return err
	} else if tr == nil {
		return xerrors.Errorf("empty tree")
	}

	if err := tr.IsValid(); err != nil {
		return err
	}

	return nil
}

func benchmarkBasicTreeStorageSaveN(n int, b *testing.B) {
	treeName := []byte("find me")
	for i := 0; i < b.N; i++ {
		ls, err := NewLevelDBStorage("memory://")
		if err != nil {
			continue
		}
		_ = ls.SetLogger(log)

		if err := testBasicTreeStorageSaveN(n, treeName, ls); err != nil {
			panic(err)
		}
	}
}

func benchmarkBasicTreeStorageSaveLoadN(n int, b *testing.B) {
	treeName := []byte("find me")
	for i := 0; i < b.N; i++ {
		ls, err := NewLevelDBStorage("memory://")
		if err != nil {
			continue
		}
		_ = ls.SetLogger(log)

		if err := testBasicTreeStorageSaveN(n, treeName, ls); err != nil {
			panic(err)
		}
		if err := testBasicTreeStorageLoadN(treeName, ls); err != nil {
			panic(err)
		}
	}
}

func BenchmarkBasicTreeStorageSave10(b *testing.B)   { benchmarkBasicTreeStorageSaveN(10, b) }
func BenchmarkBasicTreeStorageSave100(b *testing.B)  { benchmarkBasicTreeStorageSaveN(100, b) }
func BenchmarkBasicTreeStorageSave500(b *testing.B)  { benchmarkBasicTreeStorageSaveN(500, b) }
func BenchmarkBasicTreeStorageSave1000(b *testing.B) { benchmarkBasicTreeStorageSaveN(1000, b) }
func BenchmarkBasicTreeStorageSave1500(b *testing.B) { benchmarkBasicTreeStorageSaveN(1500, b) }
func BenchmarkBasicTreeStorageSave2000(b *testing.B) { benchmarkBasicTreeStorageSaveN(2000, b) }
func BenchmarkBasicTreeStorageSave2500(b *testing.B) { benchmarkBasicTreeStorageSaveN(2500, b) }
func BenchmarkBasicTreeStorageSave3000(b *testing.B) { benchmarkBasicTreeStorageSaveN(3000, b) }
func BenchmarkBasicTreeStorageSave3500(b *testing.B) { benchmarkBasicTreeStorageSaveN(3500, b) }
func BenchmarkBasicTreeStorageSave4000(b *testing.B) { benchmarkBasicTreeStorageSaveN(4000, b) }
func BenchmarkBasicTreeStorageSave4500(b *testing.B) { benchmarkBasicTreeStorageSaveN(4500, b) }
func BenchmarkBasicTreeStorageSave5000(b *testing.B) { benchmarkBasicTreeStorageSaveN(5000, b) }

func BenchmarkBasicTreeStorageSaveLoad10(b *testing.B)   { benchmarkBasicTreeStorageSaveLoadN(10, b) }
func BenchmarkBasicTreeStorageSaveLoad100(b *testing.B)  { benchmarkBasicTreeStorageSaveLoadN(100, b) }
func BenchmarkBasicTreeStorageSaveLoad500(b *testing.B)  { benchmarkBasicTreeStorageSaveLoadN(500, b) }
func BenchmarkBasicTreeStorageSaveLoad1000(b *testing.B) { benchmarkBasicTreeStorageSaveLoadN(1000, b) }
func BenchmarkBasicTreeStorageSaveLoad1500(b *testing.B) { benchmarkBasicTreeStorageSaveLoadN(1500, b) }
func BenchmarkBasicTreeStorageSaveLoad2000(b *testing.B) { benchmarkBasicTreeStorageSaveLoadN(2000, b) }
func BenchmarkBasicTreeStorageSaveLoad2500(b *testing.B) { benchmarkBasicTreeStorageSaveLoadN(2500, b) }
func BenchmarkBasicTreeStorageSaveLoad3000(b *testing.B) { benchmarkBasicTreeStorageSaveLoadN(3000, b) }
func BenchmarkBasicTreeStorageSaveLoad3500(b *testing.B) { benchmarkBasicTreeStorageSaveLoadN(3500, b) }
func BenchmarkBasicTreeStorageSaveLoad4000(b *testing.B) { benchmarkBasicTreeStorageSaveLoadN(4000, b) }
func BenchmarkBasicTreeStorageSaveLoad4500(b *testing.B) { benchmarkBasicTreeStorageSaveLoadN(4500, b) }
func BenchmarkBasicTreeStorageSaveLoad5000(b *testing.B) { benchmarkBasicTreeStorageSaveLoadN(5000, b) }

func TestBasicTreeStorage(t *testing.T) {
	suite.Run(t, new(testBasicTreeStorage))
}
