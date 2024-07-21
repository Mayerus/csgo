package csgo

import "errors"

type Node[T comparable] struct {
	data T
	next *Node[T]
}

type LinkedList[T comparable] struct {
	head *Node[T]
}

// Add value at the end of the list
func (l *LinkedList[T]) Add(data T) {
	node := &Node[T]{data, nil}
	if l.head == nil {
		l.head = node
		return
	}
	var iterator *Node[T] = l.head
	for ; iterator.next != nil; iterator = iterator.next {
	}
	iterator.next = node
}

// Add value to the list at the specified index
func (l *LinkedList[T]) AddAt(index int, data T) error {
	if index == 0 {
		l.head = &Node[T]{data, l.head}
	}
	if l.head == nil {
		return errors.New("index out of range")
	}
	listIterator := l.head
	for i := 0; i < index-2; i++ {
		if listIterator.next == nil {
			return errors.New("index is out of range")
		}
		listIterator = listIterator.next
	}
	old_next := listIterator.next
	listIterator.next = &Node[T]{data, old_next}
	return nil
}

// Add value at the start of the list
func (l *LinkedList[T]) Insert(data T) {
	l.head = &Node[T]{data, l.head}
}

// Delete value at the end of the list
func (l *LinkedList[T]) Delete() error {
	if l.head == nil {
		return errors.New("list is empty")
	}
	node := l.head
	for ; node.next.next != nil; node = node.next {
	}
	node.next = nil
	return nil
}

func (l *LinkedList[T]) DeleteAt(index int) error {
	// TODO: function
	if l.head == nil {
		return errors.New("list is empty")
	}
	if index == 0 {
		l.head = l.head.next
		return nil
	}
	iterator := l.head
	for i := 0; i < index-1; i++ {
		if iterator.next == nil {
			return errors.New("index is out of range")
		}
		iterator = iterator.next
	}
	iterator.next = iterator.next.next
	return nil
}

func (l *LinkedList[T]) Get(index int) (T, error) {
	var value T
	if l.head == nil {
		return value, errors.New("list is empty")
	}
	listIterator := l.head
	for i := 0; i < index; i++ {
		if listIterator.next == nil {
			return value, errors.New("index is out of range")
		}
		listIterator = listIterator.next
	}
	return listIterator.data, nil
}

func (l *LinkedList[T]) Count(data T) int {
	node := l.head
	count := 0
	for ; node != nil; node = node.next {
		count++
	}
	return count
}
