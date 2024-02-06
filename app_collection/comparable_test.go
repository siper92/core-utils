package app_collection

import "testing"

var _ Collection[TestItem] = (*NamedTypeCollection[TestItem])(nil)

type TestItem struct {
	key string
}

func (t TestItem) Key() string {
	return t.key
}

func (t TestItem) Compare(other Comparable) Comparison {
	val, ok := other.(TestItem)
	if !ok {
		return DifferentTypes
	}

	if t.key < val.key {
		return LessThan
	} else if t.key > val.key {
		return GreaterThan
	}

	return Equal
}

func Test_BasicCollection(t *testing.T) {
	collection := NewCollection[TestItem]()

	collection.Add(TestItem{key: "a"})
	collection.Add(TestItem{key: "z"})
	collection.Add(TestItem{key: "s"})
	collection.Add(TestItem{key: "d"})

	if collection == nil {
		t.Fatal("Collection not created")
	} else if collection.Len() != 4 {
		t.Fatalf("Invalid length: %d", collection.Len())
	}

	if !collection.ContainsKey("d") {
		t.Fatal("Collection does not contain key b")
	}

	collection.Remove("s")
	if collection.ContainsKey("s") {
		t.Fatal("Collection remove failed")
	} else if collection.Len() != 3 {
		t.Fatalf("Invalid length after remove: %d", collection.Len())
	}

	collection.Sort()
	if collection.Items()[0].Key() != "a" &&
		collection.Items()[1].Key() != "d" &&
		collection.Items()[2].Key() != "s" {
		t.Fatal("Collection sort failed")
	}

	newItem := TestItem{key: "zzzz"}
	collection.Add(newItem)
	if !collection.Contains(newItem) {
		t.Errorf("Collection does not contain new item %s", newItem.Key())
	}

	collection.Sort()
	if collection.Items()[3].Key() != "zzzz" {
		t.Fatal("Collection sort failed")
	}

	collection.Add(TestItem{key: "bbb"})
	collection.Sort()
	if collection.Items()[0].Key() != "a" &&
		collection.Items()[1].Key() != "bbb" &&
		collection.Items()[2].Key() != "d" &&
		collection.Items()[3].Key() != "s" &&
		collection.Items()[4].Key() != "zzzz" {
		t.Fatal("Collection sort failed")
	}
}
