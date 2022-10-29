package envconfig_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yakumo-saki/go-envconfig"
)

// same as keymerge
type EnvMapSliceDefaultMergeConfig struct {
	StrSliceMap map[string][]string `cfg:"STR_SLICE_MAP_"`
	IntSliceMap map[string][]int    `cfg:"INT_SLICE_MAP_"`
	Check       string              `cfg:"CHECK_READ"`
}

// same as default merge
type EnvMapSliceKeyMergeConfig struct {
	StrSliceMap map[string][]string `cfg:"STR_SLICE_MAP_,keymerge"`
	IntSliceMap map[string][]int    `cfg:"INT_SLICE_MAP_,keymerge"`
	Check       string              `cfg:"CHECK_READ"`
}

func TestMapStrSliceDefaultMerge(t *testing.T) {
	cfg := EnvMapSliceDefaultMergeConfig{}

	assert := assert.New(t)
	ec := envconfig.New()
	ec.AddPath("../data/mapslice/map_merge_slice_test.env")

	t.Setenv("STR_SLICE_MAP_STRKEY1_1", "STR1-1")
	t.Setenv("STR_SLICE_MAP_STRKEY1_2", "STR1-2")
	t.Setenv("INT_SLICE_MAP_INTKEY1_1", "100")
	t.Setenv("INT_SLICE_MAP_INTKEY1_2", "101")

	ec.EnableLogWithDefaultLogger()
	err := ec.LoadConfig(&cfg)
	if err != nil {
		assert.Fail(err.Error())
	}

	assert.Equal("OK", cfg.Check)

	assert.Equal(1, len(cfg.StrSliceMap))
	slice, ok := cfg.StrSliceMap["STRKEY1"]
	assert.True(ok)
	assert.Equal(2, len(slice))
	assert.Equal("STR1-1", slice[0])
	assert.Equal("STR1-2", slice[1])

	assert.Equal(1, len(cfg.IntSliceMap))
	intslice, ok := cfg.IntSliceMap["INTKEY1"]
	assert.True(ok)
	assert.Equal(2, len(intslice))
	assert.Equal(100, intslice[0])
	assert.Equal(101, intslice[1])
}

// 中身は3行目から上と同じ
func TestMapStrSliceKeyMerge(t *testing.T) {
	cfg := EnvMapSliceKeyMergeConfig{}

	assert := assert.New(t)
	ec := envconfig.New()
	ec.AddPath("../data/mapslice/map_merge_slice_test.env")

	t.Setenv("STR_SLICE_MAP_STRKEY1_1", "STR1-1")
	t.Setenv("STR_SLICE_MAP_STRKEY1_2", "STR1-2")
	t.Setenv("INT_SLICE_MAP_INTKEY1_1", "100")
	t.Setenv("INT_SLICE_MAP_INTKEY1_2", "101")

	ec.EnableLogWithDefaultLogger()
	err := ec.LoadConfig(&cfg)
	if err != nil {
		assert.Fail(err.Error())
	}

	assert.Equal("OK", cfg.Check)

	assert.Equal(1, len(cfg.StrSliceMap))
	slice, ok := cfg.StrSliceMap["STRKEY1"]
	assert.True(ok)
	assert.Equal(2, len(slice))
	assert.Equal("STR1-1", slice[0])
	assert.Equal("STR1-2", slice[1])

	assert.Equal(1, len(cfg.IntSliceMap))
	intslice, ok := cfg.IntSliceMap["INTKEY1"]
	assert.True(ok)
	assert.Equal(2, len(intslice))
	assert.Equal(100, intslice[0])
	assert.Equal(101, intslice[1])

}
