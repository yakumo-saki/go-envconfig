package envconfig

import (
	"fmt"
	"reflect"
)

// createMergedMapRV creates new map and return as reflect.Value
// if map value is slice and both maps has same key, slices are append.
// if map value is NOT slice and both maps has same key, value of map2 is set
func createMergedMapRV(lesserMap, priorMap reflect.Value) (reflect.Value, error) {

	map1ValueType := lesserMap.Type().Elem()
	map2ValueType := priorMap.Type().Elem()

	if map1ValueType.Kind() != map2ValueType.Kind() {
		retRv := reflect.MakeMap(lesserMap.Type())
		err := fmt.Errorf("map1 and 2 has different types, %s and %s",
			map1ValueType.String(), map2ValueType.String())
		return retRv, err
	}

	l("map1 and 2 types %s %s\n", map1ValueType.String(), map2ValueType.String())

	switch map1ValueType.Kind() {
	case reflect.Slice:
		return createMergedMapRVSlice(lesserMap, priorMap)
	default:
		return createMergedMapRVNonSlice(lesserMap, priorMap)
	}
}

func createMergedMapRVNonSlice(lesserMap, priorMap reflect.Value) (reflect.Value, error) {
	retRv := reflect.MakeMap(lesserMap.Type())

	iter := lesserMap.MapRange()
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()
		retRv.SetMapIndex(k, v)
	}

	iter = priorMap.MapRange()
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()
		retRv.SetMapIndex(k, v)
	}

	return retRv, nil
}

func createMergedMapRVSlice(lesserMap, priorMap reflect.Value) (reflect.Value, error) {
	retRv := reflect.MakeMap(lesserMap.Type())

	iter := lesserMap.MapRange()
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()
		retRv.SetMapIndex(k, v)
	}

	iter = priorMap.MapRange()
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()

		l("k,v = %s %s", k, v)

		orgSlice := retRv.MapIndex(k)
		if orgSlice == reflect.Zero(v.Type()) {
			// no need to merge only on priorMap
			retRv.SetMapIndex(k, v)
			continue
		}

		newSlice := reflect.AppendSlice(orgSlice, v)

		retRv.SetMapIndex(k, newSlice)
	}

	ret := retRv.Interface()
	l("%v", ret)

	return retRv, nil
}
