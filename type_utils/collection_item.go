package type_utils

import (
	"fmt"
	core_utils "github.com/siper92/core-utils"
	"strings"
)

func getValueKey(v any) string {
	if v == nil {
		return ""
	}

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

func CompareItemsKeys(a, b any) int {
	compValA := getValueKey(a)
	compValB := getValueKey(b)

	return CompareString(compValA, compValB)
}

func CompareKeysInsensitive(a, b any) int {
	valA := strings.ToLower(strings.TrimSpace(getValueKey(a)))
	valB := strings.ToLower(strings.TrimSpace(getValueKey(b)))

	return CompareString(valA, valB)
}

func CompareString(a, b string) int {
	if a == b {
		return Equal
	} else if a < b {
		return LessThan
	}

	return GreaterThan
}
