package envconfig_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yakumo-saki/go-envconfig"
)

type EnvMapIntConfig struct {
	Map map[string]int `cfg:"MAP_"`
}

func TestEnvMapIntConfig(t *testing.T) {
	assert := assert.New(t)

	ec := envconfig.New()

	t.Setenv("MAP_INT_KEY1-1", "11")
	t.Setenv("MAP_INT_KEY1-2", "12")
	t.Setenv("MAP_INT_KEY2-1", "21")
	t.Setenv("MAP_INT_KEY2-2", "22")

	// confusing key
	t.Setenv("MAP_MAP_MAP_MAP", "12345")

	cfg := EnvMapIntConfig{}
	ec.EnableLogWithDefaultLogger()
	err := ec.LoadConfig(&cfg)
	if err != nil {
		assert.Fail(err.Error())
	}

	assert.Equal(5, len(cfg.Map))
	assert.Equal(11, cfg.Map["INT_KEY1-1"])
	assert.Equal(12, cfg.Map["INT_KEY1-2"])
	assert.Equal(21, cfg.Map["INT_KEY2-1"])
	assert.Equal(22, cfg.Map["INT_KEY2-2"])
	assert.Equal(12345, cfg.Map["MAP_MAP_MAP"])
}
