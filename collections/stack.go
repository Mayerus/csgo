package csgo

import "errors"

type Stack[T comparable] []T

func (s *Stack[T]) Push(item T) {
	*s = append(*s, item)
}

func (s *Stack[T]) Pop() (T, error) {
	if len(*s) == 0 {
		var v T
		return v, errors.New("Stack is empty")
	}

	lastIndex := len(*s) - 1
	popped := (*s)[lastIndex]
	*s = (*s)[:lastIndex]
	return popped, nil
}

func (s *Stack[T]) Peek() (T, error) {
	if len(*s) == 0 {
		var v T
		return v, errors.New("Stack is empty")
	}
	return (*s)[s.Count()-1], nil
}

func (s *Stack[T]) Count() int {
	return len(*s)
}
