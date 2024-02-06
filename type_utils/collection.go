package type_utils

import (
	"fmt"
	core_utils "github.com/siper92/core-utils"
	"slices"
	"sort"
)

type Comparison int

const (
	Equal       Comparison = 0
	LessThan    Comparison = -1
	GreaterThan Comparison = 1
)

type Comparable interface {
	Compare(other Comparable) Comparison
}

type HasName interface {
	Name() string
}

type CollectionItem interface {
	Comparable
	HasName
}

type Collection[T CollectionItem] interface {
	sort.Interface
	Slice() []T
	Map() map[string]T
	Contains(i T) bool
	ContainsKey(k string) bool
	Remove(k string)
	Add(item T)
	Get(k string) T
	Sort()
}

var _ Collection[CollectionItem] = (*SimpleCollection[CollectionItem])(nil)

type SimpleCollection[T CollectionItem] struct {
	items []T
}

func NewCollection[T CollectionItem](items ...T) *SimpleCollection[T] {
	if len(items) == 0 {
		return &SimpleCollection[T]{items: make([]T, 0)}
	}

	return &SimpleCollection[T]{items: items}
}

func (c *SimpleCollection[T]) Slice() []T {
	return c.items
}

func (c *SimpleCollection[T]) Map() map[string]T {
	m := make(map[string]T)
	for _, item := range c.items {
		m[item.Name()] = item
	}

	return m
}

func (c *SimpleCollection[T]) Add(item T) {
	c.items = append(c.items, item)
}

func (c *SimpleCollection[T]) Remove(n string) {
	for i, val := range c.items {
		if val.Name() == n {
			c.items = append(c.items[:i], c.items[i+1:]...)
			break
		}
	}
}

func (c *SimpleCollection[T]) Get(n string) (noVal T) {
	for _, item := range c.items {
		if item.Name() == n {
			return item
		}
	}

	return noVal
}

func (c *SimpleCollection[T]) Sort() {
	slices.SortFunc(c.items, func(i, j T) int {
		return int(i.Compare(j))
	})
}

func (c *SimpleCollection[T]) Len() int { return len(c.items) }
func (c *SimpleCollection[T]) Less(i, j int) bool {
	a := c.items[i]
	b := c.items[j]

	return a.Compare(b) == LessThan
}
func (c *SimpleCollection[T]) Swap(i, j int) { c.items[i], c.items[j] = c.items[j], c.items[i] }

func (c *SimpleCollection[T]) Contains(i T) bool {
	for _, item := range c.items {
		if item.Compare(i) == Equal {
			return true
		}
	}

	return false
}

func (c *SimpleCollection[T]) ContainsKey(n string) bool {
	for _, item := range c.items {
		if item.Name() == n {
			return true
		}
	}

	return false
}

func getCompareValue(v interface{}) string {
	switch val := v.(type) {
	case CollectionItem:
		return val.Name()
	case HasName:
		return val.Name()
	default:
		core_utils.Debug("type comparison not implemented %T", v)
		return fmt.Sprintf("%+v", v)
	}
}

func CompareItems(a, b interface{}) Comparison {
	if a == nil {
		if a == b {
			return Equal
		}

		return LessThan
	} else if b == nil {
		return GreaterThan
	}

	compValA := getCompareValue(a)
	compValB := getCompareValue(b)

	if compValA == compValB {
		return Equal
	} else if compValA < compValB {
		return LessThan
	}

	return GreaterThan
}
