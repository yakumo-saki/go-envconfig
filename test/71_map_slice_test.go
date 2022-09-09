package envconfig_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yakumo-saki/go-envconfig"
)

type EnvMapSliceStringConfig struct {
	StrSliceMap map[string][]string `cfg:"STR_SLICE_MAP_"`
	IntSliceMap map[string][]int    `cfg:"INT_SLICE_MAP_"`
}

func TestEnvMapStrSlice(t *testing.T) {
	assert := assert.New(t)
	envconfig.ClearPath()

	t.Setenv("STR_SLICE_MAP_STRKEY1_1", "STR1-1")
	t.Setenv("STR_SLICE_MAP_STRKEY1_2", "STR1-2")
	t.Setenv("STR_SLICE_MAP_STRKEY1_3", "STR1-3")
	t.Setenv("STR_SLICE_MAP_STRKEY2_1", "STR2-1")
	t.Setenv("STR_SLICE_MAP_STRKEY2_2", "STR2-2")

	cfg := EnvMapSliceStringConfig{}
	envconfig.EnableLogWithDefaultLogger()
	err := envconfig.LoadConfig(&cfg)
	if err != nil {
		assert.Fail(err.Error())
	}

	assert.Equal(2, len(cfg.StrSliceMap))
	slice, ok := cfg.StrSliceMap["STRKEY1"]
	assert.True(ok)
	assert.Equal(len(slice), 3)
	assert.Equal("STR1-1", slice[0])
	assert.Equal("STR1-2", slice[1])
	assert.Equal("STR1-3", slice[2])

	slice, ok = cfg.StrSliceMap["STRKEY2"]
	assert.True(ok)
	assert.Equal(len(slice), 2)
	assert.Equal("STR2-1", slice[0])
	assert.Equal("STR2-2", slice[1])
}

func TestEnvMapIntSlice(t *testing.T) {
	assert := assert.New(t)
	envconfig.ClearPath()

	t.Setenv("INT_SLICE_MAP_INTKEY1_1", "100")
	t.Setenv("INT_SLICE_MAP_INTKEY1_2", "101")
	t.Setenv("INT_SLICE_MAP_INTKEY1_3", "102")
	t.Setenv("INT_SLICE_MAP_INTKEY2_1", "200")
	t.Setenv("INT_SLICE_MAP_INTKEY2_2", "201")

	cfg := EnvMapSliceStringConfig{}
	envconfig.EnableLogWithDefaultLogger()
	err := envconfig.LoadConfig(&cfg)
	if err != nil {
		assert.Fail(err.Error())
	}

	assert.Equal(2, len(cfg.IntSliceMap))
	slice, ok := cfg.IntSliceMap["INTKEY1"]
	assert.True(ok)
	assert.Equal(len(slice), 3)
	assert.Equal(100, slice[0])
	assert.Equal(101, slice[1])
	assert.Equal(102, slice[2])

	slice, ok = cfg.IntSliceMap["INTKEY2"]
	assert.True(ok)
	assert.Equal(len(slice), 2)
	assert.Equal(200, slice[0])
	assert.Equal(201, slice[1])
}
