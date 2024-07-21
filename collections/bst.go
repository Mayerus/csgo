package collections

import "io"

// Binary Search Tree
type BSTree[T Numeric] struct {
	Root *TreeNode[T]
}

func (t *BSTree[T]) String() string {
	return t.Root.String()
}

func (t *BSTree[T]) Print(w io.Writer) {
	t.Root.Print(w, 0, 'M')
}

func (t *BSTree[T]) Insert(value T) {
	t.Root = t.Root.insert(value)
}

func (n *TreeNode[T]) insert(value T) *TreeNode[T] {
	if n == nil {
		return &TreeNode[T]{Value: value}
	}
	if value < n.Value {
		if n.Left == nil {
			n.Left = &TreeNode[T]{value, nil, nil}
			return n
		}
		n.Left.insert(value)
	}
	if value > n.Value {
		if n.Right == nil {
			n.Right = &TreeNode[T]{value, nil, nil}
			return n
		}
		n.Right.insert(value)
	}
	return n
}

func (t *BSTree[T]) Delete(value T) bool {
	node, parent := t.Search(value)
	// node does not exist
	if node == nil {
		return false
	}

	//node has 2 children (remove right subtree minLeaf then )
	if node.HasLeft() && node.HasRight() {
		t.deleteBiChild(node)
		return true
	}

	defer node.Dispose()
	// node has only left child
	if node.HasLeft() {
		t.deleteLeftChild(node, parent)
		return true
	}
	// node has only right child
	if node.HasRight() {
		t.deleteRightChild(node, parent)
		return true
	}
	// node is leaf
	t.deleteLeaf(node, parent)
	return true
}

func (t *BSTree[T]) deleteBiChild(node *TreeNode[T]) {
	minLeaf, minParent := node.Right.Min()
	// remove minLeaf which must be a left child
	if minParent != nil {
		minParent.Left = nil
	}
	// By overriding the node value using minLeaf.Value
	//we perform an action which is equivalent to the nodes deletion.
	node.Value = minLeaf.Value
}

func (t *BSTree[T]) deleteLeftChild(node, parent *TreeNode[T]) {
	// delete root node with left child
	if parent == nil {
		t.Root = node.Left
		return
	}
	if parent.Left == node {
		parent.Left = node.Left
	}
	parent.Right = node.Left
}

func (t *BSTree[T]) deleteRightChild(node, parent *TreeNode[T]) {
	// delete root node with left child
	if parent == nil {
		t.Root = node.Right
		return
	}
	if parent.Left == node {
		parent.Left = node.Right
	}
	parent.Right = node.Right
}

func (t *BSTree[T]) deleteLeaf(node, parent *TreeNode[T]) {
	// node is leaf
	// delete root leaf
	if parent == nil {
		t.Root = nil
		return
	}
	// delete non root leaf
	if parent.Left == node {
		parent.Left = nil
	}
	parent.Right = nil
}

func (n *TreeNode[T]) Min() (min, parent *TreeNode[T]) {
	return n.min(parent), parent
}

func (n *TreeNode[T]) min(parent *TreeNode[T]) *TreeNode[T] {
	if n.Left != nil {
		parent = n
		n.min(parent)
	}
	return n
}

func (n *TreeNode[T]) Max() *TreeNode[T] {
	if n == nil {
		return nil
	}
	return n.Right.Max()
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

func (t *BSTree[T]) Search(value T) (node, parent *TreeNode[T]) {
	node, parent = t.Root.search(value, parent)
	return
}

func (n *TreeNode[T]) search(value T, recurse_parent *TreeNode[T]) (node, parent *TreeNode[T]) {
	if n == nil {
		return nil, nil
	}

	if value < n.Value {
		n.Left.search(value, n)
	}
	if value > n.Value {
		n.Right.search(value, n)
	}
	return n, recurse_parent
}
