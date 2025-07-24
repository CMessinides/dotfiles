package main

type stack[T any] struct {
	slice []T
}

func newStack[T any]() *stack[T] {
	return &stack[T]{
		slice: make([]T, 0),
	}
}

func (s *stack[T]) Len() int {
	return len(s.slice)
}

func (s *stack[T]) Top() T {
	return s.slice[len(s.slice)-1]
}

func (s *stack[T]) Push(value T) int {
	s.slice = append(s.slice, value)
	return s.Len()
}

func (s *stack[T]) Pop() (value T, ok bool) {
	if len(s.slice) == 0 {
		return *new(T), false
	}

	value = s.slice[len(s.slice)-1]
	s.slice = s.slice[:len(s.slice)-1]
	return value, true
}
