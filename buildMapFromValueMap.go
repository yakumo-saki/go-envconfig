package envconfig

import (
	"fmt"
	"reflect"
	"strings"
)

// buildMapFromValueMap make slice[string] from config suffix by _00 ~ _99
func buildMapFromValueMap(prefix string, mapType reflect.Type, valueMap map[string]string) any {

	retRv := reflect.MakeMap(mapType)
	valueType := mapType.Elem()

	keyType := mapType.Key()
	if keyType.Kind() != reflect.String {
		msg := fmt.Sprintf("map is supported. but key must be string. cant use %s", mapType.Kind().String())
		panic(msg)
	}

	switch valueType.Kind() {
	case reflect.Slice:

		fmt.Println("slice inside map")
	case reflect.Map:
		fmt.Println("map inside map")
		panic("map inside map is not supported")
	default:
		fmt.Printf("map key -> %s\n", valueType.String())
		keyMap := scanMapFromValueMap(prefix, valueMap)
		for k, v := range keyMap {
			key := k
			retRv.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(valueMap[v]))
		}
	}

	var ret any = retRv.Interface()

	return ret
}

// scanMapKeysFromValueMap scans value keys then make map
// for slice values
// example (prefix=MAP)
// MAP_GRP1_1, MAP_GRP1_2, MAP_GROUPB_1 ->  map{"GRP1": "MAP_GRP1_", "GROUPB": "MAP_GROUPB_"}
func scanMapPrefixsFromValueMap(prefix string, valueMap map[string]string) map[string]string {

	ret := make(map[string]string)

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
