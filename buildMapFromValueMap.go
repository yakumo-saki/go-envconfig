package envconfig

import (
	"fmt"
	"reflect"
	"strings"
)

// buildMapFromValueMap makes map of map type settings for specified single config field
// map Key must be string. map value must not be map, map[string][]type is accepted.
func buildMapFromValueMap(prefix string, mapType reflect.Type, valueMap map[string]string) (any, error) {

	retRv := reflect.MakeMap(mapType)
	valueType := mapType.Elem()

	keyType := mapType.Key()
	if keyType.Kind() != reflect.String {
		msg := fmt.Sprintf("map is supported. but key must be string. cant use %s", mapType.Kind().String())
		panic(msg)
	}

	// process map value
	switch valueType.Kind() {
	case reflect.Map:
		return nil, fmt.Errorf("map inside map is not supported")
	case reflect.Slice:
		sliceValueType := valueType.Elem()
		l("slice inside map. slice of %s\n", sliceValueType.String())

		keyMap := scanMapPrefixesFromValueMap(prefix, valueMap)

		for k, envkey := range keyMap {
			key := k
			strslice := buildSliceFromValueMap(envkey, valueMap)

			rv, err := convertSlice(strslice, sliceValueType)
			if err != nil {
				err = fmt.Errorf("CONVERT ERROR %w", err)
				return nil, err
			}
			retRv.SetMapIndex(reflect.ValueOf(key), rv)
		}

	default:
		l("map value type -> %s\n", valueType.String())
		keyMap := scanMapFromValueMap(prefix, valueMap)

		for k, v := range keyMap {
			key := k
			rv, err := convertTo(valueMap[v], valueType)
			if err != nil {
				err = fmt.Errorf("CONVERT ERROR %w", err)
				return nil, err
			}
			retRv.SetMapIndex(reflect.ValueOf(key), rv)
		}
	}

	var ret any = retRv.Interface()

	return ret, nil
}

// scanMapKeysFromValueMap scans value keys then make map for specified single config field
// for slice values
// example (prefix=MAP)
// MAP_GRP1_1, MAP_GRP1_2, MAP_GROUPB_1 ->  map{"GRP1": "MAP_GRP1_", "GROUPB": "MAP_GROUPB_"}
func scanMapPrefixesFromValueMap(prefix string, valueMap map[string]string) map[string]string {

	ret := make(map[string]string)

	for envname := range valueMap {
		if !strings.HasPrefix(envname, prefix) {
			continue
		}

		// MAP_GRP1_1 -> GRP1_1
		mapKeyWithIndex := strings.Replace(envname, prefix, "", 1)

		// GRP1_1 -> GRP1 is config map key
		keyIndex := strings.LastIndex(mapKeyWithIndex, "_")
		mapKey := string([]rune(mapKeyWithIndex)[:keyIndex])

		// GRP1's environment value name prefix is MAP_GRP1_
		valueIndex := strings.LastIndex(envname, "_")
		mapValuePrefix := string([]rune(envname)[:valueIndex+1])

		ret[mapKey] = mapValuePrefix
	}

	return ret
}

// scanMapKeysFromValueMap scans value keys then make map
// for one value
// example (prefix=VALMAP)
// VALMAP_GRP1, VALMAP_GRP2, VALMAP_GROUPA ->  map{"GRP1": "MAP_GRP1", "GROUPB": "MAP_GROUPB"}
func scanMapFromValueMap(prefix string, valueMap map[string]string) map[string]string {

	ret := make(map[string]string)
	for k := range valueMap {
		if !strings.HasPrefix(k, prefix) {
			continue
		}

		mapKey := strings.Replace(k, prefix, "", 1)
		ret[mapKey] = k
	}

	return ret
}
