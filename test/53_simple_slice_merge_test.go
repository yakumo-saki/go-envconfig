package envconfig_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yakumo-saki/go-envconfig"
)

type MergeSliceConfig struct {
	StrArrayConf   []string  `cfg:"STR_ARRAY_,merge"`
	IntArrayConf   []int     `cfg:"INT_ARRAY_,merge"`
	FloatArrayConf []float64 `cfg:"FLOAT_ARRAY_,merge"`
}

func TestMergeSlice(t *testing.T) {
	assert := assert.New(t)
	envconfig.ClearPath()
	envconfig.AddPath("data/simple/merge_slice_config1.env")
	envconfig.AddPath("data/simple/merge_slice_config2.env")

	cfg := MergeSliceConfig{}
	envconfig.EnableLogWithDefaultLogger()
	err := envconfig.LoadConfig(&cfg)
	if err != nil {
		assert.Fail(err.Error())
	}

	assert.Equal(4, len(cfg.StrArrayConf))
	assert.Equal("config1_abc", cfg.StrArrayConf[0])
	assert.Equal("config1_def", cfg.StrArrayConf[1])
	assert.Equal("config2_abc", cfg.StrArrayConf[2])
	assert.Equal("config2_def", cfg.StrArrayConf[3])

	assert.Equal(4, len(cfg.IntArrayConf))
	assert.Equal(11, cfg.IntArrayConf[0])
	assert.Equal(12, cfg.IntArrayConf[1])
	assert.Equal(21, cfg.IntArrayConf[2])
	assert.Equal(22, cfg.IntArrayConf[3])

	assert.Equal(4, len(cfg.FloatArrayConf))
	assert.Equal(1.1, cfg.FloatArrayConf[0])
	assert.Equal(1.2, cfg.FloatArrayConf[1])
	assert.Equal(2.1, cfg.FloatArrayConf[2])
	assert.Equal(2.2, cfg.FloatArrayConf[3])
}
