package type_utils

import "fmt"

func BaseTypeToString(raw any) string {
	val := ""
	switch v := raw.(type) {
	case string:
		val = v
	case *string:
		val = *v
	case []byte:
		val = string(v)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64,
		*int, *int8, *int16, *int32, *int64, *uint, *uint8, *uint16, *uint32, *uint64:
		val = fmt.Sprintf("%d", v)
	case float32, float64:
		val = fmt.Sprintf("%f", v)
	case bool:
		val = fmt.Sprintf("%t", v)
	}

	return val
}
