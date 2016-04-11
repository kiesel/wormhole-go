package wormhole

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Error interface {
	Error() string
}

type WormholeConfig struct {
	Addr    string            `yaml:"listen,omitempty"`
	Mapping map[string]string `yaml:"mapping"`
	App     map[string]string `yaml:"apps"`
}

func GetDefaultConfig() string {
	return path.Join(os.Getenv("HOME"), ".wormhole.yml")
}

func (this *WormholeConfig) GetAddr() string {
	if "" == this.Addr {
		return "127.0.0.1:5115"
	}

	return this.Addr
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

	return filepath.FromSlash(path)
}
