package parse

import "encoding/json"

func readJson[T any](byteJson []byte) (T, error) {
	var answer T
	err := json.Unmarshal(byteJson, &answer)
	return answer, err
}
