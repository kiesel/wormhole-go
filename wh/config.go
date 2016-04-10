package wormhole

import (
	"errors"
	"strings"
)

type Error interface {
	Error() string
}

type WormholeConfig struct {
	Port    int               `yaml:"port,omitempty"`
	Mapping map[string]string `yaml:"mapping"`
	Editors map[string]string `yaml:"editors"`
}

func (this *WormholeConfig) GetPort() int {
	if 0 == this.Port {
		return 5115
	}

	return this.Port
}

func (this *WormholeConfig) GetMapping(key string) (executable string, err Error) {
	executable, ok := this.Mapping[key]

	if !ok {
		return "", errors.New("No mapping for '" + key + "'")
	}

	return executable, nil
}

func (this *WormholeConfig) AvailableMappings() string {
	var keys []string

	for key, _ := range this.Mapping {
		keys = append(keys, key)
	}

	return strings.Join(keys, ", ")
}
