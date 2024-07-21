package csgo

import "golang.org/x/exp/constraints"

type Numeric interface {
	constraints.Float | constraints.Integer
}

type BSTree[T Numeric] struct {
	Value       T
	Left, Right *BSTree[T]
}

func (t *BSTree[T]) Insert(value T) {
	if value < t.Value {
		if t.Left == nil {
			t.Left = &BSTree[T]{value, nil, nil}
			return
		}
		t.Left.Insert(value)
	}
	if value > t.Value {
		if t.Right == nil {
			t.Right = &BSTree[T]{value, nil, nil}
			return
		}
		t.Right.Insert(value)
	}

}

func (t *BSTree[T]) Delete(value T) bool {
	if t.Value == value {
		// TODO: delete root
	}
	if t.Left != nil {
		if t.Left.Value == value {
			t.deleteLeft()
		}
	}
	if t.Right != nil {
		if t.Right.Value == value {
			t.deleteRight()
		}
	}

	var deleted bool
	if value < t.Value {
		if t.Left == nil {
			return false
		}
		deleted = t.Left.Delete(value)
	}
	if value > t.Value {
		if t.Right == nil {
			return false
		}
		deleted = t.Right.Delete(value)
	}
	return deleted
}

func (t *BSTree[T]) deleteLeft() bool {
	if t.Left.isLeaf() {
		t.Left = nil
		return true
	}
	deletedNode := t.Left
	if t.Left.Left != nil && t.Left.Right == nil {
		t.Left = t.Left.Left
		deletedNode.Left = nil
		return true
	}
	if t.Left.Left == nil && t.Left.Right != nil {
		t.Left = t.Left.Right
		deletedNode.Right = nil
		return true
	}
	// root has 2 children
	t.rebaseMin()
	return true
}

func (t *BSTree[T]) deleteRight() bool {
	if t.Right.isLeaf() {
		t.Right = nil
		return true
	}
	deletedNode := t.Right
	if t.Right.Left != nil && t.Right.Right == nil {
		t.Right = t.Right.Left
		deletedNode.Left = nil
		return true
	}
	if t.Right.Left == nil && t.Right.Right != nil {
		t.Right = t.Right.Right
		deletedNode.Right = nil
		return true
	}
	// root has 2 children
	t.rebaseMin()
	return true
}

func (t *BSTree[T]) rebaseMin() {
	if t.Right.Left == nil {
		oldRight := t.Right
		t.Value = oldRight.Value
		t.Right = oldRight.Right
		oldRight.Right = nil
		return
	}
	min := t.Right.deleteMin()
	t.Value = min
}

func (t *BSTree[T]) rebaseMax() {
	if t.Left.Right == nil {
		oldLeft := t.Left
		t.Value = oldLeft.Value
		t.Left = oldLeft.Left
		oldLeft.Left = nil
		return
	}
	max := t.Left.deleteMax()
	t.Value = max
}

func (t *BSTree[T]) deleteMin() T {
	if t.Left.Left == nil {
		min := t.Left.Value
		// reatach Min's right
		if t.Left.Right != nil {
			t.Left = t.Left.Right
			return min
		}
		t.Left = nil
		return min
	}
	return t.Left.deleteMin()
}

func (t *BSTree[T]) deleteMax() T {
	if t.Right.Right == nil {
		max := t.Right.Value
		// reatach Max's left
		if t.Right.Left != nil {
			t.Right = t.Right.Left
			return max
		}
		t.Right = nil
		return max
	}
	return t.Right.deleteMax()
}

func (t *BSTree[T]) isLeaf() bool {
	return t.Left == nil && t.Right == nil
}

func (t *BSTree[T]) Search(value T) bool {
	if value == t.Value {
		return true
	}

	var leftSearch, rightSearch bool
	if value < t.Value {
		if t.Left == nil {
			return false
		}
		leftSearch = t.Left.Search(value)
	}
	if value > t.Value {
		if t.Right == nil {
			return false
		}
		rightSearch = t.Right.Search(value)
	}

	return leftSearch || rightSearch
}
