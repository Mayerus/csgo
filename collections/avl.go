package collections

import (
	"bytes"
	"cmp"
	"fmt"
	"math"
	"sync"
)

type AvlTree[T cmp.Ordered] struct {
	root *AvlNode[T]
}

// type AvlNode[T Numeric] struct {
type AvlNode[T cmp.Ordered] struct {
	Value         T
	Left, Right   *AvlNode[T]
	Parent        *AvlNode[T]
	balanceFactor int8
}

type safeBuffer struct {
	bytes.Buffer
	mu sync.Mutex
}

func (n *AvlNode[T]) Dispose() {
	n.Parent = nil
	n.Left = nil
	n.Right = nil
}

func (adopter *AvlNode[T]) Adopt(giver *AvlNode[T]) {
	adopter.Left = giver.Left
	if adopter.Left != nil {
		adopter.Left.Parent = adopter
	}
	adopter.Right = giver.Right
	if adopter.Right != nil {
		adopter.Right.Parent = adopter
	}
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

func (t *AvlTree[T]) Min() (node *AvlNode[T]) {
	if t.root == nil {
		return nil
	}
	return t.root.Min()
}

func (n *AvlNode[T]) Min() (node *AvlNode[T]) {
	node = n
	for node.Left != nil {
		node = node.Left
	}
	return
}

func (t *AvlTree[T]) Max() (node *AvlNode[T]) {
	if t.root == nil {
		return nil
	}
	return t.root.Max()
}

func (n *AvlNode[T]) Max() (node *AvlNode[T]) {
	node = n
	for node.Right != nil {
		node = node.Right
	}
	return
}

func (tree *AvlTree[T]) Height() int {
	return tree.root.height()
}

func (node *AvlNode[T]) height() int {
	if node == nil {
		return -1
	}
	return int(math.Max(float64(node.Left.height()), float64(node.Right.height()))) + 1
}

// The successor of node n is the node with the smallest value greater than n's.
func (n *AvlNode[T]) Successor() *AvlNode[T] {
	if n.Right != nil {
		return n.Right.Min()
	}
	node := n
	successor := n.Parent
	for successor != nil && node == successor.Right {
		node = successor
		successor = successor.Parent
	}
	return successor
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
	defer t.retraceInsert(node)
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

func (t *AvlTree[T]) retraceInsert(inserted *AvlNode[T]) {
	node := inserted
	for parent := node.Parent; parent != nil; parent = node.Parent {
		//var pivot *AvlNode[T]
		if parent.Left == node {
			if t.balanceLeftSubtreeInsert(parent, node) {
				break
			}
			node = parent
			continue
		}
		if t.balanceRightSubtreeInsert(parent, node) {
			break
		}
		node = parent
	}
}

func (t *AvlTree[T]) balanceLeftSubtreeInsert(parent, node *AvlNode[T]) (quit bool) {
	if parent.balanceFactor < 0 {
		// the newly inserted node balanced the subtree
		parent.balanceFactor = 0
		return true
	}
	if parent.balanceFactor == 0 {
		parent.balanceFactor = 1
		return false
	}
	var pivot *AvlNode[T]
	//defer t.updateSubtreeParent(pivot, parent, subtreeParent)
	//  left is heavy
	if node.balanceFactor < 0 {
		// right is heavy
		pivot = parent.rotateLeftRight()
		t.updateSubtreeParent(pivot)
		return true
	}
	pivot = parent.rotateRight()
	t.updateSubtreeParent(pivot)
	return true
}

func (t *AvlTree[T]) balanceRightSubtreeInsert(parent, node *AvlNode[T]) (quit bool) {
	if parent.balanceFactor > 0 {
		//the newly inserted node balanced the subtree
		parent.balanceFactor = 0
		return true
	}
	if parent.balanceFactor == 0 {
		parent.balanceFactor = -1
		return false
	}
	var pivot *AvlNode[T]
	//defer t.setRightSubtreeParent(pivot)
	// parent is right heavy
	if node.balanceFactor > 0 {
		// node is left heavy
		pivot = parent.rotateRightLeft()
		t.updateSubtreeParent(pivot)
		return true
	}
	pivot = parent.rotateLeft()
	t.updateSubtreeParent(pivot)
	return true
}

func (t *AvlTree[T]) Delete(value T) bool {
	// remove node
	if t.Search(value) == nil {
		return false
	}

	//pivot := t.deleteNode(t.root, value)
	//if pivot != nil {
	//	pivot.Parent = nil
	//	t.root = pivot
	//}
	return t.deleteNodeIter(value)
	return true
}

func (t *AvlTree[T]) deleteNode(node *AvlNode[T], value T) *AvlNode[T] {
	// TODO: make the function iterative instead of recursive
	if node == nil {
		return nil
	}
	if value < node.Value {
		node.Left = t.deleteNode(node.Left, value)
		if node.Left != nil {
			node.Left.Parent = node
		}
	} else if value > node.Value {
		node.Right = t.deleteNode(node.Right, value)
		if node.Right != nil {
			node.Right.Parent = node
		}
	} else {
		// Node with no child or one child
		if node.Left == nil {
			successor := node.Right
			node.Dispose()
			node = nil
			return successor
		} else if node.Right == nil {
			successor := node.Left
			node.Dispose()
			node = nil
			return successor
		}
		// Node with two children: Get the in-order successor (smallest in the right subtree)
		successor := node.Successor()
		node.Value = successor.Value

		node.Right = t.deleteNode(node.Right, successor.Value)
		if node.Right != nil {
			node.Right.Parent = node
		}
	}

	//if node == nil {
	//	return nil
	//}
	t.updateSubtreeParent(node)

	// balance node
	return node.balanceDelete()
}

func (t *AvlTree[T]) deleteNodeIter( /*node *AvlNode[T],*/ value T) bool {
	if t.root == nil {
		return false
	}

	stack := &Stack[*AvlNode[T]]{}
	current := t.root
	//callStack.Push(t.root)
	var parent *AvlNode[T]

	// traverse to the we wish to delete and push its ancestorial branch
	for current != nil {
		if value < current.Value {
			stack.Push(current)
			parent = current
			current = current.Left
			continue
		}
		if value > current.Value {
			stack.Push(current)
			parent = current
			current = current.Right
			continue
		}
		break
	}

	// return if no node corresponds to the given value
	if current == nil {
		return false
	}

	// Leaf or a node with single child
	if current.Left == nil || current.Right == nil {
		if t.removeLeafOrSingleChildNode(current, parent) {
			return true
		}
	} else {
		remove2ChildrenNode(current, parent, stack)
	}

	// Balance the nodes on the path to the root
	for !stack.Empty() {
		current, _ = stack.Pop()
		t.updateSubtreeParent(current)
		current = current.balanceDelete()
		t.updateSubtreeParent(current)
	}

	return true
}

// Returns whether a rebalance is required or not
func (t *AvlTree[T]) removeLeafOrSingleChildNode(node, parent *AvlNode[T]) bool {
	var successor *AvlNode[T]
	if node.Left == nil {
		successor = node.Right
	} else {
		successor = node.Left
	}

	if parent == nil {
		// Deleting the root node with a single child
		if successor != nil {
			successor.Parent = nil
		}

		node.Dispose()
		//if successor != nil {
		t.updateSubtreeParent(successor)
		//}
		return true
	}

	// Update parent's reference
	if parent.Left == node {
		parent.Left = successor
	} else {
		parent.Right = successor
	}
	if successor != nil {
		successor.Parent = parent
	}

	node.Dispose()
	return false
}

func remove2ChildrenNode[T cmp.Ordered](node, parent *AvlNode[T], stack *Stack[*AvlNode[T]]) {
	stack.Push(node)
	var successor *AvlNode[T]
	for successor = node.Right; successor.Left != nil; successor = successor.Left {
		stack.Push(successor)
	}

	node.Value = successor.Value

	parent, _ = stack.Peek()
	if parent.Left == successor {
		parent.Left = successor.Right
	} else {
		parent.Right = successor.Right
	}

	if successor.Right != nil {
		successor.Right.Parent = parent
	}

	successor.Dispose()
}

// Returns a balanced subtree with "node" being the original root
// before the possible rebalance action
func (node *AvlNode[T]) balanceDelete() *AvlNode[T] {
	//TODO: consider using the existing balance factor?
	// Update height and balance factor
	node.balanceFactor = int8(node.Left.height() - node.Right.height())

	if node.balanceFactor > 1 {
		if node.Left.balanceFactor >= 0 {
			return node.rotateRight()
		}
		return node.rotateLeftRight()
	}
	if node.balanceFactor < -1 {
		if node.Right.balanceFactor <= 0 {
			return node.rotateLeft()
		}
		return node.rotateRightLeft()
	}
	return node
}

//func Join[T Numeric](tL, tR *AvlTree[T], k T) (bool, *AvlTree[T]) {
//	if tL.Max().Value >= k || tR.Min().Value <= k {
//		return false, nil
//	}
//
//	if tL.Height() > tR.Height() {
//		return true, joinRightAvl(tL, tR, k)
//	}
//	if tL.Height() < tR.Height() {
//		return true, joinLefttAvl(tL, tR, k)
//	}
//
//	root := &AvlNode[T]{Value: k,
//		Left:          tL.root,
//		Right:         tR.root,
//		Parent:        nil,
//		balanceFactor: 0,
//	}
//	tree := &AvlTree[T]{root}
//	return true, tree
//}
//
//func joinRightAvl[T Numeric](tL, tR *AvlTree[T], k T) *AvlTree[T] {
//	// l, c, kTag := tL.root.Left, tL.root.Right, tL.Root.Value
//	l, c, kTag := tL.expose()
//	// stopping condition - if height difference is 1 at most (left may be heigher than right)
//
//	balanceFactor := tL.root.Right.height() - tR.Height()
//	//if tL.root.Right.height() <= tR.Height()+1 {
//	if balanceFactor <= 1 {
//		joinRoot := &AvlNode[T]{Value: k,
//			Left:          tL.root.Right,
//			Right:         tR.root,
//			balanceFactor: int8(balanceFactor),
//		} // TODO: restructure properly
//
//		joinBF := joinRoot.height() - tL.root.Left.height()
//
//		//if joinRoot.height() <= tL.root.Left.height()+1 {
//
//		if joinBF <= 1 {
//			joinRoot = &AvlNode[T]{Value: kTag,
//				Left:          l,
//				Right:         joinRoot,
//				balanceFactor: int8(joinBF),
//			}
//			return &AvlTree[T]{root: joinRoot}
//		}
//		pivot := tTag.rotateRight()
//		node := &AvlNode[T]{Value: kTag, Left: l, Right: pivot, balanceFactor: 0}
//		return node.rotateLeft()
//	}
//	tTag := joinRightAvl(c, tR, k)
//	tTagTag := &AvlNode[T]{Value: kTag, Left: l, Right: tTag} // restructure properly
//	if tTag.Height() <= l.Height()+1 {
//		return tTagTag
//	}
//	return tTagTag.rotateLeft()
//}
//
//func (t *AvlTree[T]) expose() (l, r *AvlNode[T], k T) {
//
//}
//func joinLeftAvl[T Numeric](tL, tR *AvlTree[T], k T) *AvlTree[T] {
//
//}
//
//func (t *AvlTree[T]) Split() {
//	// TODO:
//}
//
//func (t *AvlTree[T]) Union() {
//	// TODO:
//}

// Sets pivot as the tree root if pivot.Parent is nil
func (t *AvlTree[T]) updateSubtreeParent(pivot *AvlNode[T]) {
	if pivot == nil {
		t.root = nil
		return
	}
	if pivot.Parent == nil {
		t.root = pivot
		return
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

	if pivot.balanceFactor == 0 {
		root.balanceFactor = -1
		pivot.balanceFactor = 1
	} else {
		root.balanceFactor = 0
		pivot.balanceFactor = 0
	}

	root.Parent = pivot
	return
}

// Performs a right rotation on (root *AvlNode[T])'s subtree,
// which makes (pivot *AvlNode[T]) the new subtree root
func (root *AvlNode[T]) rotateRight() (pivot *AvlNode[T]) {
	pivot = root.Left
	innerChild := pivot.Right

	root.Left = innerChild
	if innerChild != nil {
		// Reposition inner child: pivot.Right -> root.Left
		innerChild.Parent = root
	}

	// set root as pivot's right child
	pivot.Right = root
	pivot.setPivotParent(root)

	if pivot.balanceFactor == 0 {
		root.balanceFactor = 1
		pivot.balanceFactor = -1
	} else {
		root.balanceFactor = 0
		pivot.balanceFactor = 0
	}

	root.Parent = pivot
	return
}

func (root *AvlNode[T]) rotateLeftRight() (pivot *AvlNode[T]) {
	// root is by 2 higher than its sibling
	pivotRoot, pivot := root.Left, root.Left.Right
	// 1 - Reorder nodes
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
	buffer := safeBuffer{}
	node.string(&buffer, 0, 'M')
	return buffer.String()
}

func (node *AvlNode[T]) string(buffer *safeBuffer, spaces int, ch rune) {
	if node == nil {
		return
	}
	buffer.mu.Lock()
	for i := 0; i < spaces; i++ {
		buffer.WriteByte(' ')
	}
	if node.Parent != nil {
		fmt.Fprintf(buffer, "%c:%v \t\tP:%v\tBF:%v\n", ch, node.Value, node.Parent.Value, node.balanceFactor)
	} else {
		fmt.Fprintf(buffer, "%c:%v \t\tP:nil\tBF:%v\n", ch, node.Value, node.balanceFactor)
	}
	buffer.mu.Unlock()

	node.Left.string(buffer, spaces+2, 'L')
	node.Right.string(buffer, spaces+2, 'R')
}
