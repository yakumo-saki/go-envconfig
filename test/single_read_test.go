package envconfig_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yakumo-saki/go-envconfig"
)

type SimpleConfig struct {
	StrConf string `cfg:"STR_CONF"`
	IntConf int    `cfg:"INT_CONF"`
}

func TestLoadSimpleOne(t *testing.T) {
	assert := assert.New(t)
	envconfig.ClearPath()
	envconfig.AddPath("data/simple/simple.env")

	cfg := SimpleConfig{}
	err := envconfig.LoadConfig(&cfg)

	fmt.Printf("after %v\n", cfg)
	assert.Nil(err)
	assert.Equal("abc", cfg.StrConf)
	assert.Equal(123, cfg.IntConf)
}
