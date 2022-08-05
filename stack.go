package rulex

import (
	"container/list"
)

type Stack struct {
	ll *list.List
}

func NewStack() *Stack {
	return &Stack{
		ll: list.New(),
	}
}

func (s Stack) Size() int {
	return s.ll.Len()
}

func (s Stack) Top() interface{} {
	return s.ll.Back().Value
}

func (s *Stack) Push(v interface{}) {
	s.ll.PushBack(v)
}

func (s *Stack) Empty() bool {
	return s.ll.Len() == 0
}

func (s *Stack) Pop() interface{} {
	elem := s.ll.Back()
	s.ll.Remove(elem)

	return elem.Value
}
