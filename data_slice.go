package core_utils

func isContainedIn[T comparable](str T, arr []T) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}

	return false
}
