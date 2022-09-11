package envconfig_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yakumo-saki/go-envconfig"
)

type StrArrayConfig struct {
	StrArrayConf []string `cfg:"STR_ARRAY_CONF_"`
	Dummy1       string   `cfg:""`
	Dummy2       string   `cfg:""`
}

func TestLoadSimpleStringSlice(t *testing.T) {
	assert := assert.New(t)

	ec := envconfig.New()

	ec.ClearPath()
	ec.AddPath("data/simple/simpleStrArray.env")

	cfg := StrArrayConfig{}
	ec.EnableLogWithDefaultLogger()
	err := ec.LoadConfig(&cfg)
	if err != nil {
		assert.Fail(err.Error())
	}

	assert.Equal(3, len(cfg.StrArrayConf))
	assert.Equal("abc0", cfg.StrArrayConf[0])
	assert.Equal("def1", cfg.StrArrayConf[1])
	assert.Equal("ghi2", cfg.StrArrayConf[2])
}

// ゼロパディングのテスト
func TestLoadSimpleStringSliceZeroPadding(t *testing.T) {
	assert := assert.New(t)
	ec := envconfig.New()
	ec.AddPath("data/simple/simpleStrArrayPadding.env")

	cfg := StrArrayConfig{}
	ec.EnableLogWithDefaultLogger()
	err := ec.LoadConfig(&cfg)
	if err != nil {
		assert.Fail(err.Error())
	}

	assert.Equal(5, len(cfg.StrArrayConf))
	assert.Equal("abc0", cfg.StrArrayConf[0])
	assert.Equal("def00", cfg.StrArrayConf[1])
	assert.Equal("ghi1", cfg.StrArrayConf[2])
	assert.Equal("jkl01", cfg.StrArrayConf[3])
	assert.Equal("mno2", cfg.StrArrayConf[4])
}
