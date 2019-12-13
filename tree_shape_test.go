package avl

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.org/x/xerrors"
)

type testTreeShape struct {
	suite.Suite
	printCount int32 // nolint
}

func (t *testTreeShape) TestShapes() {
	cases := []struct {
		name    string
		keys    []int
		shape   map[int]shape
		updated []int
	}{
		{
			name: "add #0",
			keys: []int{100},
			shape: map[int]shape{
				100: {height: 0},
			},
			updated: []int{},
		},
		{
			name: "add #1",
			keys: []int{100, 50},
			shape: map[int]shape{
				100: {height: 1, left: 50},
				50:  {height: 0},
			},
			updated: []int{100},
		},
		{
			name: "add #2",
			keys: []int{100, 50, 150},
			shape: map[int]shape{
				100: {height: 1, left: 50, right: 150},
				50:  {height: 0},
				150: {height: 0},
			},
			updated: []int{100},
		},
		{
			name: "add #3",
			keys: []int{100, 50, 150, 30},
			shape: map[int]shape{
				100: {height: 2, left: 50, right: 150},
				50:  {height: 1, left: 30},
				150: {height: 0},
				30:  {height: 0},
			},
			updated: []int{100, 50},
		},
		{
			name: "single rotation left #0",
			keys: []int{100, 50, 30},
			shape: map[int]shape{
				100: {height: 0},
				50:  {height: 1, left: 30, right: 100},
				30:  {height: 0},
			},
			updated: []int{100, 50},
		},
		{
			name: "single rotation left #1",
			keys: []int{100, 50, 80},
			shape: map[int]shape{
				100: {height: 0},
				50:  {height: 0},
				80:  {height: 1, left: 50, right: 100},
			},
			updated: []int{100, 50},
		},
		{
			name: "single rotation right #0",
			keys: []int{100, 150, 180},
			shape: map[int]shape{
				100: {height: 0},
				150: {height: 1, left: 100, right: 180},
				180: {height: 0},
			},
			updated: []int{100, 150},
		},
		{
			name: "single rotation right #1",
			keys: []int{100, 150, 110},
			shape: map[int]shape{
				100: {height: 0},
				150: {height: 0},
				110: {height: 1, left: 100, right: 150},
			},
			updated: []int{100, 150},
		},
		{
			name: "single rotation left #2",
			keys: []int{100, 50, 150, 30, 10},
			shape: map[int]shape{
				100: {height: 2, left: 30, right: 150},
				50:  {height: 0},
				150: {height: 0},
				30:  {height: 1, left: 10, right: 50},
				10:  {height: 0},
			},
			updated: []int{100, 50, 30},
		},
		{
			name: "single rotation right #2",
			keys: []int{100, 50, 150, 180, 200},
			shape: map[int]shape{
				100: {height: 2, left: 50, right: 180},
				50:  {height: 0},
				150: {height: 0},
				180: {height: 1, left: 150, right: 200},
				200: {height: 0},
			},
			updated: []int{100, 150, 180},
		},
		{
			name: "single rotation deep left #0",
			keys: []int{100, 50, 150, 30, 70, 130, 10, 5},
			shape: map[int]shape{
				100: {height: 3, left: 50, right: 150},
				50:  {height: 2, left: 10, right: 70},
				150: {height: 1, left: 130},
				30:  {height: 0},
				70:  {height: 0},
				130: {height: 0},
				10:  {height: 1, left: 5, right: 30},
				5:   {height: 0},
			},
			updated: []int{100, 50, 30, 10},
		},
		{
			name: "single rotation deep left #1",
			keys: []int{100, 50, 150, 30, 70, 130, 10, 20},
			shape: map[int]shape{
				100: {height: 3, left: 50, right: 150},
				50:  {height: 2, left: 20, right: 70},
				150: {height: 1, left: 130},
				30:  {height: 0},
				70:  {height: 0},
				130: {height: 0},
				10:  {height: 0},
				20:  {height: 1, left: 10, right: 30},
			},
			updated: []int{100, 50, 30, 10},
		},
		{
			name: "single rotation deep right #0",
			keys: []int{100, 50, 150, 30, 70, 130, 170, 180, 200},
			shape: map[int]shape{
				100: {height: 3, left: 50, right: 150},
				50:  {height: 1, left: 30, right: 70},
				150: {height: 2, left: 130, right: 180},
				30:  {height: 0},
				70:  {height: 0},
				130: {height: 0},
				170: {height: 0},
				180: {height: 1, left: 170, right: 200},
				200: {height: 0},
			},
			updated: []int{100, 150, 170, 180},
		},
		{
			name: "single rotation deep right #1",
			keys: []int{100, 50, 150, 30, 70, 130, 170, 180, 175},
			shape: map[int]shape{
				100: {height: 3, left: 50, right: 150},
				50:  {height: 1, left: 30, right: 70},
				150: {height: 2, left: 130, right: 175},
				30:  {height: 0},
				70:  {height: 0},
				130: {height: 0},
				170: {height: 0},
				180: {height: 0},
				175: {height: 1, left: 170, right: 180},
			},
			updated: []int{100, 150, 170, 180},
		},
		{
			name: "left-left #0",
			keys: []int{100, 50, 150, 30, 70, 10},
			shape: map[int]shape{
				100: {height: 1, left: 70, right: 150},
				50:  {height: 2, left: 30, right: 100},
				150: {height: 0},
				30:  {height: 1, left: 10},
				70:  {height: 0},
				10:  {height: 0},
			},
			updated: []int{100, 50, 30},
		},
		{
			name: "right-right #0",
			keys: []int{100, 50, 150, 130, 180, 200},
			shape: map[int]shape{
				100: {height: 1, left: 50, right: 130},
				50:  {height: 0},
				150: {height: 2, left: 100, right: 180},
				130: {height: 0},
				180: {height: 1, right: 200},
				200: {height: 0},
			},
			updated: []int{100, 150, 180},
		},
		{
			name: "left-right #0",
			keys: []int{100, 50, 150, 30, 70, 80},
			shape: map[int]shape{
				100: {height: 1, left: 80, right: 150},
				50:  {height: 1, left: 30},
				150: {height: 0},
				30:  {height: 0},
				70:  {height: 2, left: 50, right: 100},
				80:  {height: 0},
			},
			updated: []int{100, 50, 70},
		},
		{
			name: "right-left #0",
			keys: []int{100, 50, 150, 130, 180, 110},
			shape: map[int]shape{
				100: {height: 1, left: 50, right: 110},
				50:  {height: 0},
				150: {height: 1, right: 180},
				130: {height: 2, left: 100, right: 150},
				180: {height: 0},
				110: {height: 0},
			},
			updated: []int{100, 150, 130},
		},
		{
			name: "upto 12",
			keys: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			shape: map[int]shape{
				1:  {height: 0},
				2:  {height: 1, left: 1, right: 3},
				3:  {height: 0},
				4:  {height: 2, left: 2, right: 6},
				5:  {height: 0},
				6:  {height: 1, left: 5, right: 7},
				7:  {height: 0},
				8:  {height: 3, left: 4, right: 10},
				9:  {height: 0},
				10: {height: 2, left: 9, right: 11},
				11: {height: 1, right: 12},
				12: {height: 0},
			},
			updated: []int{4, 8, 10, 11},
		},
		{
			name: "upto 50",
			keys: []int{
				1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
				11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
				21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
				31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
				41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
			},
			shape: map[int]shape{
				1:  {height: 0},
				2:  {height: 1, left: 1, right: 3},
				3:  {height: 0},
				4:  {height: 2, left: 2, right: 6},
				5:  {height: 0},
				6:  {height: 1, left: 5, right: 7},
				7:  {height: 0},
				8:  {height: 3, left: 4, right: 12},
				9:  {height: 0},
				10: {height: 1, left: 9, right: 11},
				11: {height: 0},
				12: {height: 2, left: 10, right: 14},
				13: {height: 0},
				14: {height: 1, left: 13, right: 15},
				15: {height: 0},
				16: {height: 4, left: 8, right: 24},
				17: {height: 0},
				18: {height: 1, left: 17, right: 19},
				19: {height: 0},
				20: {height: 2, left: 18, right: 22},
				21: {height: 0},
				22: {height: 1, left: 21, right: 23},
				23: {height: 0},
				24: {height: 3, left: 20, right: 28},
				25: {height: 0},
				26: {height: 1, left: 25, right: 27},
				27: {height: 0},
				28: {height: 2, left: 26, right: 30},
				29: {height: 0},
				30: {height: 1, left: 29, right: 31},
				31: {height: 0},
				32: {height: 5, left: 16, right: 40},
				33: {height: 0},
				34: {height: 1, left: 33, right: 35},
				35: {height: 0},
				36: {height: 2, left: 34, right: 38},
				37: {height: 0},
				38: {height: 1, left: 37, right: 39},
				39: {height: 0},
				40: {height: 4, left: 36, right: 44},
				41: {height: 0},
				42: {height: 1, left: 41, right: 43},
				43: {height: 0},
				44: {height: 3, left: 42, right: 48},
				45: {height: 0},
				46: {height: 1, left: 45, right: 47},
				47: {height: 0},
				48: {height: 2, left: 46, right: 49},
				49: {height: 1, right: 50},
				50: {height: 0},
			},
			updated: []int{32, 40, 44, 46, 48, 49},
		},
		{
			name: "nil parent exception found",
			keys: []int{
				48, 43, 15, 25, 11, 32, 2, 44, 30, 6,
				19, 12, 49, 37, 34, 20, 14, 8, 42, 22, 18, 7, 36, 33, 40, 28,
			},
			shape: map[int]shape{
				48: {height: 1, left: 44, right: 49},
				43: {height: 2, left: 42, right: 48},
				15: {height: 4, left: 11, right: 20},
				25: {height: 5, left: 15, right: 37},
				11: {height: 3, left: 6, right: 12},
				32: {height: 2, left: 30, right: 34},
				2:  {height: 0},
				44: {height: 0},
				30: {height: 1, left: 28},
				6:  {height: 2, left: 2, right: 8},
				19: {height: 1, left: 18},
				12: {height: 1, right: 14},
				49: {height: 0},
				37: {height: 3, left: 32, right: 43},
				34: {height: 1, left: 33, right: 36},
				20: {height: 2, left: 19, right: 22},
				14: {height: 0},
				8:  {height: 1, left: 7},
				42: {height: 1, left: 40},
				22: {height: 0},
				18: {height: 0},
				7:  {height: 0},
				36: {height: 0},
				33: {height: 0},
				40: {height: 0},
				28: {height: 0},
			},
			updated: []int{25, 37, 34, 32, 30},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(
			c.name,
			func() {
				tr := NewTree(NewMapNodePool())
				_ = tr.SetLogger(log)
				nodes := map[int]*BaseNode{}

				var updatedNodes []Node
				for _, k := range c.keys {
					nodes[k] = t.newNode(k)

					var err error
					updatedNodes, err = tr.Add(nodes[k])
					t.NoError(err)
				}

				var updated []int
				for _, n := range updatedNodes {
					updated = append(updated, t.parseNodeKey(n.Key()))
				}

				for k, n := range nodes {
					var left, right int
					if n.LeftKey() != nil {
						left = t.parseNodeKey(n.LeftKey())
					}
					if n.RightKey() != nil {
						right = t.parseNodeKey(n.RightKey())
					}
					log.Debug().
						Int("key", k).
						Int("left", left).
						Int("right", right).
						Send()
				}

				t.compareShape(nodes, c.shape, "%s, %s", c.name, "before-shape")

				if c.updated != nil {
					if len(c.updated) < 1 {
						c.updated = nil
					}
					t.Equal(c.updated, updated)
				}

				for key, n := range nodes {
					// check orphans
					tn, err := tr.Get(n.Key())
					t.NoError(err)
					t.NotNil(tn, "orphan found: key=%d", key)

					left, right, err := tr.getLeafs(tn)
					t.NoError(err)
					err = IsValidNode(tn, left, right)
					t.NoError(err)
				}

				err := tr.IsValid()
				t.NoError(err)
			},
		)
	}
}

func TestTreeShape(t *testing.T) {
	suite.Run(t, new(testTreeShape))
}

type shape struct {
	height int16
	left   int
	right  int
}

func (t *testTreeShape) parseNodeKey(b []byte) int {
	i, _ := strconv.ParseInt(string(b), 10, 64)
	return int(i)
}

func (t *testTreeShape) print(tr *Tree, verbose bool) { // nolint
	defer atomic.AddInt32(&t.printCount, 1)

	var b bytes.Buffer

	PrintDotGraph(tr, &b)

	_ = ioutil.WriteFile(
		fmt.Sprintf("/tmp/%03d.dot", t.printCount),
		b.Bytes(), 0644,
	)

	if verbose {
		fmt.Fprintln(os.Stderr, b.String())
	}
}

func (t *testTreeShape) newNode(i int) *BaseNode {
	return NewBaseNode([]byte(fmt.Sprintf("%03d", i)))
}

func (t *testTreeShape) compareShape(nodes map[int]*BaseNode, shapes map[int]shape, args ...interface{}) {
	var prefix string
	if len(args) > 0 {
		prefix = fmt.Sprintf(args[0].(string)+":", args[1:]...)
	}

	for k, sh := range shapes {
		n := nodes[k]
		if n == nil {
			t.NoError(xerrors.Errorf("node not found: %v", k))
			break
		}

		t.Equal(
			sh.height,
			n.Height(),
			"%s height does not matched: %v: %v != %v",
			prefix,
			k,
			sh.height,
			n.Height(),
		)
		if sh.left < 1 {
			if n.LeftKey() != nil {
				t.NoError(
					xerrors.Errorf("left is not empty: %v: left=%v", k, t.parseNodeKey(n.LeftKey())),
					prefix,
				)
			}
		} else {
			if n.LeftKey() == nil {
				t.NoError(xerrors.Errorf("left is empty: %v", sh.left))
				continue
			}

			nLeft := t.parseNodeKey(n.LeftKey())
			t.Equal(
				sh.left, nLeft,
				"%s left does not matched: %v: %v != %v",
				prefix,
				k, sh.left, nLeft,
			)
		}

		if sh.right < 1 {
			if n.RightKey() != nil {
				t.NoError(
					xerrors.Errorf("right is not empty: %v: right=%v", k, t.parseNodeKey(n.RightKey())),
					prefix,
				)
			}
		} else {
			if n.RightKey() == nil {
				t.NoError(xerrors.Errorf("right is empty: %v", sh.right))
				continue
			}

			nRight := t.parseNodeKey(n.RightKey())
			t.Equal(
				sh.right, nRight,
				"%s right does not matched: %v: %v != %v",
				prefix,
				k, sh.right, nRight,
			)
		}
	}
}
