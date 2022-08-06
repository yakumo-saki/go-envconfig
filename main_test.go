package envconfig

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertSlice(t *testing.T) {
	assert := assert.New(t)

	EnableLog()

	strSlice := []string{"1", "2", "3"}
	intSlice := []int{}

	sliceType := reflect.TypeOf(intSlice).Elem()
	retRefVal, err := convertSlice(strSlice, sliceType)

	assert.Nil(err)

	ret := retRefVal.Interface().([]int)
	assert.Equal(3, len(ret))
	assert.Equal(1, ret[0])
	assert.Equal(2, ret[1])
	assert.Equal(3, ret[2])
}
