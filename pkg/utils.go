package pkg

func GetKeys[T comparable, N any](v map[T]N) []T {
	keys := make([]T, 0, len(v))
	for key := range v {
		keys = append(keys, key)
	}
	return keys
}
