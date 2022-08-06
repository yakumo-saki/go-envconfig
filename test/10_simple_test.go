package envconfig_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yakumo-saki/go-envconfig"
)

type SimpleConfig struct {
	StrConf         string `cfg:"STR_CONF"`
	IntConf         int    `cfg:"INT_CONF"`
	TestConvertCase string
}

// slice以外の階層化されていない単純なstructの読み込みテスト
// 値の変換にバグがないかテストする
func TestLoadSimpleOne(t *testing.T) {
	assert := assert.New(t)
	envconfig.ClearPath()
	envconfig.AddPath("data/simple/simple.env")

	cfg := SimpleConfig{}
	envconfig.EnableLog()
	err := envconfig.LoadConfig(&cfg)
	if err != nil {
		assert.Fail(err.Error())
	}

	fmt.Printf("after %v\n", cfg)
	assert.Nil(err)
	assert.Equal("abc", cfg.StrConf)
	assert.Equal(123, cfg.IntConf)
	assert.Equal("this must be read", cfg.TestConvertCase)
}
