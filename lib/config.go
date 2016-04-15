package wormhole

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
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

func GetDefaultLog() string {
	return path.Join(os.Getenv("HOME"), "wormhole.log")
}

func ReadConfigurationFrom(filename string) (config *WormholeConfig, err Error) {
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Critical(err.Error())

		return nil, err
	}

	return ReadConfiguration(source)
}

func ReadConfiguration(source []byte) (config *WormholeConfig, err Error) {
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		log.Critical(err.Error())
		return nil, err
	}

	return config, nil
}

func (this *WormholeConfig) GetAddr() string {
	if "" == this.Addr {
		return "127.0.0.1:5115"
	}

	return this.Addr
}

func (this *WormholeConfig) GetApp(key string) (app *App, err Error) {
	cmdline, ok := this.App[key]

	if !ok {
		return nil, errors.New("No mapping for '" + key + "'")
	}

	return &App{
		Executable: cmdline,
		Args:       []string{},
	}, nil
}

func (this *WormholeConfig) GetAppWith(key string, args []string) (app *App, err Error) {
	if app, err = this.GetApp(key); err != nil {
		return app, err
	}

	app.MergeArguments(this.translatePaths(args))
	return app, nil
}

func (this *WormholeConfig) AvailableApps() string {
	var keys []string

	for key, _ := range this.App {
		keys = append(keys, key)
	}

	return strings.Join(keys, ", ")
}

func (this *WormholeConfig) translatePath(path string) string {
	for from, to := range this.Mapping {
		path = strings.Replace(path, from, to, 1)
	}

	return filepath.FromSlash(path)
}

func (this *WormholeConfig) translatePaths(paths []string) []string {
	for index, arg := range paths {
		paths[index] = this.translatePath(arg)
	}

	return paths
}

type App struct {
	Executable string
	Args       []string
}

func (this *App) MergeArguments(args []string) {
	this.Args = append(this.Args, args...)
}
