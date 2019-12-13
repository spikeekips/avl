package avl

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.org/x/xerrors"
)

type testTree struct {
	suite.Suite
}

func newNode(i int) *BaseNode {
	return NewBaseNode(
		[]byte(fmt.Sprintf("%03d", i)),
	)
}

func (t *testTree) newTree() *Tree {
	tr := NewTree(NewMapNodePool())
	_ = tr.SetLogger(log)

	return tr
}

func (t *testTree) parseNodeKey(b []byte) int {
	i, _ := strconv.ParseInt(string(b), 10, 64)
	return int(i)
}

func (t *testTree) TestNew() {
	tr := t.newTree()

	n := newNode(10)
	_, err := tr.Add(n)
	t.NoError(err)
}

func (t *testTree) TestInvalidNode() {
	tr := t.newTree()

	n := newNode(10)
	n.key = nil
	_, err := tr.Add(n)
	t.NotNil(err)
	t.True(xerrors.Is(err, InvalidNodeError))
}

func (t *testTree) TestAddLeft() {
	tr := t.newTree()

	n := newNode(10)
	_, err := tr.Add(n)
	t.NoError(err)

	l0 := newNode(5)
	_, err = tr.Add(l0)
	t.NoError(err)

	l1 := newNode(3)
	_, err = tr.Add(l1)
	t.NoError(err)
}

func (t *testTree) TestGet() {
	tr := t.newTree()

	keys := []int{100, 50, 150, 20, 70, 120, 170}

	for _, k := range keys {
		_, _ = tr.Add(newNode(k))
	}

	{
		n, err := tr.Get([]byte("unknown"))
		t.Nil(n)
		t.NoError(err)
	}

	for _, k := range keys {
		key := []byte(fmt.Sprintf("%03d", k))
		n, err := tr.Get(key)
		t.NoError(err)
		t.Equal(n.Key(), key)
	}
}

func (t *testTree) TestGetWithParents() {
	tr := t.newTree()

	keys := []int{100, 50, 150, 20, 70, 120, 170}

	for _, k := range keys {
		_, _ = tr.Add(newNode(k))
	}

	{
		n, parents, err := tr.GetWithParents([]byte("unknown"))
		t.Nil(n)
		t.Nil(parents)
		t.NoError(err)
	}

	shapes := map[int][]int{
		100: nil,
		50:  {100},
		150: {100},
		20:  {100, 50},
		70:  {100, 50},
		120: {100, 150},
		170: {100, 150},
	}

	for k, expectedParents := range shapes {
		key := []byte(fmt.Sprintf("%03d", k))
		n, parents, err := tr.GetWithParents(key)
		t.NoError(err)

		t.Equal(k, t.parseNodeKey(n.Key()))

		var ps []int
		for _, p := range parents {
			ps = append(ps, t.parseNodeKey(p.Key()))
		}
		t.Equal(expectedParents, ps)
	}
}

func (t *testTree) TestTraverse() {
	tr := t.newTree()

	keys := []int{100, 50, 150, 20, 70, 120, 170}

	for _, k := range keys {
		_, _ = tr.Add(newNode(k))
	}

	{ // tranverse all
		traversed := map[int]Node{}
		err := tr.Traverse(func(node Node) (bool, error) {
			traversed[t.parseNodeKey(node.Key())] = node
			return true, nil
		})
		t.NoError(err)

		t.Equal(len(keys), len(traversed))

		for _, k := range keys {
			_, found := traversed[k]
			t.True(found)
		}
	}

	{ // filter; lesser than 100; 100, 50, 20, 70
		traversed := map[int]Node{}
		err := tr.Traverse(func(node Node) (bool, error) {
			n := t.parseNodeKey(node.Key())
			if n > 100 {
				return true, nil
			}

			traversed[n] = node
			return true, nil
		})
		t.NoError(err)

		t.Equal(4, len(traversed))

		for _, k := range []int{100, 50, 20, 70} {
			_, found := traversed[k]
			t.True(found)
		}
	}

	{ // compare by key; lesser than 100; 100, 50, 20, 70
		traversed := map[int]Node{}

		point := []byte(fmt.Sprintf("%03d", 100))
		err := tr.Traverse(func(node Node) (bool, error) {
			// prevent to traverse the right branch
			if bytes.Compare(node.Key(), point) > 0 {
				return false, nil
			}

			traversed[t.parseNodeKey(node.Key())] = node
			return true, nil
		})
		t.NoError(err)

		t.Equal(4, len(traversed))

		for _, k := range []int{100, 50, 20, 70} {
			_, found := traversed[k]
			t.True(found)
		}
	}
}

func (t *testTree) TestValidate() {
	{ // ok
		tr := t.newTree()
		keys := []int{100, 50, 150, 20, 70, 120, 170}
		for _, k := range keys {
			_, _ = tr.Add(newNode(k))
		}

		err := tr.IsValid()
		t.NoError(err)
	}

	{ // 170 has 100 as left
		tr := t.newTree()

		nodes := map[int]*BaseNode{}
		keys := []int{100, 50, 150, 20, 70, 120, 170}
		for _, k := range keys {
			n := newNode(k)
			_, _ = tr.Add(n)
			nodes[k] = n
		}
		_, _ = tr.Add(nodes[170])
		err := nodes[170].SetLeftKey([]byte(fmt.Sprintf("%03d", 100)))
		t.NoError(err)

		_, _ = tr.resetNodeHeight(nodes[170], false)
		_, _ = tr.resetNodeHeight(nodes[150], false)
		_, _ = tr.resetNodeHeight(nodes[100], false)

		err = tr.IsValid()
		t.True(xerrors.Is(err, InvalidNodeError))
	}
}

func TestTree(t *testing.T) {
	suite.Run(t, new(testTree))
}

func benchmarkTreeN(n int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		tr := NewTree(NewMapNodePool())
		for j := 1; j <= n; j++ {
			if _, err := tr.Add(newNode(j)); err != nil {
				panic(err)
			}
		}
	}
}

func BenchmarkTree10(b *testing.B)    { benchmarkTreeN(10, b) }
func BenchmarkTree100(b *testing.B)   { benchmarkTreeN(100, b) }
func BenchmarkTree200(b *testing.B)   { benchmarkTreeN(200, b) }
func BenchmarkTree300(b *testing.B)   { benchmarkTreeN(300, b) }
func BenchmarkTree400(b *testing.B)   { benchmarkTreeN(400, b) }
func BenchmarkTree500(b *testing.B)   { benchmarkTreeN(500, b) }
func BenchmarkTree600(b *testing.B)   { benchmarkTreeN(600, b) }
func BenchmarkTree700(b *testing.B)   { benchmarkTreeN(700, b) }
func BenchmarkTree800(b *testing.B)   { benchmarkTreeN(800, b) }
func BenchmarkTree900(b *testing.B)   { benchmarkTreeN(900, b) }
func BenchmarkTree1000(b *testing.B)  { benchmarkTreeN(1000, b) }
func BenchmarkTree1100(b *testing.B)  { benchmarkTreeN(1100, b) }
func BenchmarkTree1200(b *testing.B)  { benchmarkTreeN(1200, b) }
func BenchmarkTree1300(b *testing.B)  { benchmarkTreeN(1300, b) }
func BenchmarkTree1400(b *testing.B)  { benchmarkTreeN(1400, b) }
func BenchmarkTree1500(b *testing.B)  { benchmarkTreeN(1500, b) }
func BenchmarkTree1600(b *testing.B)  { benchmarkTreeN(1600, b) }
func BenchmarkTree1700(b *testing.B)  { benchmarkTreeN(1700, b) }
func BenchmarkTree1800(b *testing.B)  { benchmarkTreeN(1800, b) }
func BenchmarkTree1900(b *testing.B)  { benchmarkTreeN(1900, b) }
func BenchmarkTree10000(b *testing.B) { benchmarkTreeN(10000, b) }
