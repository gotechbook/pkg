package yaml

import (
	"github.com/gotechbook/pkg/codec"
	"gopkg.in/yaml.v3"
)

// Name is the name registered for the yaml codec.
const Name = "yaml"

func init() {
	codec.RegisterCodec(code{})
}

// codec is a Codec implementation with yaml.
type code struct{}

func (code) Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

func (code) Unmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}

func (code) Name() string {
	return Name
}
