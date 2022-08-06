package envconfig_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yakumo-saki/go-envconfig"
)

func TestAddPath(t *testing.T) {
	assert := assert.New(t)
	ret := envconfig.AddPath("aaa")
	assert.Equal(ret, 1)
	ret = envconfig.AddPath("aaa")
	assert.Equal(ret, 2)

	envconfig.ClearPath()
}
