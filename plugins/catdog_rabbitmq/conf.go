package catdog_rabbitmq

import (
	"encoding/json"
	"errors"
)

type rabbitConfig struct {
	URL string
}

func Parse(value []byte) (*rabbitConfig, error) {
	cfg := new(rabbitConfig)
	err := json.Unmarshal(value, cfg)
	if err != nil {
		return nil, errors.New("json unmarshal, error=" + err.Error())
	}

	return cfg, nil
}
