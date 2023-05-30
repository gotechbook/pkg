package json

import (
	"encoding/json"
	"github.com/gotechbook/pkg/codec"
)

// Name is the name registered for the yaml codec.
const Name = "json"

func init() {
	codec.RegisterCodec(code{})
}

// codec is a Codec implementation with yaml.
type code struct{}

func (code) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (code) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (code) Name() string {
	return Name
}
