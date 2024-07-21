package collections

import (
	"fmt"
	"io"

	"golang.org/x/exp/constraints"
)

type Numeric interface {
	constraints.Float | constraints.Integer
}

type TreeNode[T Numeric] struct {
	Value T
	Left  *TreeNode[T]
	Right *TreeNode[T]
}

func (n *TreeNode[T]) String() string {
	if n == nil {
		return "[]"
	}
	s := ""
	if n.Left != nil {
		s += n.Left.String() + " "
	}
	s += fmt.Sprintf("%v", n.Value)
	if n.Right != nil {
		s += " " + n.Right.String()
	}
	return fmt.Sprintf("[%s]", s)

}

func (node *TreeNode[T]) Print(w io.Writer, spaces int, ch rune) {
	if node == nil {
		return
	}
	fmt.Println("test")
	for i := 0; i < spaces; i++ {
		fmt.Fprint(w, " ")
	}
	fmt.Fprintf(w, "%c:%v\n", ch, node.Value)
	node.Left.Print(w, spaces+2, 'L')
	node.Right.Print(w, spaces+2, 'R')
}

func (n *TreeNode[T]) IsLeaf() bool {
	return n.Left == nil && n.Right == nil
}

func (n *TreeNode[T]) HasLeft() bool {
	return n.Left != nil
}

func (n *TreeNode[T]) HasRight() bool {
	return n.Right != nil
}

func (n *TreeNode[T]) Dispose() {
	n.Right = nil
	n.Left = nil
}
