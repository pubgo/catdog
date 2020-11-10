package yaml

import (
	"github.com/ghodss/yaml"
	"github.com/pubgo/xerror"

	"github.com/asim/nitro/v3/config/encoder"
)

type yamlEncoder struct{}

func (y yamlEncoder) Encode(v interface{}) ([]byte, error) {
	dt, err := yaml.Marshal(v)
	return dt, xerror.Wrap(err)
}

func (y yamlEncoder) Decode(d []byte, v interface{}) error {
	return xerror.Wrap(yaml.Unmarshal(d, v))
}

func (y yamlEncoder) String() string {
	return "yaml"
}

func NewEncoder() encoder.Encoder {
	return yamlEncoder{}
}
