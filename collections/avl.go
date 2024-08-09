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
	// TODO:
}

func (t *AvlTree[T]) Join() {
	// TODO:
}

func (t *AvlTree[T]) Split() {
	// TODO:
}

func (t *AvlTree[T]) Union() {
	// TODO:
}

func (t AvlTree[T]) InsertList(values ...T) {
	for _, value := range values {
		t.Insert(value)
	}
}

func (t *AvlTree[T]) Insert(value T) {
	// step 1: add value to the tree
	var iteratorParent *AvlNode[T]
	iterator := t.root

	for iterator != nil {
		iteratorParent = iterator
		if value == iterator.Value {
			return
		}
		if value < iterator.Value {
			iterator = iterator.Left
			continue
		}
		iterator = iterator.Right
	}

	node := &AvlNode[T]{value, nil, nil, iteratorParent, 0}
	defer t.balance(node)
	if iteratorParent == nil {
		t.root = node
		return
	}
	if value < iteratorParent.Value {
		iteratorParent.Left = node
		return
	}
	iteratorParent.Right = node
}

func (t *AvlTree[T]) balance(inserted *AvlNode[T]) {
	node := inserted
	for parent := node.Parent; parent != nil; parent = node.Parent {
		//var pivot *AvlNode[T]
		if parent.Left == node {
			if t.balanceLeftSubtree(parent, node) {
				break
			}
			node = parent
			continue
		}
		if t.balanceRightSubtree(parent, node) {
			break
		}
		node = parent
	}
}

func (t *AvlTree[T]) balanceLeftSubtree(parent, node *AvlNode[T]) (quit bool) {
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
	var pivot *AvlNode[T]
	//defer t.updateSubtreeParent(pivot, parent, subtreeParent)
	//  left is heavy
	if node.balanceFactor < 0 {
		// right is heavy
		pivot = parent.rotateLeftRight()
		//t.setLeftSubtreeParent(pivot)
		t.updateSubtreeParent(pivot)
		return true
	}
	pivot = parent.rotateRight()
	//t.trySettingRoot(pivot)
	t.updateSubtreeParent(pivot)
	return true
}

func (t *AvlTree[T]) balanceRightSubtree(parent, node *AvlNode[T]) (quit bool) {
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
	var pivot *AvlNode[T]
	//defer t.setRightSubtreeParent(pivot)
	// parent is right heavy
	if node.balanceFactor > 0 {
		// node is left heavy
		pivot = parent.rotateRightLeft()
		//t.setRightSubtreeParent(pivot)
		t.updateSubtreeParent(pivot)
		return true
	}
	pivot = parent.rotateLeft()
	t.updateSubtreeParent(pivot)
	return true
}

// Sets pivot as the tree root if pivot.Parent is nil
func (t *AvlTree[T]) updateSubtreeParent(pivot *AvlNode[T]) {
	if pivot.Parent == nil {
		t.root = pivot
	}
}

func (pivot *AvlNode[T]) setPivotParent(root *AvlNode[T]) {
	pivot.Parent = root.Parent
	if root.Parent != nil {
		if root.Parent.Left == root {
			root.Parent.Left = pivot
		} else {
			root.Parent.Right = pivot
		}
	}
}

// Performs a left rotation on (root *AvlNode[T])'s subtree,
// which makes (pivot *AvlNode[T]) the new subtree root
func (root *AvlNode[T]) rotateLeft() (pivot *AvlNode[T]) {
	// Todo: balance
	pivot = root.Right
	innerChild := pivot.Left

	root.Right = innerChild
	if innerChild != nil {
		// Reposition inner child: pivot.Left -> root.Right
		innerChild.Parent = root
	}

	// set root as pivot's right child
	pivot.Left = root
	pivot.setPivotParent(root)

	// TODO: balance factor
	if pivot.balanceFactor == 0 {
		root.balanceFactor = 1
		pivot.balanceFactor = -1
	} else {
		root.balanceFactor = 0
		pivot.balanceFactor = 0
	}

	//subtreeParent := root.Parent
	root.Parent = pivot
	//if subtreeParent == nil {
	//	return
	//}
	//if subtreeParent.Left == root {
	//	subtreeParent.Left = pivot
	//	return
	//}
	//subtreeParent.Right = pivot
	return
}

// Performs a right rotation on (root *AvlNode[T])'s subtree,
// which makes (pivot *AvlNode[T]) the new subtree root
func (root *AvlNode[T]) rotateRight() (pivot *AvlNode[T]) {
	// Todo: balance
	pivot = root.Left
	innerChild := pivot.Right

	root.Left = innerChild
	if innerChild != nil {
		// Reposition inner child: pivot.Right -> root.Left
		innerChild.Parent = root
	}

	// set root as pivot's right child
	pivot.Right = root
	//pivot.Parent = root.Parent
	pivot.setPivotParent(root)

	// TODO: balance factor
	if pivot.balanceFactor == 0 {
		root.balanceFactor = 1
		pivot.balanceFactor = -1
	} else {
		root.balanceFactor = 0
		pivot.balanceFactor = 0
	}

	//subtreeParent := root.Parent
	root.Parent = pivot
	//if subtreeParent == nil {
	//	return
	//}
	//if subtreeParent.Left == root {
	//	subtreeParent.Left = pivot
	//	return
	//}
	//subtreeParent.Right = pivot
	return
}

func (root *AvlNode[T]) rotateLeftRight() (pivot *AvlNode[T]) {
	// root is by 2 higher than its sibling
	pivotRoot, pivot := root.Left, root.Left.Right
	// 1 - Reorder nodes
	//pivot.Parent = root.Parent
	pivot.setPivotParent(root)
	// y is by 1 higher than sibling
	t2 := pivot.Left
	pivotRoot.Right = t2
	if t2 != nil {
		t2.Parent = pivotRoot
	}
	pivot.Left = pivotRoot
	pivotRoot.Parent = pivot

	t3 := pivot.Right
	root.Left = t3
	if t3 != nil {
		t3.Parent = root
	}
	pivot.Right = root
	root.Parent = pivot

	// 2 - Revaluate balance factors
	if pivot.balanceFactor == 0 {
		root.balanceFactor = 0
		pivotRoot.balanceFactor = 0
	} else if pivot.balanceFactor > 0 {
		// t3 was higher
		root.balanceFactor = -1
		pivotRoot.balanceFactor = 0
	} else {
		// t2 was higher
		root.balanceFactor = 0
		pivotRoot.balanceFactor = 1
	}
	pivot.balanceFactor = 0
	return
}

// Performs a right-left rotation on (root *AvlNode[T])'s subtree,
// which makes (pivot *AvlNode[T]) the new subtree root
func (root *AvlNode[T]) rotateRightLeft() (pivot *AvlNode[T]) {
	// root is by 2 higher than its sibling
	pivotRoot, pivot := root.Right, root.Right.Left
	// 1 - Reorder nodes
	//pivot.Parent = root.Parent
	pivot.setPivotParent(root)
	// y is by 1 higher than sibling
	t2 := pivot.Left
	root.Right = t2
	if t2 != nil {
		t2.Parent = root
	}
	pivot.Left = root
	root.Parent = pivot

	t3 := pivot.Right
	pivotRoot.Left = t3
	if t3 != nil {
		t3.Parent = pivotRoot
	}
	pivot.Right = pivotRoot
	pivotRoot.Parent = pivot

	// 2 - Revaluate balance factors
	// 1st case BF(y) == 0
	if pivot.balanceFactor == 0 {
		root.balanceFactor = 0
		pivotRoot.balanceFactor = 0
	} else if pivot.balanceFactor > 0 {
		// t3 was higher
		root.balanceFactor = 0
		pivotRoot.balanceFactor = -1
	} else {
		// t2 was higher
		root.balanceFactor = 1
		pivotRoot.balanceFactor = 0
	}
	pivot.balanceFactor = 0
	return
}

func (t *AvlTree[T]) String() string {
	return t.root.String()
}

func (node *AvlNode[T]) String() string {
	var buffer bytes.Buffer
	//var buffer safeBuffer
	node.string(&buffer, 0, 'M')
	return buffer.String()
}

// func (node *AvlNode[T]) string(buffer *safeBuffer, spaces int, ch rune) {
func (node *AvlNode[T]) string(buffer *bytes.Buffer, spaces int, ch rune) {
	if node == nil {
		return
	}
	//buffer.mu.Lock()
	for i := 0; i < spaces; i++ {
		buffer.WriteByte(' ')
	}
	if node.Parent != nil {
		fmt.Fprintf(buffer, "%c:%v \t\tP:%v\tBF:%v\n", ch, node.Value, node.Parent.Value, node.balanceFactor)
	} else {
		fmt.Fprintf(buffer, "%c:%v \t\tP:nil\tBF:%v\n", ch, node.Value, node.balanceFactor)
	}
	//buffer.mu.Unlock()

	node.Left.string(buffer, spaces+2, 'L')
	node.Right.string(buffer, spaces+2, 'R')
}
