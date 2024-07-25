package collections

type AvlTree[T Numeric] struct {
	Root *TreeNode[T]
}

func (t *AvlTree[T]) Search(value T) {
	t.Root.Search(value)
}

func (t *AvlTree[T]) Delete() {

}

func (t *AvlTree[T]) Insert() {

}

func (t *TreeNode[T]) rotateLeft() {
	if t.Left == nil {
		return
	}
	// root is the initial parent and pivot is the child to take the root's place
	pivot := t.Right
	root := t

	// reposition pivot's right child as root's left child
	root.Right = pivot.Left
	root.Right.Parent = root

	// set root as pivot's right child
	pivot.Left = root
	root.Parent = pivot
	pivot.Parent = root.Parent
}

func (t *TreeNode[T]) rotateRight() {
	if t.Left == nil {
		return
	}
	// root is the initial parent and pivot is the child to take the root's place
	pivot := t.Left
	root := t

	// reposition pivot's right child as root's left child
	root.Left = pivot.Right
	root.Left.Parent = root

	// set root as pivot's right child
	pivot.Right = root
	root.Parent = pivot
	pivot.Parent = root.Parent
}
