package avl

import (
	"fmt"
	"io"
	"os"
)

func printDotGraphNode(w io.Writer, tr *Tree, node Node) {
	fmt.Fprintf(w, `  "%s" [label="%s (%d)"];
`,
		string(node.Key()),
		string(node.Key()),
		node.Height(),
	)

	if node.LeftKey() == nil && node.RightKey() == nil {
		fmt.Fprintf(
			w,
			`  "%s" [label=" " style="filled" color="white" bgcolor="white"];
  "%s" -- "%s" [style="solid" color="white" bgcolor="white"];
  "%s" [label=" " style="filled" color="white" bgcolor="white"];
  "%s" -- "%s" [style="solid" color="white" bgcolor="white"];
`,
			string(node.Key())+"0",
			string(node.Key()),
			string(node.Key())+"0",
			string(node.Key())+"1",
			string(node.Key()),
			string(node.Key())+"1",
		)

		return
	}

	if node.LeftKey() != nil {
		fmt.Fprintf(
			w,
			`  "%s" -- "%s";
`,
			string(node.Key()),
			string(node.LeftKey()),
		)

		if node.RightKey() == nil {
			fmt.Fprintf(
				w,
				`  "%s" [label=" " style="filled" color="white" bgcolor="white"];
  "%s" -- "%s" [style="solid" color="white" bgcolor="white"];
`,
				string(node.Key())+"0",
				string(node.Key()),
				string(node.Key())+"0",
			)
		}

		left, _ := tr.NodePool().Get(node.LeftKey())
		printDotGraphNode(w, tr, left)
	}

	if node.RightKey() != nil {
		if node.LeftKey() == nil {
			fmt.Fprintf(
				w,
				`  "%s" [label=" " style="filled" color="white" bgcolor="white"];
  "%s" -- "%s" [style="solid" color="white" bgcolor="white"];
`,
				string(node.Key())+"1",
				string(node.Key()),
				string(node.Key())+"1",
			)
		}

		fmt.Fprintf(
			w,
			`  "%s" -- "%s";
`,
			string(node.Key()),
			string(node.RightKey()),
		)

		right, _ := tr.NodePool().Get(node.RightKey())
		printDotGraphNode(w, tr, right)
	}
}

func PrintDotGraph(tr *Tree, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}

	fmt.Fprintln(w, "graph graphname {")
	defer fmt.Fprintln(w, "}")

	if tr.root == nil {
		return
	}

	printDotGraphNode(w, tr, tr.root)
}
