package envconfig_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yakumo-saki/go-envconfig"
)

type EnvMapSliceMergeConfig struct {
	StrSliceMap map[string][]string `cfg:"STR_SLICE_MAP_"`
	IntSliceMap map[string][]int    `cfg:"INT_SLICE_MAP_"`
}

func TestMapStrSliceMerge(t *testing.T) {
	assert := assert.New(t)
	ec := envconfig.New()
	ec.AddPath("data/map/map_merge_slice_test.env")

	t.Setenv("STR_SLICE_MAP_STRKEY1_1", "STR1-1")
	t.Setenv("STR_SLICE_MAP_STRKEY1_2", "STR1-2")

	cfg := EnvMapSliceMergeConfig{}
	ec.EnableLogWithDefaultLogger()
	err := ec.LoadConfig(&cfg)
	if err != nil {
		assert.Fail(err.Error())
	}

	assert.Equal(1, len(cfg.StrSliceMap))
	slice, ok := cfg.StrSliceMap["STRKEY1"]
	assert.True(ok)
	assert.Equal("FILESTR1", slice[0])
	assert.Equal("FILESTR2", slice[1])
	assert.Equal("STR1-1", slice[2])
	assert.Equal("STR1-2", slice[3])
}

func TestMapIntSliceMerge(t *testing.T) {
	assert := assert.New(t)
	ec := envconfig.New()
	ec.AddPath("data/map/map_merge_slice_test.env")

	t.Setenv("INT_SLICE_MAP_INTKEY1_1", "100")
	t.Setenv("INT_SLICE_MAP_INTKEY1_2", "101")

	cfg := EnvMapSliceMergeConfig{}
	ec.EnableLogWithDefaultLogger()
	err := ec.LoadConfig(&cfg)
	if err != nil {
		assert.Fail(err.Error())
	}

	assert.Equal(1, len(cfg.IntSliceMap))
	slice, ok := cfg.IntSliceMap["INTKEY1"]
	assert.True(ok)
	assert.Equal(len(slice), 4)
	assert.Equal(99999, slice[0])
	assert.Equal(88888, slice[1])
	assert.Equal(100, slice[2])
	assert.Equal(101, slice[3])
}
