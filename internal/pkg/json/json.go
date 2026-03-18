package json

import (
	"encoding/json"
)

// ToJSONString converts any interface to a pretty-printed JSON string.
func ToJSONString(i interface{}) string {
	bytes, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}
