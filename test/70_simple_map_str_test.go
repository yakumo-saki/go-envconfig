package envconfig_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yakumo-saki/go-envconfig"
)

type EnvMapStringConfig struct {
	StrMap map[string]string `cfg:"STR_MAP_,map"`
}

func TestEnvMapStringConfig(t *testing.T) {
	assert := assert.New(t)
	envconfig.ClearPath()

	SetEnv("STR_MAP_STR_KEY1-1", "STR1-1")
	SetEnv("STR_MAP_STR_KEY1-2", "STR1-2")
	SetEnv("STR_MAP_STR_KEY1-3", "STR1-3")
	SetEnv("STR_MAP_STR_KEY2-1", "STR2-1")
	SetEnv("STR_MAP_STR_KEY2-2", "STR2-2")

	cfg := EnvMapStringConfig{}
	envconfig.EnableLogWithDefaultLogger()
	err := envconfig.LoadConfig(&cfg)
	if err != nil {
		assert.Fail(err.Error())
	}

	assert.Equal(5, len(cfg.StrMap))
	assert.Equal("STR1-1", cfg.StrMap["KEY1-1"])
	assert.Equal("STR1-2", cfg.StrMap["KEY1-2"])
	assert.Equal("STR1-3", cfg.StrMap["KEY1-3"])
	assert.Equal("STR2-1", cfg.StrMap["KEY2-1"])
	assert.Equal("STR2-2", cfg.StrMap["KEY2-2"])
}
