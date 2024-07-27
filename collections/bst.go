package collections

// Binary Search Tree
type BSTree[T Numeric] struct {
	root *TreeNode[T]
}

func (t *BSTree[T]) String() string {
	return t.root.String()
}

func (t *BSTree[T]) Insert(value T) {
	var iteratorParent *TreeNode[T]
	iterator := t.root
	for iterator != nil {
		iteratorParent = iterator
		if value < iterator.Value {
			iterator = iterator.Left
			continue
		}
		iterator = iterator.Right
	}
	node := &TreeNode[T]{value, nil, nil, iteratorParent}
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

func (t *BSTree[T]) Delete(value T) bool {
	deleted := t.Search(value)
	//defer deleted.Dispose()
	if deleted.Left == nil {
		t.shiftNodes(deleted, deleted.Right)
		return true
	}
	if deleted.Right == nil {
		t.shiftNodes(deleted, deleted.Left)
		return true
	}
	successor := deleted.Successor()
	if successor.Parent != deleted {
		t.shiftNodes(successor, successor.Right)
		successor.Right = deleted.Right
		successor.Right.Parent = successor
		return true
	}
	t.shiftNodes(deleted, successor)
	successor.Left = deleted.Left
	successor.Left.Parent = successor
	return true
}

func (t *BSTree[T]) shiftNodes(original, successor *TreeNode[T]) {
	if original.Parent == nil {
		t.root = successor
		//return
	} else if original == original.Parent.Left {
		// if a left-child is replaced
		original.Parent.Left = successor
		//return
	} else {
		// if a right-child is replaced
		original.Parent.Right = successor
	}
	if successor != nil {
		successor.Parent = original.Parent
	}
}

func (n *TreeNode[T]) Min() (node *TreeNode[T]) {
	node = n
	for node.Left != nil {
		node = node.Left
	}
	return
}

func (n *TreeNode[T]) Max() (node *TreeNode[T]) {
	node = n
	for node.Right != nil {
		node = node.Right
	}
	return
}

// The successor of node n is the node with the smallest value greater than n's.
func (n *TreeNode[T]) Successor() *TreeNode[T] {
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

// The predecessor of node n is the node with the largest value smaller than n's.
func (n *TreeNode[T]) Predecessor() *TreeNode[T] {
	if n.Left == nil {
		return n.Left.Max()
	}
	node := n
	predecessor := n.Parent
	for predecessor != nil && node == n.Left {
		node = predecessor
		predecessor = predecessor.Parent
	}
	return predecessor
}

// Iterative binary tree search (Considered more efficient on most machines)
func (t *BSTree[T]) Search(value T) (node *TreeNode[T]) {
	return t.root.Search(value)
}

func (t *TreeNode[T]) Search(value T) (node *TreeNode[T]) {
	node = t
	for node != nil && value != node.Value {
		if value < node.Value {
			node = node.Left
			continue
		}
		node = node.Right
	}
	return
}

func (t *BSTree[T]) ConSearch(value T, ch chan *TreeNode[T]) {
	conSearch(t.root, value, ch)
	close(ch)
}

func conSearch[T Numeric](n *TreeNode[T], value T, ch chan *TreeNode[T]) {
	if n == nil {
		return
	}
	if n.Value == value {
		ch <- n
		return
	}
	conSearch(n.Left, value, ch)
	conSearch(n.Right, value, ch)
}

func (t *BSTree[T]) InorderTraversal() {
	t.root.InorderTraversal()
}

func (t *BSTree[T]) PreorderTraversal() {
	t.root.PreorderTraversal()
}

func (t *BSTree[T]) PostorderTraversal() {
	t.root.PostorderTraversal()
}
