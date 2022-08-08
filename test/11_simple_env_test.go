package envconfig_test

import (
	"fmt"
	"os"
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

	os.Setenv("STR_CONF", "ENV123")
	os.Setenv("INT_CONF", "555")
	os.Setenv("MY_CONF", "MYCONF")

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
