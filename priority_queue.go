package main

type Sortable interface {
	ComesBefore(other Sortable) bool
}

type PriorityQueue[T Sortable] interface {
	Push(other T)
	Pop() T
	Len() uint
	PopAt(idx uint)
	LookAt(idx uint)
}

type LinkedNode struct {
	Prev *LinkedNode
	Next *LinkedNode
	Val  Sortable
}

type Queue struct {
}
