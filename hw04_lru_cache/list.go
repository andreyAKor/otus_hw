package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	Len() int                      // длина списка
	Front() *Item                  // первый Item
	Back() *Item                   // последний Item
	PushFront(v interface{}) *Item // добавить значение в начало
	PushBack(v interface{}) *Item  // добавить значение в конец
	Remove(i *Item)                // удалить элемент
	MoveToFront(i *Item)           // переместить элемент в начало
}

type Item struct {
	Value interface{} // значение
	Next  *Item       // следующий элемент
	Prev  *Item       // предыдущий элемент
}

type list struct {
	front *Item // Первый элемент списка
	back  *Item // Последний элемент списка
	size  int   // Размер списка
}

// Длина списка
func (l *list) Len() int {
	return l.size
}

// Первый элемент списка
func (l *list) Front() *Item {
	return l.front
}

// Последний элемент списка
func (l *list) Back() *Item {
	return l.back
}

// Добавляет значение в начало списка
func (l *list) PushFront(v interface{}) *Item {
	if l.front == nil {
		l.front = &Item{
			Value: v,
		}
		l.back = l.front
	} else {
		l.front.Next = &Item{
			Value: v,
			Prev:  l.front,
		}

		l.front = l.front.Next
	}

	l.size++

	return l.front
}

// Добавляет значение в конец списка
func (l *list) PushBack(v interface{}) *Item {
	if l.back == nil {
		l.back = &Item{
			Value: v,
		}
		l.front = l.back
	} else {
		l.back.Prev = &Item{
			Value: v,
			Next:  l.back,
		}

		l.back = l.back.Prev
	}

	l.size++

	return l.back
}

// Удаляет элемент из списка
func (l *list) Remove(i *Item) {
	if i == nil {
		return
	}

	l.remove(i)
	l.size--
}

func (l *list) remove(i *Item) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.back = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.front = i.Prev
	}
}

// переместить элемент в начало
func (l *list) MoveToFront(i *Item) {
	if i == nil || i == l.front {
		return
	}

	l.remove(i)

	l.front.Next = i
	i.Prev = l.front
	i.Next = nil
	l.front = i
}

func NewList() List {
	return &list{}
}
