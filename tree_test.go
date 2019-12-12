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
)

type testTree struct {
	suite.Suite
}

func (t *testTree) TestShape() {
	cases := []shapeCase{
		{
			shape: map[int]shape{
				100: {height: 0},
			},
			root: 100,
		},
		{
			shape: map[int]shape{
				100: {height: 1, left: 50},
				50:  {height: 0},
			},
			root: 100,
		},
		{
			shape: map[int]shape{
				100: {height: 1, left: 50, right: 150},
				50:  {height: 0},
				150: {height: 0},
			},
			root: 100,
		},
		{
			shape: map[int]shape{
				100: {height: 2, left: 50, right: 150},
				50:  {height: 1, left: 30},
				150: {height: 0},
				30:  {height: 0},
			},
			root: 100,
		},
		{
			shape: map[int]shape{
				100: {height: 0},
				50:  {height: 1, left: 30, right: 100},
				30:  {height: 0},
			},
			root: 50,
		},
		{
			shape: map[int]shape{
				100: {height: 0},
				50:  {height: 0},
				80:  {height: 1, left: 50, right: 100},
			},
			root: 80,
		},
		{
			shape: map[int]shape{
				100: {height: 0},
				150: {height: 1, left: 100, right: 180},
				180: {height: 0},
			},
			root: 150,
		},
		{
			shape: map[int]shape{
				100: {height: 0},
				150: {height: 0},
				110: {height: 1, left: 100, right: 150},
			},
			root: 110,
		},
		{
			shape: map[int]shape{
				100: {height: 2, left: 30, right: 150},
				50:  {height: 0},
				150: {height: 0},
				30:  {height: 1, left: 10, right: 50},
				10:  {height: 0},
			},
			root: 100,
		},
		{
			shape: map[int]shape{
				100: {height: 2, left: 50, right: 180},
				50:  {height: 0},
				150: {height: 0},
				180: {height: 1, left: 150, right: 200},
				200: {height: 0},
			},
			root: 100,
		},
		{
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
			root: 100,
		},
		{
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
			root: 100,
		},
		{
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
			root: 100,
		},
		{
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
			root: 100,
		},
		{
			shape: map[int]shape{
				100: {height: 1, left: 70, right: 150},
				50:  {height: 2, left: 30, right: 100},
				150: {height: 0},
				30:  {height: 1, left: 10},
				70:  {height: 0},
				10:  {height: 0},
			},
			root: 50,
		},
		{
			shape: map[int]shape{
				100: {height: 1, left: 50, right: 130},
				50:  {height: 0},
				150: {height: 2, left: 100, right: 180},
				130: {height: 0},
				180: {height: 1, right: 200},
				200: {height: 0},
			},
			root: 150,
		},
		{
			shape: map[int]shape{
				100: {height: 1, left: 80, right: 150},
				50:  {height: 1, left: 30},
				150: {height: 0},
				30:  {height: 0},
				70:  {height: 2, left: 50, right: 100},
				80:  {height: 0},
			},
			root: 70,
		},
		{
			shape: map[int]shape{
				100: {height: 1, left: 50, right: 110},
				50:  {height: 0},
				150: {height: 1, right: 180},
				130: {height: 2, left: 100, right: 150},
				180: {height: 0},
				110: {height: 0},
			},
			root: 130,
		},
		{
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
			root: 8,
		},
		{
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
			root: 32,
		},
		{
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
			root: 25,
		},
	}

	for i, c := range cases {
		c := c
		t.Run(
			c.name,
			func() {
				np := NewMapNodePool()
				for k, v := range c.shape {
					node := newExampleNode(k)
					node.height = v.height

					if v.left > 0 {
						node.left = nodeIntKey(v.left)
					}
					if v.right > 0 {
						node.right = nodeIntKey(v.right)
					}

					_ = np.Set(node)
				}

				tr, err := NewTree(nodeIntKey(c.root), np)
				t.NoError(err, "case %d: NewTree failed", i)
				if err != nil {
					return
				}

				t.Equal(nodeIntKey(c.root), tr.Root().Key(), "case %d: root key mismatch", i)
				t.NoError(tr.IsValid(), "case %d: IsValid() failed", i)
			},
		)
	}
}

func TestTree(t *testing.T) {
	suite.Run(t, new(testTree))
}

type shape struct {
	height int16
	left   int
	right  int
}

type shapeCase struct {
	name    string
	keys    []int
	shape   map[int]shape
	updated []int
	root    int
}

func parseNodeIntKey(b []byte) int {
	i, _ := strconv.ParseInt(string(b), 10, 64)
	return int(i)
}

func nodeIntKey(i int) []byte {
	return []byte(fmt.Sprintf("%03d", i))
}

func newExampleNode(i int) *ExampleNode {
	return &ExampleNode{key: nodeIntKey(i)}
}

func newExampleMutableNode(i int) *ExampleMutableNode {
	return &ExampleMutableNode{key: nodeIntKey(i)}
}

var printCount int32

func printTree(tg *TreeGenerator, verbose bool) error {
	defer atomic.AddInt32(&printCount, 1)

	var b bytes.Buffer

	tr, err := tg.Tree()
	if err != nil {
		return err
	}

	PrintDotGraph(tr, &b)

	count := atomic.LoadInt32(&printCount)
	_ = ioutil.WriteFile(
		fmt.Sprintf("/tmp/%03d.dot", count),
		b.Bytes(), 0644,
	)

	if verbose {
		fmt.Fprintln(os.Stderr, b.String())
	}

	return nil
}
