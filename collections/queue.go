package csgo

import "errors"

type Queuer[T any] interface {
	Push(T)
	Pop() (T, error)
	Count() int
}

// Double-ended queue, (incomplete)
type deque[T comparable] []T

type Queue[T comparable] []T

func (q *Queue[T]) Push(item T) {
	*q = append(*q, item)
}

func (q *Queue[T]) Pop() (T, error) {
	var item T
	if len(*q) == 0 {
		return item, errors.New("Queue is empty")
	}
	item = (*q)[0]
	*q = (*q)[1:]
	return item, nil
}

func (q *Queue[T]) Count() int {
	return len(*q)
}
