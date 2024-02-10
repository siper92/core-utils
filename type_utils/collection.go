package type_utils

import (
	"slices"
	"sort"
)

const (
	Equal       int = 0
	LessThan    int = -1
	GreaterThan int = 1
)

type ItemWithName interface {
	Name() string
}

type ItemWithKey interface {
	KeyEqual(other any) bool
	Key() string
}

type ComparableItem interface {
	Compare(other any) int
}

type CollectionItem interface {
	ItemWithKey
	ComparableItem
}

type ICollection[T CollectionItem] interface {
	sort.Interface
	Add(item T)
	Get(k string) T
	Remove(k string)
	Contains(i T) bool
	ContainsKey(k string) bool
	Sort() ICollection[T]
	Slice() []T
	Map() map[string]T
}

var _ ICollection[CollectionItem] = (*Collection[CollectionItem])(nil)

type Collection[T CollectionItem] struct {
	items []T
}

func NewCollection[T CollectionItem](items ...T) *Collection[T] {
	if len(items) == 0 {
		return &Collection[T]{items: make([]T, 0)}
	}

	return &Collection[T]{items: items}
}

func (c *Collection[T]) Slice() []T {
	return c.items
}

func (c *Collection[T]) Map() map[string]T {
	m := make(map[string]T)
	for _, item := range c.items {
		m[item.Key()] = item
	}

	return m
}

func (c *Collection[T]) Add(item T) {
	c.items = append(c.items, item)
}

func (c *Collection[T]) Remove(key string) {
	for i, val := range c.items {
		if val.KeyEqual(key) {
			c.items = append(c.items[:i], c.items[i+1:]...)
			break
		}
	}
}

func (c *Collection[T]) Get(key string) (noVal T) {
	for _, item := range c.items {
		if item.KeyEqual(key) {
			return item
		}
	}

	return noVal
}

func (c *Collection[T]) Sort() ICollection[T] {
	slices.SortFunc(c.items, func(i, j T) int {
		return i.Compare(j)
	})

	return c
}

func (c *Collection[T]) Len() int { return len(c.items) }
func (c *Collection[T]) Less(i, j int) bool {
	a := c.items[i]
	b := c.items[j]

	return a.Compare(b) == LessThan
}
func (c *Collection[T]) Swap(i, j int) { c.items[i], c.items[j] = c.items[j], c.items[i] }

func (c *Collection[T]) Contains(i T) bool {
	for _, item := range c.items {
		if item.KeyEqual(i) {
			return true
		}
	}

	return false
}

func (c *Collection[T]) ContainsKey(n string) bool {
	for _, item := range c.items {
		if item.KeyEqual(n) {
			return true
		}
	}

	return false
}
