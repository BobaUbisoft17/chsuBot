package pkg

import "time"

func GetKeys[T comparable, N any](v map[T]N) []T {
	keys := make([]T, 0, len(v))
	for key := range v {
		keys = append(keys, key)
	}
	return keys
}

func StringToTimestamp(date string) (int, error) {
	timeObject, err := time.Parse("02.01.2006", date)
	return int(timeObject.Unix()), err
}
