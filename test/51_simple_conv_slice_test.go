package envconfig_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yakumo-saki/go-envconfig"
)

type ConvArrayConfig struct {
	IntArrayConf   []int     `cfg:"INT_ARRAY_,slice"`
	FloatArrayConf []float64 `cfg:"FLOAT_ARRAY_,slice"`
}

// Sliceがnilのまま、シンプルに読み込むテスト
func TestLoadSliceWithConvert(t *testing.T) {
	assert := assert.New(t)
	envconfig.ClearPath()
	envconfig.AddPath("data/simple/simpleArrayConv.env")

	cfg := ConvArrayConfig{}
	envconfig.EnableLog()
	err := envconfig.LoadConfig(&cfg)
	if err != nil {
		assert.Fail(err.Error())
	}

	assert.Equal(3, len(cfg.IntArrayConf))
	assert.Equal(1, cfg.IntArrayConf[0])
	assert.Equal(2, cfg.IntArrayConf[1])
	assert.Equal(3, cfg.IntArrayConf[2])

	assert.Equal(3, len(cfg.FloatArrayConf))
	assert.Equal(1.1, cfg.FloatArrayConf[0])
	assert.Equal(2.2, cfg.FloatArrayConf[1])
	assert.Equal(3.3, cfg.FloatArrayConf[2])

}
