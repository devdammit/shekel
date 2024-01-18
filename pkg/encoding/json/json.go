package json

import (
	json "encoding/json"

	jsoniter "github.com/json-iterator/go"
	"github.com/mailru/easyjson"
)

var Marshal = func(v interface{}) ([]byte, error) {
	if em, ok := v.(easyjson.Marshaler); ok {
		return easyjson.Marshal(em)
	}
	return json.Marshal(v)
}

var MarshalIndent = json.MarshalIndent

var Unmarshal = func(data []byte, v interface{}) error {
	if em, ok := v.(easyjson.Unmarshaler); ok {
		return easyjson.Unmarshal(data, em)
	}
	return jsoniter.ConfigFastest.Unmarshal(data, v)
}

var NewEncoder = json.NewEncoder
var NewDecoder = jsoniter.ConfigFastest.NewDecoder

type RawMessage = json.RawMessage
type Number = json.Number
