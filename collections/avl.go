package collections

import (
	"bytes"
	"fmt"
)

type AvlTree[T Numeric] struct {
	root *AvlNode[T]
}

// balanceFactor = left height = right height
type AvlNode[T Numeric] struct {
	Value         T
	Left, Right   *AvlNode[T]
	Parent        *AvlNode[T]
	balanceFactor int
}

func (t *AvlTree[T]) Search(value T) (node *AvlNode[T]) {
	node = t.root
	for node != nil && value != node.Value {
		if value < node.Value {
			node = node.Left
			continue
		}
		node = node.Right
	}
	return
}

func (t *AvlTree[T]) Delete() {

}

func (t *AvlTree[T]) Insert(value T) {
	// step 1: add value to the tree
	var iteratorParent *AvlNode[T]
	iterator := t.root

	for iterator != nil {
		iteratorParent = iterator
		if value < iterator.Value {
			iterator = iterator.Left
			continue
		}
		iterator = iterator.Right
	}
	node := &AvlNode[T]{value, nil, nil, iteratorParent, 0}
	if iteratorParent == nil {
		t.root = node
		return
	}
	if value < iteratorParent.Value {
		iteratorParent.Left = node
		return
	}
	iteratorParent.Right = node

	// step 2: rebalance the tree
	t.balance(node)
}

func (t *AvlTree[T]) balance(inserted *AvlNode[T]) {
	node := inserted
	for parent := node.Parent; parent != nil; parent = node.Parent {
		var pivot *AvlNode[T]
		if parent.Left == node {
			if t.balanceLeftSubtree(parent, node, pivot) {
				break
			}
			continue
		}
		if t.balanceRightSubtree(parent, node, pivot) {
			break
		}
	}
}

func (t *AvlTree[T]) balanceLeftSubtree(parent, node, pivot *AvlNode[T]) (quit bool) {
	if parent.balanceFactor < 0 {
		// the newly inserted node balanced the subtree
		parent.balanceFactor = 0
		return true
	}
	if parent.balanceFactor == 0 {
		parent.balanceFactor++
		node = parent
		return false
	}
	//  left is heavy
	if node.balanceFactor < 0 {
		// right is heavy
		pivot = parent.rotateLeftRight()
		return t.setSubtreeParent(pivot)
	}
	pivot = parent.rotateRight()
	return t.setSubtreeParent(pivot)
}

func (t *AvlTree[T]) balanceRightSubtree(parent, node, pivot *AvlNode[T]) (quit bool) {
	if parent.balanceFactor > 0 {
		//the newly inserted node balanced the subtree
		parent.balanceFactor = 0
		return true
	}
	if parent.balanceFactor == 0 {
		parent.balanceFactor--
		node = parent
		return false
	}
	// parent is right heavy
	if node.balanceFactor > 0 {
		// node is left heavy
		pivot = parent.rotateRightLeft()
		return t.setSubtreeParent(pivot)
	}
	pivot = parent.rotateLeft()
	return t.setSubtreeParent(pivot)
}

func (t *AvlTree[T]) setSubtreeParent(pivot *AvlNode[T]) (quit bool) {
	if pivot.Parent != nil {
		return true
	}
	t.root = pivot
	return false
}

// Performs a left rotation on (root *AvlNode[T])'s subtree,
// which makes (pivot *AvlNode[T]) the new subtree root
func (root *AvlNode[T]) rotateLeft() (pivot *AvlNode[T]) {
	// Todo: balance
	pivot = root.Right
	innerChild := pivot.Left

	if innerChild != nil {
		// Reposition inner child: pivot.Left -> root.Right
		root.Right = innerChild
		innerChild.Parent = root
	}

	// set root as pivot's right child
	pivot.Left = root
	pivot.Parent = root.Parent

	// TODO: balance factor
	if pivot.balanceFactor == 0 {
		root.balanceFactor = 1
		pivot.balanceFactor = -1
	} else {
		root.balanceFactor = 0
		pivot.balanceFactor = 0
	}

	subtreeParent := root.Parent
	root.Parent = pivot
	if subtreeParent == nil {
		return
	}
	if subtreeParent.Left == root {
		subtreeParent.Left = pivot
		return
	}
	subtreeParent.Right = pivot
	return
}

// Performs a right rotation on (root *AvlNode[T])'s subtree,
// which makes (pivot *AvlNode[T]) the new subtree root
func (root *AvlNode[T]) rotateRight() (pivot *AvlNode[T]) {
	// Todo: balance
	pivot = root.Left
	innerChild := pivot.Right

	if innerChild != nil {
		// Reposition inner child: pivot.Right -> root.Left
		root.Left = innerChild
		innerChild.Parent = root
	}

	// set root as pivot's right child
	pivot.Right = root
	pivot.Parent = root.Parent

	// TODO: balance factor
	if pivot.balanceFactor == 0 {
		root.balanceFactor = 1
		pivot.balanceFactor = -1
	} else {
		root.balanceFactor = 0
		pivot.balanceFactor = 0
	}

	subtreeParent := root.Parent
	root.Parent = pivot
	if subtreeParent == nil {
		return
	}
	if subtreeParent.Left == root {
		subtreeParent.Left = pivot
		return
	}
	subtreeParent.Right = pivot
	return
}

func (root *AvlNode[T]) rotateLeftRight() (pivot *AvlNode[T]) {
	// TODO:
}

func (root *AvlNode[T]) rotateRightLeft() (pivot *AvlNode[T]) {
	// TODO:
}

func (t *AvlTree[T]) String() string {
	return t.root.String()
}

func (node *AvlNode[T]) String() string {
	var buffer bytes.Buffer
	node.string(&buffer, 0, 'M')
	return buffer.String()
}

func (node *AvlNode[T]) string(buffer *bytes.Buffer, spaces int, ch rune) {
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
