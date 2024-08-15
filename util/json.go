package util

import "encoding/json"

func ToJson(obj any) string {
	b, _ := json.Marshal(obj)
	return string(b)
}

func ToPrettyJson(obj any) string {
	b, _ := json.MarshalIndent(obj, "", "  ")
	return string(b)
}
