package type_utils

func InArray[T comparable](val T, arr []T) bool {
	for _, s := range arr {
		if s == val {
			return true
		}
	}

	return false
}

func FindElement[T any](arr []T, predicate func(T) bool) *T {
	for _, el := range arr {
		if predicate(el) {
			return &el
		}
	}

	return nil
}

func GetFiltered[T any](arr []T, predicate func(T) bool) []T {
	var result []T

	for _, el := range arr {
		if predicate(el) {
			result = append(result, el)
		}
	}

	return result
}
