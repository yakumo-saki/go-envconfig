package envconfig_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yakumo-saki/go-envconfig"
)

type TestDefaultConfig struct {
	StrConf     string
	Sub         TestDefaultSubConfig
	warningTest string `cfg:"IGNORED"`
}

type TestDefaultSubConfig struct {
	SubStructStrConf string
	SubConfig        string `cfg:"SUB_CONFIG"`
}

// デフォルト値を保存しているかテスト
func TestSimpleDefaultValue(t *testing.T) {
	assert := assert.New(t)

	ec := envconfig.New()
	ec.ClearPath()
	ec.AddPath("data/simple/default_value.env")

	cfg := TestDefaultConfig{StrConf: "DEFAULT_STR_CONF"}
	cfg.Sub = TestDefaultSubConfig{SubStructStrConf: "SUB_STRUCT_DEFAULT", SubConfig: "must be overwrite"}
	ec.EnableLogWithDefaultLogger()
	err := ec.LoadConfig(&cfg)
	if err != nil {
		assert.Fail(err.Error())
	}

	assert.Equal("DEFAULT_STR_CONF", cfg.StrConf)
	assert.Equal("SUB_STRUCT_DEFAULT", cfg.Sub.SubStructStrConf)
	assert.Equal("123abc", cfg.Sub.SubConfig)
	assert.Equal("", cfg.warningTest) // not inject. and output warning log
}
