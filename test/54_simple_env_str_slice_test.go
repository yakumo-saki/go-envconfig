package envconfig_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yakumo-saki/go-envconfig"
)

type EnvStringSliceConfig struct {
	StrSlice []string `cfg:"STR_SLICE_,slice"`
}

func TestLoadEnvStringSlice(t *testing.T) {
	assert := assert.New(t)
	envconfig.ClearPath()

	os.Setenv("STR_SLICE_0", "STR0")
	os.Setenv("STR_SLICE_1", "STR1")

	cfg := EnvStringSliceConfig{}
	envconfig.EnableLogWithDefaultLogger()
	err := envconfig.LoadConfig(&cfg)
	if err != nil {
		assert.Fail(err.Error())
	}

	assert.Equal(2, len(cfg.StrSlice))
	assert.Equal("STR0", cfg.StrSlice[0])
	assert.Equal("STR1", cfg.StrSlice[1])
}
