package avl_test

import (
	"fmt"
	"os"

	"github.com/spikeekips/avl"
)

var numberOfNodesForDot = 10

// ExamplePrintDotGraph print a dot graph from the generated Tree.
func ExamplePrintDotGraph() {
	// create new TreeGenerator
	tg := avl.NewTreeGenerator()

	// generate 10 new MutableNodes and add to TreeGenerator.
	for i := 0; i < numberOfNodesForDot; i++ {
		node := &ExampleMutableNode{
			key: []byte(fmt.Sprintf("%03d", i)),
		}
		if _, err := tg.Add(node); err != nil {
			return
		}
	}

	// Get Tree from TreeGenerator.
	tree, err := tg.Tree()
	if err != nil {
		return
	}

	avl.PrintDotGraph(tree, os.Stdout)
	// Output:
	// graph graphname {
	//   "003" [label="003 (3)"];
	//   "003" -- "001";
	//   "001" [label="001 (1)"];
	//   "001" -- "000";
	//   "000" [label="000 (0)"];
	//   "0000" [label=" " style="filled" color="white" bgcolor="white"];
	//   "000" -- "0000" [style="solid" color="white" bgcolor="white"];
	//   "0001" [label=" " style="filled" color="white" bgcolor="white"];
	//   "000" -- "0001" [style="solid" color="white" bgcolor="white"];
	//   "001" -- "002";
	//   "002" [label="002 (0)"];
	//   "0020" [label=" " style="filled" color="white" bgcolor="white"];
	//   "002" -- "0020" [style="solid" color="white" bgcolor="white"];
	//   "0021" [label=" " style="filled" color="white" bgcolor="white"];
	//   "002" -- "0021" [style="solid" color="white" bgcolor="white"];
	//   "003" -- "007";
	//   "007" [label="007 (2)"];
	//   "007" -- "005";
	//   "005" [label="005 (1)"];
	//   "005" -- "004";
	//   "004" [label="004 (0)"];
	//   "0040" [label=" " style="filled" color="white" bgcolor="white"];
	//   "004" -- "0040" [style="solid" color="white" bgcolor="white"];
	//   "0041" [label=" " style="filled" color="white" bgcolor="white"];
	//   "004" -- "0041" [style="solid" color="white" bgcolor="white"];
	//   "005" -- "006";
	//   "006" [label="006 (0)"];
	//   "0060" [label=" " style="filled" color="white" bgcolor="white"];
	//   "006" -- "0060" [style="solid" color="white" bgcolor="white"];
	//   "0061" [label=" " style="filled" color="white" bgcolor="white"];
	//   "006" -- "0061" [style="solid" color="white" bgcolor="white"];
	//   "007" -- "008";
	//   "008" [label="008 (1)"];
	//   "0081" [label=" " style="filled" color="white" bgcolor="white"];
	//   "008" -- "0081" [style="solid" color="white" bgcolor="white"];
	//   "008" -- "009";
	//   "009" [label="009 (0)"];
	//   "0090" [label=" " style="filled" color="white" bgcolor="white"];
	//   "009" -- "0090" [style="solid" color="white" bgcolor="white"];
	//   "0091" [label=" " style="filled" color="white" bgcolor="white"];
	//   "009" -- "0091" [style="solid" color="white" bgcolor="white"];
	// }
}
