package envconfig_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yakumo-saki/go-envconfig"
)

type EnvConfig struct {
	StrConf string `cfg:"STR_CONF"`
	IntConf int    `cfg:"INT_CONF"`
	MyConf  string
}

// slice以外の階層化されていない単純なstructの読み込みテスト
// 値の変換にバグがないかテストする
func TestEnvOnly(t *testing.T) {
	assert := assert.New(t)
	envconfig.ClearPath()

	SetEnv("MY_CONF", "MYCONF")
	SetEnv("STR_CONF", "ENV123")
	SetEnv("INT_CONF", "555")

	cfg := EnvConfig{}
	envconfig.EnableLogWithDefaultLogger()
	err := envconfig.LoadConfig(&cfg)
	if err != nil {
		assert.Fail(err.Error())
	}

	fmt.Printf("after %v\n", cfg)
	assert.Nil(err)
	assert.Equal("ENV123", cfg.StrConf)
	assert.Equal(555, cfg.IntConf)
	assert.Equal("MYCONF", cfg.MyConf)
}
