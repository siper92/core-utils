package type_utils

import (
	"strings"
	"testing"
)

var _ ICollection[TestItem] = (*Collection[TestItem])(nil)

type TestItem struct {
	name string
}

func (t TestItem) CompareKey(other any) Comparison {
	return CompareItems(t, other)
}

func (t TestItem) Key() string {
	return strings.ToLower(t.name)
}

func (t TestItem) Name() string {
	return t.name
}

func Test_BasicCollection(t *testing.T) {
	collection := NewCollection[TestItem]()

	collection.Add(TestItem{name: "a"})
	collection.Add(TestItem{name: "z"})
	collection.Add(TestItem{name: "s"})
	collection.Add(TestItem{name: "d"})

	if collection == nil {
		t.Fatal("ICollection not created")
	} else if collection.Len() != 4 {
		t.Fatalf("Invalid length: %d", collection.Len())
	}

	if !collection.ContainsKey("d") {
		t.Fatal("ICollection does not contain name b")
	}

	collection.Remove("s")
	if collection.ContainsKey("s") {
		t.Fatal("ICollection remove failed")
	} else if collection.Len() != 3 {
		t.Fatalf("Invalid length after remove: %d", collection.Len())
	}

	collection.Sort()
	if collection.Slice()[0].Name() != "a" &&
		collection.Slice()[1].Name() != "d" &&
		collection.Slice()[2].Name() != "s" {
		t.Fatal("ICollection sort failed")
	}

	newItem := TestItem{name: "zzzz"}
	collection.Add(newItem)
	if !collection.Contains(newItem) {
		t.Errorf("ICollection does not contain new item %s", newItem.Name())
	}

	collection.Sort()
	if collection.Slice()[3].Name() != "zzzz" {
		t.Fatal("ICollection sort failed")
	}

	collection.Add(TestItem{name: "bbb"})
	collection.Sort()
	if collection.Slice()[0].Name() != "a" &&
		collection.Slice()[1].Name() != "bbb" &&
		collection.Slice()[2].Name() != "d" &&
		collection.Slice()[3].Name() != "s" &&
		collection.Slice()[4].Name() != "zzzz" {
		t.Fatal("ICollection sort failed")
	}
}
