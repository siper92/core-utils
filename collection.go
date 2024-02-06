package core_utils

import (
	"slices"
	"sort"
)

type Comparison int

const (
	Equal          Comparison = 0
	LessThan                  = -1
	GreaterThan               = 1
	DifferentTypes            = 2
)

type Comparable interface {
	Compare(other Comparable) Comparison
}

type CollectionItem interface {
	Comparable
	Key() string
}

type Collection[T CollectionItem] interface {
	sort.Interface
	Items() []T
	Contains(i T) bool
	ContainsKey(k string) bool
	Remove(k string)
	Add(item T)
	Get(k string) T
	Sort()
}

var _ Collection[CollectionItem] = (*NamedTypeCollection[CollectionItem])(nil)

type NamedTypeCollection[T CollectionItem] struct {
	items []T
}

func NewCollection[T CollectionItem](items ...T) *NamedTypeCollection[T] {
	if len(items) == 0 {
		return &NamedTypeCollection[T]{items: make([]T, 0)}
	}

	return &NamedTypeCollection[T]{items: items}
}

func (c *NamedTypeCollection[T]) Items() []T {
	return c.items
}

func (c *NamedTypeCollection[T]) Add(item T) {
	c.items = append(c.items, item)
}

func (c *NamedTypeCollection[T]) Remove(n string) {
	for i, val := range c.items {
		if val.Key() == n {
			c.items = append(c.items[:i], c.items[i+1:]...)
			break
		}
	}
}

func (c *NamedTypeCollection[T]) Get(n string) (noVal T) {
	for _, item := range c.items {
		if item.Key() == n {
			return item
		}
	}

	return noVal
}

func (c *NamedTypeCollection[T]) Sort() {
	slices.SortFunc(c.items, func(i, j T) int {
		return int(i.Compare(j))
	})
}

func (c *NamedTypeCollection[T]) Len() int { return len(c.items) }
func (c *NamedTypeCollection[T]) Less(i, j int) bool {
	a := c.items[i]
	b := c.items[j]

	return a.Compare(b) == LessThan
}
func (c *NamedTypeCollection[T]) Swap(i, j int) { c.items[i], c.items[j] = c.items[j], c.items[i] }

func (c *NamedTypeCollection[T]) Contains(i T) bool {
	for _, item := range c.items {
		if item.Compare(i) == Equal {
			return true
		}
	}

	return false
}

func (c *NamedTypeCollection[T]) ContainsKey(n string) bool {
	for _, item := range c.items {
		if item.Key() == n {
			return true
		}
	}

	return false
}

func CompareItems(a, b CollectionItem) Comparison {
	if a == nil || b == nil {
		return DifferentTypes
	}

	if a.Key() < b.Key() {
		return LessThan
	} else if a.Key() > b.Key() {
		return GreaterThan
	}

	return Equal
}
