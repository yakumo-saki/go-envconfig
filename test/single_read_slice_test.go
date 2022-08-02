package envconfig_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yakumo-saki/go-envconfig"
)

type ArrayConfig struct {
	StrArrayConf string `cfg:"STR_ARRAY_CONF,slice"`
}

// Sliceがnilのまま、シンプルに読み込むテスト
func TestLoadSimpleSliceNil(t *testing.T) {
	assert := assert.New(t)
	envconfig.ClearPath()
	envconfig.AddPath("data/simple/simpleArray.env")

	cfg := ArrayConfig{}
	envconfig.LoadConfig(&cfg)

	assert.Equal(3, len(cfg.StrArrayConf))
	assert.Equal("abc", cfg.StrArrayConf[0])
	assert.Equal("def", cfg.StrArrayConf[1])
	assert.Equal("ghi", cfg.StrArrayConf[2])
}
