package sliceutil

func Contains[E comparable, S ~[]E](s S, e E) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}
