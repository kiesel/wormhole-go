package wormhole

import (
	"fmt"
	"testing"
)

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
		t.Errorf("Failure during parsing " + err.Error())
	}

	fmt.Println(config)
}
