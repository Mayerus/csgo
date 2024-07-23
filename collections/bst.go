package collections

// Binary Search Tree
type BSTree[T Numeric] struct {
	Root *TreeNode[T]
}

type searchResult[T Numeric] struct {
	Node, Parent *TreeNode[T]
}

func (t *BSTree[T]) String() string {
	return t.Root.String()
}

func (t *BSTree[T]) Insert(value T) {
	var iteratorParent *TreeNode[T]
	iterator := t.Root
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
		t.Root = node
		return
	}
	if value < iterator.Parent.Value {
		iteratorParent.Left = node
	}
	iteratorParent.Right = node
}

func (t *BSTree[T]) Delete(value T) bool {
	deleted := t.Search(value)
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
		t.shiftNodes(successor, successor.Parent)
		successor.Right = deleted.Right
		successor.Right.Parent = successor
	}
	t.shiftNodes(successor, successor.Parent)
	successor.Left = deleted.Left
	successor.Left.Parent = successor
	return true
}

func (t *BSTree[T]) shiftNodes(original, successor *TreeNode[T]) {
	defer original.Dispose()
	if successor != nil {
		successor.Parent = original.Parent
	}
	if original.Parent == nil {
		t.Root = successor
		return
	}
	if original == original.Parent.Left {
		// if a left-child is replaced
		original.Parent.Left = successor
		return
	}
	// if a right-child is replaced
	original.Parent.Right = successor

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
	node = t.Root
	for node != nil && value != node.Value {
		if value == node.Value {
			node = node.Left
			continue
		}
		node = node.Right
	}
	return
}

func (t *BSTree[T]) ConSearch(value T, ch chan *TreeNode[T]) {
	conSearch(t.Root, value, ch)
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
	t.Root.InorderTraversal()
}

func (t *BSTree[T]) PreorderTraversal() {
	t.Root.PreorderTraversal()
}

func (t *BSTree[T]) PostorderTraversal() {
	t.Root.PostorderTraversal()
}
