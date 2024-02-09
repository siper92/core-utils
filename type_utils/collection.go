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

type ItemWithName interface {
	Name() string
}

type ItemWithKey interface {
	CompareKey(other any) Comparison
	Key() string
}

type CollectionItem interface {
	ItemWithKey
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
		if val.CompareKey(key) == Equal {
			c.items = append(c.items[:i], c.items[i+1:]...)
			break
		}
	}
}

func (c *Collection[T]) Get(key string) (noVal T) {
	for _, item := range c.items {
		if item.CompareKey(key) == Equal {
			return item
		}
	}

	return noVal
}

func (c *Collection[T]) Sort() ICollection[T] {
	slices.SortFunc(c.items, func(i, j T) int {
		return int(i.CompareKey(j))
	})

	return c
}

func (c *Collection[T]) Len() int { return len(c.items) }
func (c *Collection[T]) Less(i, j int) bool {
	a := c.items[i]
	b := c.items[j]

	return a.CompareKey(b) == LessThan
}
func (c *Collection[T]) Swap(i, j int) { c.items[i], c.items[j] = c.items[j], c.items[i] }

func (c *Collection[T]) Contains(i T) bool {
	for _, item := range c.items {
		if item.CompareKey(i) == Equal {
			return true
		}
	}

	return false
}

func (c *Collection[T]) ContainsKey(n string) bool {
	for _, item := range c.items {
		if item.CompareKey(n) == Equal {
			return true
		}
	}

	return false
}

func getValueKey(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", val)
	case CollectionItem:
		return val.Key()
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

	compValA := getValueKey(a)
	compValB := getValueKey(b)

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
