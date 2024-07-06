package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	head *ListItem // Указатель на первый элемент
	tail *ListItem // Указатель на последний элемент
	size int
}

func (l *list) Len() int {
	return l.size
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	oldFrontItem := l.Front()
	item := &ListItem{v, oldFrontItem, nil}

	if oldFrontItem != nil {
		oldFrontItem.Prev = item
	}

	l.head = item
	if l.Back() == nil {
		l.tail = item
	}
	l.size++
	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	oldTailItem := l.Back()
	item := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  oldTailItem,
	}
	if oldTailItem != nil {
		oldTailItem.Next = item
	}

	l.tail = item
	if l.Front() == nil {
		l.head = item
	}

	l.size++
	return item
}

func (l *list) Remove(i *ListItem) {
	nextItem := i.Next
	prevItem := i.Prev

	if prevItem != nil {
		prevItem.Next = nextItem
	}

	if nextItem != nil {
		nextItem.Prev = prevItem
	}

	if l.Front() == i && nextItem != nil {
		l.head = nextItem
	}

	if l.Back() == i && prevItem != nil {
		l.tail = prevItem
	}

	l.size--
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return new(list)
}
