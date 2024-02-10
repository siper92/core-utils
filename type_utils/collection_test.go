package type_utils

import (
	"strings"
	"testing"
)

var _ ICollection[TestItem] = (*Collection[TestItem])(nil)

type TestItem struct {
	name string
}

func (t TestItem) KeyEqual(other any) bool {
	return CompareKeysInsensitive(t, other)
}

func (t TestItem) Compare(other any) int {
	return CompareItemsKeys(t, other)
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

func Test_BasicCollection_Get_And_Keys(t *testing.T) {
	collection := NewCollection[TestItem]()

	collection.Add(TestItem{name: "aRt"})
	collection.Add(TestItem{name: "Zaf"})
	collection.Add(TestItem{name: "s"})
	collection.Add(TestItem{name: "DaS"})
	collection.Add(TestItem{name: "DaS aS "})

	test := []struct {
		name  string
		key   string
		found bool
	}{
		{
			name:  "aRt",
			key:   "art",
			found: true,
		},
		{
			name:  "Zaf",
			key:   "zAf",
			found: true,
		},
		{
			name:  "s",
			key:   "S",
			found: true,
		},
		{
			name:  "DaS",
			key:   "daS",
			found: true,
		},
		{
			name:  "DaS aS ",
			key:   "das as",
			found: true,
		},
		{
			name:  " NotFound ",
			key:   "notfound",
			found: false,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			testItem := collection.Get(tt.key)
			if testItem.Name() == "" {
				if tt.found {
					t.Fatalf("Test failed to find %s", tt.key)
				} else {
					collection.Add(TestItem{name: tt.name})
					testItem = collection.Get(tt.key)
					if testItem.Name() != tt.name {
						t.Fatalf("Test failed to insert %s", tt.name)
					}
				}
			} else {
				if testItem.Name() != tt.name {
					t.Fatalf("Test failed for %s, name is %s", tt.key, testItem.Name())
				} else if !tt.found {
					t.Fatalf("Test failed to not find %s", tt.key)
				}
			}
		})
	}
}
