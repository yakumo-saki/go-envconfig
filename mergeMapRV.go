package envconfig

import (
	"fmt"
	"reflect"
)

// createMergedMapRV creates new map and return as reflect.Value
// if map value is slice and both maps has same key, slices are append.
// if map value is NOT slice and both maps has same key, value of map2 is set
func createMergedMapRV(map1, map2 reflect.Value) (reflect.Value, error) {

	map1ValueType := map1.Type().Elem()
	map2ValueType := map2.Type().Elem()

	if map1ValueType.Kind() != map2ValueType.Kind() {
		retRv := reflect.MakeMap(map1.Type())
		err := fmt.Errorf("map1 and 2 has different types, %s and %s",
			map1ValueType.String(), map2ValueType.String())
		return retRv, err
	}

	l("map1 and 2 types %s %s\n", map1ValueType.String(), map2ValueType.String())

	switch map1ValueType.Kind() {
	case reflect.Slice:
		return createMergedMapRVSlice(map1, map2)
	default:
		return createMergedMapRVNonSlice(map1, map2)
	}
}

func createMergedMapRVNonSlice(map1, map2 reflect.Value) (reflect.Value, error) {
	retRv := reflect.MakeMap(map1.Type())

	iter := map1.MapRange()
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()
		retRv.SetMapIndex(k, v)
	}

	iter = map2.MapRange()
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()
		retRv.SetMapIndex(k, v)
	}

	return retRv, nil
}

func createMergedMapRVSlice(map1, map2 reflect.Value) (reflect.Value, error) {
	retRv := reflect.MakeMap(map1.Type())

	iter := map1.MapRange()
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()
		retRv.SetMapIndex(k, v)
	}

	iter = map2.MapRange()
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()
		retRv.SetMapIndex(k, v)
	}

	return retRv, nil
}
