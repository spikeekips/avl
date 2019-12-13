package hashstorage

import (
	"testing"

	"github.com/spikeekips/avl"
	"github.com/spikeekips/avl/hashable"
	"github.com/stretchr/testify/suite"
	"golang.org/x/xerrors"
)

type testDumpTreeStorage struct {
	suite.Suite
}

func (t *testDumpTreeStorage) TestSimpleSave() {
	ls, err := NewLevelDBStorage("memory://")
	_ = ls.SetLogger(log)
	t.NoError(err)
	nodePool := NewStorageNodePool(ls)

	tr := NewTree(nodePool)
	_ = tr.SetLogger(log)
	_, _ = tr.Add(newNode(100))

	treeName := []byte("new ts")
	ts, err := NewDumpTreeStorage(treeName, tr, ls)
	t.NoError(err)
	_ = ts.SetLogger(log)

	err = ts.Save()
	t.NoError(err)

	dump, err := ls.GetRaw(treeName)
	t.NoError(err)
	t.NotEmpty(dump)
}

func (t *testDumpTreeStorage) TestSimpleLoad() {
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

		ts, err := NewDumpTreeStorage(treeName, tr, ls)
		t.NoError(err)
		_ = ts.SetLogger(log)

		err = ts.Save()
		t.NoError(err)
	}

	{ // load
		tr, err := LoadTreeFromDumpTreeStorage(treeName, ls, nodePool)
		t.NoError(err)
		t.NotNil(tr)

		err = tr.IsValid()
		t.NoError(err)
	}
}

func (t *testDumpTreeStorage) TestLoad100() {
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

		ts, err := NewDumpTreeStorage(treeName, tr, ls)
		t.NoError(err)
		_ = ts.SetLogger(log)

		err = ts.Save()
		t.NoError(err)
	}

	{ // load
		newNodePool := NewStorageNodePool(ls)
		_ = newNodePool.SetLogger(log)

		tr, err := LoadTreeFromDumpTreeStorage(treeName, ls, newNodePool)
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

func testDumpTreeStorageSaveN(n int, name []byte, storage Storage) error {
	var nodes []hashable.HashableNode
	nodePool := NewStorageNodePool(storage)
	_ = nodePool.SetLogger(log)

	for i := 0; i < n; i++ {
		nodes = append(nodes, newNode(i))
	}

	tr := NewTree(nodePool)
	_ = tr.SetLogger(log)

	_, _ = tr.Adds(nodes)

	ts, err := NewDumpTreeStorage(name, tr, storage)
	if err != nil {
		return err
	}
	_ = ts.SetLogger(log)

	if err := ts.Save(); err != nil {
		return err
	}

	return nil
}

func testDumpTreeStorageLoadN(name []byte, storage Storage) error {
	nodePool := NewStorageNodePool(storage)
	_ = nodePool.SetLogger(log)

	tr, err := LoadTreeFromDumpTreeStorage(name, storage, nodePool)
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

func benchmarkDumpTreeStorageSaveN(n int, b *testing.B) {
	treeName := []byte("find me")
	for i := 0; i < b.N; i++ {
		ls, err := NewLevelDBStorage("memory://")
		if err != nil {
			continue
		}
		_ = ls.SetLogger(log)

		if err := testDumpTreeStorageSaveN(n, treeName, ls); err != nil {
			panic(err)
		}
	}
}

func benchmarkDumpTreeStorageSaveLoadN(n int, b *testing.B) {
	treeName := []byte("find me")
	for i := 0; i < b.N; i++ {
		ls, err := NewLevelDBStorage("memory://")
		if err != nil {
			continue
		}
		_ = ls.SetLogger(log)

		if err := testDumpTreeStorageSaveN(n, treeName, ls); err != nil {
			panic(err)
		}
		if err := testDumpTreeStorageLoadN(treeName, ls); err != nil {
			panic(err)
		}
	}
}

func BenchmarkDumpTreeStorageSave10(b *testing.B)   { benchmarkDumpTreeStorageSaveN(10, b) }
func BenchmarkDumpTreeStorageSave100(b *testing.B)  { benchmarkDumpTreeStorageSaveN(100, b) }
func BenchmarkDumpTreeStorageSave500(b *testing.B)  { benchmarkDumpTreeStorageSaveN(500, b) }
func BenchmarkDumpTreeStorageSave1000(b *testing.B) { benchmarkDumpTreeStorageSaveN(1000, b) }
func BenchmarkDumpTreeStorageSave1500(b *testing.B) { benchmarkDumpTreeStorageSaveN(1500, b) }
func BenchmarkDumpTreeStorageSave2000(b *testing.B) { benchmarkDumpTreeStorageSaveN(2000, b) }
func BenchmarkDumpTreeStorageSave2500(b *testing.B) { benchmarkDumpTreeStorageSaveN(2500, b) }
func BenchmarkDumpTreeStorageSave3000(b *testing.B) { benchmarkDumpTreeStorageSaveN(3000, b) }
func BenchmarkDumpTreeStorageSave3500(b *testing.B) { benchmarkDumpTreeStorageSaveN(3500, b) }
func BenchmarkDumpTreeStorageSave4000(b *testing.B) { benchmarkDumpTreeStorageSaveN(4000, b) }
func BenchmarkDumpTreeStorageSave4500(b *testing.B) { benchmarkDumpTreeStorageSaveN(4500, b) }
func BenchmarkDumpTreeStorageSave5000(b *testing.B) { benchmarkDumpTreeStorageSaveN(5000, b) }

func BenchmarkDumpTreeStorageSaveLoad10(b *testing.B)   { benchmarkDumpTreeStorageSaveLoadN(10, b) }
func BenchmarkDumpTreeStorageSaveLoad100(b *testing.B)  { benchmarkDumpTreeStorageSaveLoadN(100, b) }
func BenchmarkDumpTreeStorageSaveLoad500(b *testing.B)  { benchmarkDumpTreeStorageSaveLoadN(500, b) }
func BenchmarkDumpTreeStorageSaveLoad1000(b *testing.B) { benchmarkDumpTreeStorageSaveLoadN(1000, b) }
func BenchmarkDumpTreeStorageSaveLoad1500(b *testing.B) { benchmarkDumpTreeStorageSaveLoadN(1500, b) }
func BenchmarkDumpTreeStorageSaveLoad2000(b *testing.B) { benchmarkDumpTreeStorageSaveLoadN(2000, b) }
func BenchmarkDumpTreeStorageSaveLoad2500(b *testing.B) { benchmarkDumpTreeStorageSaveLoadN(2500, b) }
func BenchmarkDumpTreeStorageSaveLoad3000(b *testing.B) { benchmarkDumpTreeStorageSaveLoadN(3000, b) }
func BenchmarkDumpTreeStorageSaveLoad3500(b *testing.B) { benchmarkDumpTreeStorageSaveLoadN(3500, b) }
func BenchmarkDumpTreeStorageSaveLoad4000(b *testing.B) { benchmarkDumpTreeStorageSaveLoadN(4000, b) }
func BenchmarkDumpTreeStorageSaveLoad4500(b *testing.B) { benchmarkDumpTreeStorageSaveLoadN(4500, b) }
func BenchmarkDumpTreeStorageSaveLoad5000(b *testing.B) { benchmarkDumpTreeStorageSaveLoadN(5000, b) }

func TestDumpTreeStorage(t *testing.T) {
	suite.Run(t, new(testDumpTreeStorage))
}
