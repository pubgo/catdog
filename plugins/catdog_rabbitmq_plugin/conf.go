package catdog_rabbitmq_plugin

import (
	"encoding/json"
	"errors"
)

type RbmqConfig struct {
	URL string
}

func Parse(value []byte) (*RbmqConfig, error) {
	rbmqConfig := new(RbmqConfig)
	err := json.Unmarshal(value, rbmqConfig)
	if err != nil {
		return nil, errors.New("json unmarshal, error=" + err.Error())
	}

	return rbmqConfig, nil
}
