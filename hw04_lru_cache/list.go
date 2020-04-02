package hw04_lru_cache //nolint:golint,stylecheck

import (
	"sync"
	"unsafe"
)

// List Интерфейс двухсвязного списка
type List interface {
	Len() int                          // длина списка
	Front() *listItem                  // первый Item
	Back() *listItem                   // последний Item
	PushFront(v interface{}) *listItem // добавить значение в начало
	PushBack(v interface{}) *listItem  // добавить значение в конец
	Remove(i *listItem)                // удалить элемент
	MoveToFront(i *listItem)           // переместить элемент в начало
	Fetch() <-chan *listItem
	Exists(i *listItem) bool
}

type listItem struct {
	Value interface{} // значение
	Next  *listItem   // следующий элемент
	Prev  *listItem   // предыдущий элемент
}

type list struct {
	mu    sync.Mutex
	first *listItem
	last  *listItem
	ptrs  map[unsafe.Pointer]struct{}
}

// Len Вернёт длину списка
func (l *list) Len() int {
	l.mu.Lock()
	length := len(l.ptrs)
	l.mu.Unlock()
	return length
}

// Front Вернёт первый элемент списка
func (l *list) Front() *listItem {
	l.mu.Lock()
	first := l.first
	l.mu.Unlock()
	return first
}

// Back Вернёт крайний элемент списка, или первый, если в списке один элемент
func (l *list) Back() *listItem {
	l.mu.Lock()
	var last *listItem
	if l.last != nil {
		last = l.last
	} else {
		last = l.first
	}
	l.mu.Unlock()
	return last
}

// PushFront Добавит значение в начало списка
func (l *list) PushFront(v interface{}) *listItem {
	l.mu.Lock()
	item := &listItem{
		Value: v,
	}

	p := l.first
	l.first = item
	item.Next = p
	if p != nil {
		p.Prev = l.first
	}

	// Если в списке ещё нет крайнего элемента, то первый и будет крайним
	if l.last == nil {
		l.last = item
	}

	l.ptrs[unsafe.Pointer(item)] = struct{}{}

	l.mu.Unlock()

	return item
}

// PushBack Добавит значение в конец списка
func (l *list) PushBack(v interface{}) *listItem {
	l.mu.Lock()
	item := &listItem{
		Value: v,
	}

	var p *listItem
	switch {
	case l.last != nil:
		p = l.last
	case l.first != nil:
		p = l.first
	default:
		l.mu.Unlock()
		return l.PushFront(v)
	}

	l.last = item
	item.Prev = p
	p.Next = l.last

	l.ptrs[unsafe.Pointer(item)] = struct{}{}

	l.mu.Unlock()

	return item
}

// Remove Удалит элемент из списка
func (l *list) Remove(i *listItem) {
	if i == nil || !l.Exists(i) {
		return
	}

	l.mu.Lock()

	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	switch {
	case i == l.first:
		l.first = i.Next
	case i == l.last:
		l.last = i.Prev
	}

	delete(l.ptrs, unsafe.Pointer(i))

	l.mu.Unlock()
}

// Переместит элемент в начало списка
func (l *list) MoveToFront(i *listItem) {
	if i == nil || !l.Exists(i) {
		return
	}
	l.Remove(i)
	l.PushFront(i.Value)
}

// Fetch Переберёт все элементы списка по порядку, для конструкций вроде for item := range list.Fetch() {...
// Just for fun ;)
func (l *list) Fetch() <-chan *listItem {
	var c = make(chan *listItem)

	go func() {
		if l.first != nil {
			for i := l.first; i != nil; i = i.Next {
				c <- i
			}
		}
		close(c)
	}()

	return c
}

// Exists Вернёт true, если элемент принадлежит этому списку
func (l *list) Exists(i *listItem) bool {
	l.mu.Lock()
	_, ok := l.ptrs[unsafe.Pointer(i)]
	l.mu.Unlock()
	return ok
}

// NewList Создаст новый список
func NewList() List {
	return &list{
		ptrs: make(map[unsafe.Pointer]struct{}),
	}
}
