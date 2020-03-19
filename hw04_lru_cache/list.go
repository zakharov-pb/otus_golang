package hw04_lru_cache //nolint:golint,stylecheck

import (
	"fmt"
	"strings"
)

type Item struct {
	Value interface{}
	Next  *Item
	Prev  *Item

	parent *list
}

type List interface {
	Len() int
	Front() *Item
	Back() *Item
	PushFront(interface{}) *Item
	PushBack(interface{}) *Item
	Remove(*Item)
	MoveToFront(*Item)
	Clear()
}

type list struct {
	count int
	front *Item
	back  *Item
}

func NewList() List {
	return &list{0, nil, nil}
}

func (l *list) Len() int {
	return l.count
}

func (l *list) Front() *Item {
	return l.front
}

func (l *list) Back() *Item {
	return l.back
}

func (l *list) PushFront(value interface{}) *Item {
	i := &Item{Value: value, Next: nil, Prev: l.front, parent: l}
	if l.front == nil {
		l.back = i
	} else {
		l.front.Next = i
	}
	l.front = i
	l.count++
	return i
}

func (l *list) PushBack(value interface{}) *Item {
	i := &Item{Value: value, Next: l.back, Prev: nil, parent: l}
	if l.back == nil {
		l.front = i
	} else {
		l.back.Prev = i
	}
	l.back = i
	l.count++
	return i
}

func (l *list) checkItem(i *Item) bool {
	if (i.parent != l) || (i.Next != nil && i.Next.Prev != i) ||
		(i.Prev != nil && i.Prev.Next != i) {
		return false
	}
	return true
}

func (l *list) Remove(i *Item) {
	if i == nil || l.count == 0 || !l.checkItem(i) {
		return
	}
	next := i.Next
	prev := i.Prev
	if next != nil {
		next.Prev = prev
	}
	if prev != nil {
		prev.Next = next
	}
	if i == l.front {
		l.front = i.Prev
	}
	if i == l.back {
		l.back = i.Next
	}
	i.Next = nil
	i.Prev = nil
	l.count--
}

func (l *list) MoveToFront(i *Item) {
	if i == nil || i == l.front || !l.checkItem(i) {
		return
	}
	l.Remove(i)
	l.PushFront(i.Value)
}

func (l *list) String() string {
	var s strings.Builder
	itm := l.front
	s.WriteRune('[')
	if itm != nil {
		s.WriteString(fmt.Sprintf("%v", itm.Value))
		itm = itm.Prev
		for itm != nil {
			s.WriteString(fmt.Sprintf(" %v", itm.Value))
			itm = itm.Prev
		}
	}
	s.WriteRune(']')
	return s.String()
}

func (l *list) Clear() {
	l.back = nil
	l.front = nil
	l.count = 0
}
