package envconfig_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yakumo-saki/go-envconfig"
)

type OverwriteSliceConfig struct {
	StrArrayConf   []string  `cfg:"STR_ARRAY_,slice"`
	IntArrayConf   []int     `cfg:"INT_ARRAY_,slice"`
	FloatArrayConf []float64 `cfg:"FLOAT_ARRAY_,slice"`
}

func TestOverwriteSlice(t *testing.T) {
	assert := assert.New(t)
	envconfig.ClearPath()
	envconfig.AddPath("data/simple/merge_slice_config1.env")
	envconfig.AddPath("data/simple/merge_slice_config2.env")

	cfg := OverwriteSliceConfig{}
	envconfig.EnableLogWithDefaultLogger()
	err := envconfig.LoadConfig(&cfg)
	if err != nil {
		assert.Fail(err.Error())
	}

	assert.Equal(2, len(cfg.StrArrayConf))
	assert.Equal("config2_abc", cfg.StrArrayConf[0])
	assert.Equal("config2_def", cfg.StrArrayConf[1])

	assert.Equal(2, len(cfg.IntArrayConf))
	assert.Equal(21, cfg.IntArrayConf[0])
	assert.Equal(22, cfg.IntArrayConf[1])

	assert.Equal(2, len(cfg.FloatArrayConf))
	assert.Equal(2.1, cfg.FloatArrayConf[0])
	assert.Equal(2.2, cfg.FloatArrayConf[1])
}
