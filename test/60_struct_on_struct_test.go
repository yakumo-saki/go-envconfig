package envconfig_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yakumo-saki/go-envconfig"
)

type MultiTierConfig struct {
	SubConfig SubConfigStruct
}

type SubConfigStruct struct {
	SubStringConf string `cfg:"SUB_STRING_CONF"`
}

func TestMultiTierConfig(t *testing.T) {
	assert := assert.New(t)
	envconfig.ClearPath()
	envconfig.AddPath("data/simple/simpleStrArray.env")

	cfg := MultiTierConfig{}
	envconfig.EnableLog()
	err := envconfig.LoadConfig(&cfg)
	if err != nil {
		assert.Fail(err.Error())
	}

	assert.Equal("abc0", cfg.SubConfig.SubStringConf)
}
