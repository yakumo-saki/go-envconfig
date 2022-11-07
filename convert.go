package envconfig

import (
	"fmt"
	"reflect"
	"strconv"
)

// convertTo converts string value to type
func convertTo(value string, typ reflect.Type) (reflect.Value, error) {

	switch typ.Kind() {
	case reflect.String:
		return reflect.ValueOf(value), nil
	case reflect.Int:
		v, err := strconv.ParseInt(value, 10, 0)
		if err != nil {
			return reflect.ValueOf(nil), createParseError(value, typ.String(), err)
		}
		ret := reflect.ValueOf(v).Convert(typ)
		return ret, nil
	case reflect.Float64:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return reflect.ValueOf(nil), createParseError(value, typ.String(), err)
		}
		ret := reflect.ValueOf(v).Convert(typ)
		return ret, nil
	default:
		return reflect.ValueOf(nil), fmt.Errorf("converTo: unsupported type %s", typ.Kind().String())
	}
}

// convertSlice converts []string to []type
func convertSlice(sliceStr []string, sliceType reflect.Type) (reflect.Value, error) {
	ret := reflect.MakeSlice(reflect.SliceOf(sliceType), 0, len(sliceStr))

	for _, v := range sliceStr {
		refVal, err := convertTo(v, sliceType)
		if err != nil {
			return ret, err
		}

		ret = reflect.Append(ret, refVal)
	}

	return ret, nil
}

func createParseError(value, typeName string, err error) error {
	return fmt.Errorf("value %s parse error. required type %s. %w", value, typeName, err)
}
