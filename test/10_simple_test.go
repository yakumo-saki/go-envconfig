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
	ec := envconfig.New()
	ec.ClearPath()
	ec.AddPath("data/simple/simple.env")

	cfg := SimpleConfig{}
	ec.EnableLogWithDefaultLogger()
	err := ec.LoadConfig(&cfg)
	if err != nil {
		assert.Fail(err.Error())
	}

	fmt.Printf("after %v\n", cfg)
	assert.Nil(err)
	assert.Equal("abc", cfg.StrConf)
	assert.Equal(123, cfg.IntConf)
	assert.Equal("this must be read", cfg.TestConvertCase)
}

// Strictモードではcfgタグがない場合読み込まない
func TestStrictMode(t *testing.T) {
	assert := assert.New(t)
	ec := envconfig.New()
	ec.ClearPath()
	ec.AddPath("data/simple/simple.env")

	cfg := SimpleConfig{}
	ec.UseStrict()
	ec.EnableLogWithDefaultLogger()
	err := ec.LoadConfig(&cfg)
	if err != nil {
		assert.Fail(err.Error())
	}

	assert.Nil(err)
	assert.Equal("", cfg.TestConvertCase)
}
