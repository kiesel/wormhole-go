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
	App     map[string]string `yaml:"apps"`
}

func (this *WormholeConfig) GetPort() int {
	if 0 == this.Port {
		return 5115
	}

	return this.Port
}

func (this *WormholeConfig) GetApp(key string) (executable string, err Error) {
	executable, ok := this.App[key]

	if !ok {
		return "", errors.New("No mapping for '" + key + "'")
	}

	return executable, nil
}

func (this *WormholeConfig) AvailableApps() string {
	var keys []string

	for key, _ := range this.App {
		keys = append(keys, key)
	}

	return strings.Join(keys, ", ")
}

func (this *WormholeConfig) TranslatePath(path string) string {
	for from, to := range this.Mapping {
		path = strings.Replace(path, from, to, 1)
	}

	return path
}
