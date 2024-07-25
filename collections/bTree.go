package collections

import (
	"bytes"
	"fmt"

	"golang.org/x/exp/constraints"
)

type Numeric interface {
	constraints.Float | constraints.Integer
}

type TreeNode[T Numeric] struct {
	Value               T
	Left, Right, Parent *TreeNode[T]
}

func (node *TreeNode[T]) String() string {
	var buffer bytes.Buffer
	node.string(&buffer, 0, 'M')
	return buffer.String()
}

func (node *TreeNode[T]) string(buffer *bytes.Buffer, spaces int, ch rune) {
	if node == nil {
		return
	}
	for i := 0; i < spaces; i++ {
		buffer.WriteByte(' ')
	}
	if node.Parent != nil {
		fmt.Fprintf(buffer, "%c:%v \tP:%v\n", ch, node.Value, node.Parent.Value)
	} else {
		fmt.Fprintf(buffer, "%c:%v \tP:nil\n", ch, node.Value)
	}

	node.Left.string(buffer, spaces+2, 'L')
	node.Right.string(buffer, spaces+2, 'R')
}

func (t *TreeNode[T]) InorderTraversal() {
	if t == nil {
		return
	}
	t.Left.InorderTraversal()
	// visit node
	t.Right.InorderTraversal()

}

func (t *TreeNode[T]) PreorderTraversal() {
	if t == nil {
		return
	}
	// visit node
	t.Left.PreorderTraversal()
	t.Right.PreorderTraversal()
}

func (t *TreeNode[T]) PostorderTraversal() {
	if t == nil {
		return
	}
	t.Left.PostorderTraversal()
	t.Right.PostorderTraversal()
	// visit node
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
	n.Parent = nil
}
