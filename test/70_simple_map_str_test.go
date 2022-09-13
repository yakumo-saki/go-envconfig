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

	ec := envconfig.New()

	t.Setenv("STR_MAP_STR_KEY1-1", "STR1-1")
	t.Setenv("STR_MAP_STR_KEY1-2", "STR1-2")
	t.Setenv("STR_MAP_STR_KEY1-3", "STR1-3")
	t.Setenv("STR_MAP_STR_KEY2-1", "STR2-1")
	t.Setenv("STR_MAP_STR_KEY2-2", "STR2-2")

	// confusing key  CONFUSING_STR_MAP_KEY -> "confusing"
	t.Setenv("STR_MAP_CONFUSING_STR_MAP_KEY", "confusing")

	cfg := EnvMapStringConfig{}
	ec.EnableLogWithDefaultLogger()
	err := ec.LoadConfig(&cfg)
	if err != nil {
		assert.Fail(err.Error())
	}

	assert.Equal(6, len(cfg.StrMap))
	assert.Equal("STR1-1", cfg.StrMap["STR_KEY1-1"])
	assert.Equal("STR1-2", cfg.StrMap["STR_KEY1-2"])
	assert.Equal("STR1-3", cfg.StrMap["STR_KEY1-3"])
	assert.Equal("STR2-1", cfg.StrMap["STR_KEY2-1"])
	assert.Equal("STR2-2", cfg.StrMap["STR_KEY2-2"])
	assert.Equal("confusing", cfg.StrMap["CONFUSING_STR_MAP_KEY"])
}
