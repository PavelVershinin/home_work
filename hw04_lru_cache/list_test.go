package hw04_lru_cache //nolint:golint,stylecheck

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, l.Len(), 0)
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, l.Len(), 3)

		// Ну никак не может у крайнего элемента Back() быть следующий элемент Next, да ещё и ссылаться в середину списка
		// Тут совершенно очевидно была опечатка.
		// Поэтому Back заметил на Front
		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, l.Len(), 2)

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, l.Len(), 7)
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		// Тут тоже, берётся крайний элемент списка Back() и от него, пытается получить следующий элемент Next
		// Уверен, что это опечатка и тут предполагалось движение по списку от крайнего элемента к первому
		for i := l.Back(); i != nil; i = i.Prev {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{50, 30, 10, 40, 60, 80, 70}, elems)
	})

	t.Run("fetch", func(t *testing.T) {
		expected := []int{50, 30, 10, 40, 60, 80, 70}
		l := NewList()
		for _, i := range expected {
			l.PushBack(i)
		}
		elems := make([]int, 0, l.Len())
		for i := range l.Fetch() {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, expected, elems)
	})

	t.Run("exists", func(t *testing.T) {
		l1, l2 := NewList(), NewList()
		li1, li2 := l1.PushFront(nil), l2.PushFront(nil)

		require.True(t, l1.Exists(li1))
		require.True(t, l2.Exists(li2))
		require.False(t, l1.Exists(li2))
		require.False(t, l2.Exists(li1))
	})

	t.Run("last", func(t *testing.T) {
		l := NewList()
		l.PushFront(0)
		l.PushFront(1)
		require.Equal(t, l.Back().Value, 0)
	})
}
