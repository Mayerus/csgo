package collections

import (
	"testing"

	"golang.org/x/exp/rand"
)

func TestAvlTree(t *testing.T) {
	avl := AvlTree[int]{}

	list := LinkedList[int]{}

	k := 1
	l := 1000
	m := 1000
	for i := 0; i < k; i++ {
		for j := 0; j < l; j++ {
			value := rand.Intn(m) - m/2
			avl.Insert(value)
			list.Add(value)
		}
	}

	if !avl.containsAllElements(&list) {
		t.Fatalf("Some values are missing from the AVL tree: \n\n%v\n\nInsert order:%v\n\n\n", avl.String(), list.String())
	}

	if avl.hasDuplicateValues() {
		t.Fatalf("AVL tree has duplicate values: \n\n%v\n\nInsert order:\n\t%v\n\n\n", avl.String(), list.String())
	}

	if !avl.root.isOrdered() {
		t.Fatalf("AVL tree is unordered: \n\n%v\n\nInsert order:\n\t%v\n\n\n", avl.String(), list.String())
	}

	if !avl.root.propperDynasty() {
		t.Fatalf("AVL tree has unpropper dynasty: \n\n%v\n\nInsert order:%v\n\n\n", avl.String(), list.String())
	}

	if !avl.root.isBalanced() {
		t.Fatalf("AVL tree is unbalanced: \n\n%v\n\nInsert order:\n\t%v\n\n\n", avl.String(), list.String())
	}
}

func (t *AvlTree[T]) hasDuplicateValues() bool {
	values := BSTree[T]{}

	return t.root.hasDuplicateValues(&values)
}

func (n *AvlNode[T]) hasDuplicateValues(values *BSTree[T]) bool {
	if n == nil {
		return false
	}
	if values.Search(n.Value) != nil {
		return true
	}
	values.Insert(n.Value)
	return n.Left.hasDuplicateValues(values) || n.Right.hasDuplicateValues(values)
}

func (n *AvlNode[T]) isBalanced() bool {
	if n == nil {
		return true
	}
	if n.Left == nil && n.Right == nil {
		return true
	}
	lHeight := n.Left.getTreeHeight(0)
	rHeight := n.Right.getTreeHeight(0)
	balanceFactor := lHeight - rHeight
	if balanceFactor < -2 || balanceFactor > 2 {
		return false
	}
	return n.Left.isBalanced() && n.Right.isBalanced()
}

func (n *AvlNode[T]) getTreeHeight(height int) int {
	if n == nil {
		return height
	}
	height++
	if n.Left == nil && n.Right == nil {
		return height
	}
	lHeight := n.Left.getTreeHeight(height)
	rHeight := n.Right.getTreeHeight(height)
	if lHeight > rHeight {
		return lHeight
	}
	return rHeight
}

func (n *AvlNode[T]) isOrdered() bool {
	if n.Left != nil {
		if n.Value < n.Left.Value {
			return false
		}

		if !n.Left.isOrdered() {
			return false
		}
	}

	if n.Right != nil {
		if n.Value > n.Right.Value {
			return false
		}
		if !n.Right.isOrdered() {
			return false
		}
	}

	return true
}

func (n *AvlNode[T]) propperDynasty() bool {
	if n.Left != nil {
		if n.Left.Parent != n {
			return false
		}
		if !n.Left.propperDynasty() {
			return false
		}
	}

	if n.Right != nil {
		if n.Right.Parent != n {
			return false
		}
		if !n.Right.propperDynasty() {
			return false
		}
	}

	return true
}

func (t *AvlTree[T]) containsAllElements(list *LinkedList[T]) bool {
	for node := list.head; node != nil; node = node.next {
		if t.Search(node.data) == nil {
			return false
		}
	}
	return true
}
