package envconfig_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yakumo-saki/go-envconfig"
)

type EnvMapMergeConfig struct {
	StrMap map[string]string `cfg:"MAP_,merge"`
}

func TestMapMergeConfig(t *testing.T) {
	assert := assert.New(t)

	ec := envconfig.New()
	ec.AddPath("data/map/map_merge_test.env")

	t.Setenv("MAP_STR_KEY1", "STR1-1")
	t.Setenv("MAP_STR_KEY2", "STR2-1")

	cfg := EnvMapMergeConfig{}
	ec.EnableLogWithDefaultLogger()
	err := ec.LoadConfig(&cfg)
	if err != nil {
		assert.Fail(err.Error())
	}

	assert.Equal("STR1-1", cfg.StrMap["STR_KEY1"])
	assert.Equal("STR2-1", cfg.StrMap["STR_KEY2"])
	assert.Equal("A1", cfg.StrMap["STR_KEY_A1"])
	assert.Equal("A2", cfg.StrMap["STR_KEY_A2"])
}
