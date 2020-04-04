package hw04_lru_cache //nolint:golint,stylecheck

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
}

type listItem struct {
	Value interface{} // значение
	Next  *listItem   // следующий элемент
	Prev  *listItem   // предыдущий элемент
}

type list struct {
	front  *listItem
	back   *listItem
	length int
}

// Len Вернёт длину списка
func (l *list) Len() int {
	return l.length
}

// Front Вернёт первый элемент списка
func (l *list) Front() *listItem {
	return l.front
}

// Back Вернёт крайний элемент списка, или первый, если в списке один элемент
func (l *list) Back() *listItem {
	return l.back
}

// PushFront Добавит значение в начало списка
func (l *list) PushFront(v interface{}) *listItem {
	i := &listItem{
		Value: v,
	}

	// nil <- (next) front <-> ... <-> elem <-> ... <-> back (prev) -> nil
	t := l.front

	i.Next = nil
	i.Prev = t

	if t != nil {
		t.Next = i
	}

	l.front = i

	l.length++

	if l.length == 1 {
		l.back = l.front
	}

	return i
}

// PushBack Добавит значение в конец списка
func (l *list) PushBack(v interface{}) *listItem {
	i := &listItem{
		Value: v,
	}

	// nil <- (next) front <-> ... <-> elem <-> ... <-> back (prev) -> nil
	t := l.back

	i.Prev = nil
	i.Next = t

	if t != nil {
		t.Prev = i
	}

	l.back = i

	l.length++

	if l.length == 1 {
		l.front = l.back
	}

	return i
}

// Remove Удалит элемент из списка
func (l *list) Remove(i *listItem) {
	if i == nil {
		return
	}

	// nil <- (next) front <-> ... <-> elem <-> ... <-> back (prev) -> nil
	prev := i.Prev
	next := i.Next

	if next != nil {
		next.Prev = prev
	}
	if prev != nil {
		prev.Next = next
	}

	switch {
	case i == l.back:
		l.back = i.Next
	case i == l.front:
		l.front = i.Prev
	}

	l.length--
}

// Переместит элемент в начало списка
func (l *list) MoveToFront(i *listItem) {
	if i == nil {
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
		defer close(c)
		if l.back != nil {
			for i := l.Back(); i != nil; i = i.Next {
				c <- i
			}
		}
	}()

	return c
}

// NewList Создаст новый список
func NewList() List {
	return &list{}
}
