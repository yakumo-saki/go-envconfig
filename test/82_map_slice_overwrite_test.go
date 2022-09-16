package envconfig_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yakumo-saki/go-envconfig"
)

type EnvMapSliceOverwriteConfig struct {
	StrSliceMap map[string][]string `cfg:"STR_SLICE_MAP_,overwrite"`
	IntSliceMap map[string][]int    `cfg:"INT_SLICE_MAP_,overwrite"`
}

func TestMapStrSliceOverwrite(t *testing.T) {
	assert := assert.New(t)
	ec := envconfig.New()
	ec.AddPath("data/map/map_merge_slice_test.env")

	t.Setenv("STR_SLICE_MAP_STRKEY2_1", "STR2-1")
	t.Setenv("STR_SLICE_MAP_STRKEY2_2", "STR2-2")

	cfg := EnvMapSliceOverwriteConfig{}
	ec.EnableLogWithDefaultLogger()
	err := ec.LoadConfig(&cfg)
	if err != nil {
		assert.Fail(err.Error())
	}

	assert.Equal(1, len(cfg.StrSliceMap))
	slice, ok := cfg.StrSliceMap["STRKEY2"]
	assert.True(ok)
	assert.Equal(2, len(slice))
	assert.Equal("STR2-1", slice[0])
	assert.Equal("STR2-2", slice[1])
}

func TestMapIntSliceOverwrite(t *testing.T) {
	assert := assert.New(t)
	ec := envconfig.New()
	ec.AddPath("data/map/map_merge_slice_test.env")

	t.Setenv("INT_SLICE_MAP_INTKEY2_1", "100")
	t.Setenv("INT_SLICE_MAP_INTKEY2_2", "101")

	cfg := EnvMapSliceOverwriteConfig{}
	ec.EnableLogWithDefaultLogger()
	err := ec.LoadConfig(&cfg)
	if err != nil {
		assert.Fail(err.Error())
	}

	assert.Equal(1, len(cfg.IntSliceMap))
	slice, ok := cfg.IntSliceMap["INTKEY2"]
	assert.True(ok)
	assert.Equal(2, len(slice))
	assert.Equal(100, slice[0])
	assert.Equal(101, slice[1])
}
