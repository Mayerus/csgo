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
	if t == nil {
		return nil
	}
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
		return 0
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

	return t.deleteNode(value)
}

func (t *AvlTree[T]) deleteNode(value T) bool {
	if t.root == nil {
		return false
	}

	stack := &Stack[*AvlNode[T]]{}
	current := t.root
	var parent *AvlNode[T]

	// traverse to the we wish to delete and push its ancestorial branch
	for !(current == nil || value == current.Value) {
		stack.Push(current)
		parent = current
		if value < current.Value {
			current = current.Left
			continue
		}
		current = current.Right
	}

	// return if no node corresponds to the given value
	if current == nil {
		return false
	}

	quit, leftDeletion := t.removeNode(current, parent, stack)

	// Balance the nodes on the path to the root
	for !(stack.Empty() || quit) {
		current, _ = stack.Pop()
		current, quit = current.balanceDelete(leftDeletion)
		if current.Parent != nil {
			leftDeletion = current.Parent.Left == current
		}
		t.updateSubtreeParent(current)
	}
	return true
}

func (t *AvlTree[T]) removeNode(current, parent *AvlNode[T], stack *Stack[*AvlNode[T]]) (quit, leftDeletion bool) {
	var successor *AvlNode[T]

	// Leaf or a node with single child
	if current.Left == nil || current.Right == nil {
		defer current.Dispose()

		if current.Left == nil {
			successor = current.Right
		} else {
			successor = current.Left
		}
		if successor != nil {
			successor.Parent = parent
		}

		if parent == nil {
			// Deleting the root node with a single child
			t.root = successor
			return true, false
		}

		// Update parent's reference
		if parent.Left == current {
			parent.Left = successor
			return false, true
		}
		parent.Right = successor
		return false, false
	}

	stack.Push(current)
	for successor = current.Right; successor.Left != nil; successor = successor.Left {
		stack.Push(successor)
	}

	current.Value = successor.Value
	parent, _ = stack.Peek()

	if successor.Right != nil {
		successor.Right.Parent = parent
	}
	defer successor.Dispose()

	leftDeletion = successor.Parent.Left == successor

	if parent.Left == successor {
		parent.Left = successor.Right
		return false, leftDeletion
	}
	parent.Right = successor.Right
	return false, leftDeletion
}

// Returns a balanced subtree with "node" being the original root
// before the possible rebalance action
func (node *AvlNode[T]) balanceDelete(leftDeletion bool) (*AvlNode[T], bool) {
	if leftDeletion {
		node.balanceFactor--
		if node.balanceFactor == -1 {
			return node, true
		}
	} else {
		node.balanceFactor++
		if node.balanceFactor == 1 {
			return node, true
		}
	}

	if node.balanceFactor > 1 {
		quit := node.Left.balanceFactor == 0
		if node.Left.balanceFactor >= 0 {
			return node.rotateRight(), quit
		}
		return node.rotateLeftRight(), quit
	}
	if node.balanceFactor < -1 {
		quit := node.Right.balanceFactor == 0
		if node.Right.balanceFactor <= 0 {
			return node.rotateLeft(), quit
		}
		return node.rotateRightLeft(), quit
	}
	return node, false
}

func AvlJoin[T cmp.Ordered](tL, tR *AvlTree[T], k T) (bool, *AvlTree[T]) {
	if tL == nil {
		//tR.Insert(k)
		//return true, tR
		return false, nil
	}
	if tL.root == nil {
		//tR.Insert(k)
		//return true, tR
		return false, nil
	}
	if tR == nil {
		//tL.Insert(k)
		//return true, tL
		return false, nil
	}
	if tR.root == nil {
		//tL.Insert(k)
		//return true, tL
		return false, nil
	}

	//a, b := tL.Max(), tR.Min()
	//fmt.Print(a.Value, b.Value)

	if tL.Max().Value >= k || tR.Min().Value <= k {
		return false, nil
	}

	if tL.Height() > tR.Height() {
		return true, joinRightAvl(tL.root, tR.root, k)
	}
	if tL.Height() < tR.Height() {
		return true, joinLeftAvl(tL.root, tR.root, k)
	}

	root := &AvlNode[T]{Value: k,
		Left:          tL.root,
		Right:         tR.root,
		Parent:        nil,
		balanceFactor: 0,
	}
	tree := &AvlTree[T]{root}
	return true, tree
}

func joinRightAvl[T cmp.Ordered](tL, tR *AvlNode[T], joinValue T) *AvlTree[T] {
	pivotLeft, pivotRight, pivotValue := tL.Left, tL.Right, tL.Value

	// stopping condition - if height difference is 1 at most (left may be heigher than right)
	joinNodeBF := pivotRight.height() - tR.height()
	if joinNodeBF <= 1 {
		joinNode := &AvlNode[T]{
			Value:         joinValue,
			Left:          pivotRight,
			Right:         tR,
			balanceFactor: int8(joinNodeBF),
		}
		if pivotRight != nil {
			pivotRight.Parent = joinNode
		}
		tR.Parent = joinNode

		joinRootBF := pivotLeft.height() - joinNode.height()
		if joinRootBF >= -1 {
			// height difference is 1 at most, complete the join
			joinRoot := &AvlNode[T]{
				Value:         pivotValue,
				Left:          pivotLeft,
				Right:         joinNode,
				Parent:        nil,
				balanceFactor: int8(joinRootBF),
			}
			joinNode.Parent = joinRoot
			pivotLeft.Parent = joinRoot
			return &AvlTree[T]{root: joinRoot}
		}

		joinNode = joinNode.rotateRightUnmodifiedTree()

		// height difference is yet too great
		joinRoot := &AvlNode[T]{
			Value:         pivotValue,
			Left:          pivotLeft,
			Right:         joinNode,
			Parent:        nil,
			balanceFactor: int8(pivotLeft.height()) - int8(joinNode.height()),
		}
		if pivotLeft != nil {
			pivotLeft.Parent = joinRoot
		}
		joinRoot.Right.Parent = joinRoot

		return &AvlTree[T]{joinRoot.rotateLeftUnmodifiedTree()}
	}

	joinedRightTree := joinRightAvl(pivotRight, tR, joinValue)
	joinedRoot := &AvlNode[T]{
		Value:         pivotValue,
		Left:          pivotLeft,
		Right:         joinedRightTree.root,
		balanceFactor: int8(pivotLeft.height()) - int8(joinedRightTree.Height()),
	}
	pivotLeft.Parent = joinedRoot
	joinedRoot.Right.Parent = joinedRoot

	if joinedRoot.balanceFactor >= -1 {
		return &AvlTree[T]{joinedRoot}
	}
	return &AvlTree[T]{joinedRoot.rotateLeftUnmodifiedTree()}
}

func joinLeftAvl[T cmp.Ordered](tL, tR *AvlNode[T], joinValue T) *AvlTree[T] {
	pivotLeft, pivotRight, pivotValue := tR.Left, tR.Right, tR.Value

	// stopping condition - if height difference is 1 at most (right may be heigher than left)
	joinNodeBF := tL.height() - pivotLeft.height()
	if joinNodeBF >= -1 {
		joinNode := &AvlNode[T]{
			Value:         joinValue,
			Left:          tL,
			Right:         pivotLeft,
			balanceFactor: int8(joinNodeBF),
		}
		tL.Parent = joinNode
		if pivotLeft != nil {
			pivotLeft.Parent = joinNode
		}

		joinRootBF := joinNode.height() - pivotRight.height()
		if joinRootBF <= 1 {
			// height difference is 1 at most, complete the join
			joinRoot := &AvlNode[T]{
				Value:         pivotValue,
				Left:          joinNode,
				Right:         pivotRight,
				Parent:        nil,
				balanceFactor: int8(joinRootBF),
			}
			joinNode.Parent = joinRoot
			pivotRight.Parent = joinRoot
			return &AvlTree[T]{root: joinRoot}
		}

		joinNode = joinNode.rotateLeftUnmodifiedTree()

		// height difference is yet too great
		joinRoot := &AvlNode[T]{
			Value:         pivotValue,
			Left:          joinNode,
			Right:         pivotRight,
			Parent:        nil,
			balanceFactor: int8(joinNode.height()) - int8(pivotRight.height()),
		}
		joinRoot.Left.Parent = joinRoot
		pivotRight.Parent = joinRoot

		return &AvlTree[T]{joinRoot.rotateRightUnmodifiedTree()}
	}

	joinedLeftTree := joinLeftAvl(tL, pivotLeft, joinValue)
	joinedRoot := &AvlNode[T]{
		Value:         pivotValue,
		Left:          joinedLeftTree.root,
		Right:         pivotRight,
		balanceFactor: int8(joinedLeftTree.Height()) - int8(pivotRight.height()),
	}
	joinedRoot.Left.Parent = joinedRoot
	pivotRight.Parent = joinedRoot

	if joinedRoot.balanceFactor <= 1 {
		return &AvlTree[T]{joinedRoot}
	}
	return &AvlTree[T]{joinedRoot.rotateRightUnmodifiedTree()}
}

func (t *AvlTree[T]) AvlSplit(wedge T) (found bool, t1, t2 *AvlTree[T]) {
	return t.root.AvlSplit(wedge)
}

func (t *AvlNode[T]) AvlSplit(wedge T) (found bool, t1, t2 *AvlTree[T]) {
	if t == nil {
		return false, nil, nil
	}

	leftSubtree, rootValue, rightSubtree := t.Left, t.Value, t.Right
	if wedge < rootValue {
		b, lTag, rTag := leftSubtree.AvlSplit(wedge)
		if rightSubtree != nil {
			//rightSubtree.Parent = nil
		}
		//fmt.Println("Split-AvlJoing(rTag: ", &rTag, ", rightSubtree: ", &rightSubtree, ", rootValue: ", rootValue, ")")
		joined, joinedTree := AvlJoin(rTag, &AvlTree[T]{rightSubtree}, rootValue)
		fmt.Println("joined wedge < rootValue?: ", joined)
		return b, lTag, joinedTree
	}
	if wedge > rootValue {
		b, lTag, rTag := rightSubtree.AvlSplit(wedge)
		//if leftSubtree != nil {
		//	leftSubtree.Parent = nil
		//}
		//if leftSubtree.Parent != nil {
		//
		//}

		joined, joinedTree := AvlJoin(&AvlTree[T]{leftSubtree}, lTag, rootValue)
		fmt.Println("joined wedge > rootValue?: ", joined)
		return b, joinedTree, rTag
	}
	fmt.Println("wedge == rootValue")
	if leftSubtree != nil {
		leftSubtree.Parent = nil
	}
	if rightSubtree != nil {
		rightSubtree.Parent = nil
	}
	return true, &AvlTree[T]{leftSubtree}, &AvlTree[T]{rightSubtree}
}

func AvlUnion[T cmp.Ordered](t1, t2 *AvlTree[T]) *AvlTree[T] {
	//return avlUnion(t1.root, t2.root)

	if t1 == nil {
		return t2
	}
	if t1.root == nil {
		return t2
	}
	if t2 == nil {
		return t1
	}
	if t2.root == nil {
		return t2
	}
	_, tL, tR := t2.AvlSplit(t1.root.Value)
	if t1.root.Left != nil {
		t1.root.Left.Parent = nil
	}
	if t1.root.Right != nil {
		t1.root.Right.Parent = nil
	}
	_, union := AvlJoin(
		AvlUnion(&AvlTree[T]{t1.root.Left}, &AvlTree[T]{tL.root}),
		AvlUnion(&AvlTree[T]{t1.root.Right}, &AvlTree[T]{tR.root}),
		t1.root.Value)
	return union

}

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
		// deletion balancing case
		root.balanceFactor = -1
		pivot.balanceFactor = 1
	} else {
		// insertion or deletion balancing case
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

// Performs a left rotation on (root *AvlNode[T])'s subtree,
// which makes (pivot *AvlNode[T]) the new subtree root
func (root *AvlNode[T]) rotateLeftUnmodifiedTree() (pivot *AvlNode[T]) {
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

	root.balanceFactor = root.balanceFactor + 1 - int8(math.Min(float64(pivot.balanceFactor), 0))
	pivot.balanceFactor = pivot.balanceFactor + 1 + int8(math.Max(float64(root.balanceFactor), 0))

	root.Parent = pivot
	return
}

// Performs a right rotation on (root *AvlNode[T])'s subtree,
// which makes (pivot *AvlNode[T]) the new subtree root
func (root *AvlNode[T]) rotateRightUnmodifiedTree() (pivot *AvlNode[T]) {
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

	root.balanceFactor = root.balanceFactor - 1 - int8(math.Max(float64(pivot.balanceFactor), 0))
	pivot.balanceFactor = pivot.balanceFactor - 1 + int8(math.Min(float64(root.balanceFactor), 0))

	root.Parent = pivot
	return
}

func (t *AvlTree[T]) String() string {
	if t == nil {
		return "<nil>"
	}
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
