package envconfig_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yakumo-saki/go-envconfig"
)

type DuplicatedBadConfig struct {
	A string `cfg:"STR_CONF"`
	B string `cfg:"STR_CONF"`
}

// test multiple same configkey => panic
func TestDuplicated(t *testing.T) {
	defer func() {
		err := recover()

		if !strings.HasPrefix(fmt.Sprintf("%s", err), "config key duplicated:") {
			t.Errorf("got %v\nwant %v", err, "config key duplicated:")
			assert.Fail(t, "wrong panic message")
		}
	}()

	envconfig.ClearPath()
	envconfig.AddPath("data/simple/simple.env")

	cfg := DuplicatedBadConfig{}
	envconfig.EnableLogWithDefaultLogger()
	envconfig.LoadConfig(&cfg)

}
