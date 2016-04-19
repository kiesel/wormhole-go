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
	App     map[string]App    `yaml:"apps"`
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

func (this *WormholeConfig) GetApp(key string) (*App, Error) {
	app, ok := this.App[key]

	if !ok {
		return nil, errors.New("No mapping for '" + key + "'")
	}

	return &app, nil
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

// Callback for YAML unmarshalling
func (this *App) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {

	// First attempt: just treat value as a single string - that is executable w/o any
	// arguments
	var line string
	if err = unmarshal(&line); err == nil {

		// Try to split string into executable and optional args...
		this.Executable = strings.TrimSpace(line)
		this.Args = []string{}
		return err
	}

	// Attempt: treat as array value, first being the executable all others arguments
	var array []string
	if err = unmarshal(&array); err == nil {
		this.Executable = array[0]
		this.Args = array[1:]
	}

	return err
}
