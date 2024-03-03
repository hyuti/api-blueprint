package tool

import "encoding/json"

func JSONStringify[T any](v T) string {
	b, err := json.Marshal(v)
	if err != nil {
		return "unable to marshal to convert to str"
	}
	return string(b)
}
