package wormhole

import (
	_ "fmt"
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

func TestReadConfigWithArray(t *testing.T) {
	config, err := ReadConfiguration([]byte(`
apps:
  start: ["cmd.exe", "/c", "start"]
  `))

	if err != nil {
		t.Error("Failure during parsing", err.Error())
	}

	app, err := config.GetApp("start")
	if err != nil {
		t.Error("Failure getting app.")
	}

	deepEqual("cmd.exe", app.Executable, t)
}
