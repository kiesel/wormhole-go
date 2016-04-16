package wormhole

import (
	"reflect"
	"testing"
)

func deepEqual(expect, actual interface{}, t *testing.T) {
	if !reflect.DeepEqual(expect, actual) {
		t.Error("Items not equal:", expect, actual)
	}
}

func TestDefaultAddress(t *testing.T) {
	config := &WormholeConfig{}
	if "127.0.0.1:5115" != config.GetAddr() {
		t.Error("Wrong default address: " + config.GetAddr())
	}
}

func TestReadSimpleConfig(t *testing.T) {
	config, err := ReadConfiguration([]byte(`
addr: :5115

mapping:
  "/home/": "A:"

apps:
  "sublime": "/opt/sublime/sublime"
  `))

	if err != nil {
		t.Error("Failure during parsing", err.Error())
	}

	deepEqual(map[string]string{"/home/": "A:"}, config.Mapping, t)
	// deepEqual("/opt/sublime/sublime", config.GetApp("sublime"), t)
}

func Test_read_config_unquoted_with_args(t *testing.T) {
	createConfigAndCompareApp(`
apps:
  app: cmd.exe /c start
  `, &App{
		Executable: "cmd.exe",
		Args:       []string{"/c", "start"},
	},
		t)
}

func Test_read_config_quoted(t *testing.T) {
	createConfigAndCompareApp(`
apps:
  app: "cmd.exe /c start"
  `, &App{
		Executable: "cmd.exe",
		Args:       []string{"/c", "start"},
	},
		t)
}

func Test_read_config_array_notation(t *testing.T) {
	createConfigAndCompareApp(`
apps:
  app: ["cmd.exe", "/c", "start"]
  `, &App{
		Executable: "cmd.exe",
		Args:       []string{"/c", "start"},
	},
		t)
}

func Test_read_config_unquoted_with_whitespace_path_leads_to_broken_executable(t *testing.T) {
	createConfigAndCompareApp(`
apps:
  app: cmd with whitespace.exe /c start
  `, &App{
		Executable: "cmd",
		Args:       []string{"with", "whitespace.exe", "/c", "start"},
	},
		t)
}

func Test_read_config_quoted_with_args_in_array_notation(t *testing.T) {
	createConfigAndCompareApp(`
apps:
  app: ["cmd with whitespace.exe", "/c", "start"]
  `, &App{
		Executable: "cmd with whitespace.exe",
		Args:       []string{"/c", "start"},
	},
		t)
}

func createConfigAndCompareApp(cstr string, expect *App, t *testing.T) {
	config, err := ReadConfiguration([]byte(cstr))

	if err != nil {
		t.Error("Failure during parsing", err.Error())
		return
	}

	var app *App
	if app, err = config.GetApp("app"); err != nil {
		t.Error("Failure to get app", err)
		return
	}

	deepEqual(expect, app, t)
}
