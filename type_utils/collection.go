package type_utils

import (
	"fmt"
	core_utils "github.com/siper92/core-utils"
	"slices"
	"sort"
)

type Comparison int

const (
	Equal          Comparison = 0
	LessThan       Comparison = -1
	GreaterThan    Comparison = 1
	NotTheSameType Comparison = 1
)

type ComparableItem interface {
	Compare(other any) Comparison
}

type ItemWithName interface {
	Name() string
}

type ItemWithKey interface {
	Key() string
}

type CollectionItem interface {
	ComparableItem
	ItemWithName
	ItemWithKey
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
	Sort() Collection[T]
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
		m[item.Key()] = item
	}

	return m
}

func (c *SimpleCollection[T]) Add(item T) {
	c.items = append(c.items, item)
}

func (c *SimpleCollection[T]) Remove(n string) {
	for i, val := range c.items {
		if val.Key() == n || val.Name() == n {
			c.items = append(c.items[:i], c.items[i+1:]...)
			break
		}
	}
}

func (c *SimpleCollection[T]) Get(n string) (noVal T) {
	for _, item := range c.items {
		if item.Name() == n || item.Key() == n {
			return item
		}
	}

	return noVal
}

func (c *SimpleCollection[T]) Sort() Collection[T] {
	slices.SortFunc(c.items, func(i, j T) int {
		return int(i.Compare(j))
	})

	return c
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
		if item.Key() == n {
			return true
		}
	}

	return false
}

func getCompareValue(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", val)
	case CollectionItem:
		if val.Key() != "" {
			return val.Key()
		}

		return val.Name()
	case ItemWithKey:
		return val.Key()
	case ItemWithName:
		return val.Name()
	default:
		core_utils.Debug("type comparison not implemented %T", v)
		return fmt.Sprintf("%+v", v)
	}
}

func CompareItems(a, b any) Comparison {
	if a == nil {
		if a == b {
			return Equal
		}

		return LessThan
	} else if b == nil {
		return GreaterThan
	}

	if fmt.Sprintf("%T", a) != fmt.Sprintf("%T", b) {
		return NotTheSameType
	}

	compValA := getCompareValue(a)
	compValB := getCompareValue(b)

	return CompareString(compValA, compValB)
}

func CompareString(a, b string) Comparison {
	if a == b {
		return Equal
	} else if a < b {
		return LessThan
	}

	return GreaterThan
}
