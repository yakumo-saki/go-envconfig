package envconfig

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/stoewer/go-strcase"
)

func buildConfigFieldMap(cfg interface{}) map[string]reflectField {
	cfgValue := reflect.ValueOf(cfg)
	return buildConfigFieldMapImpl(cfgValue)
}

// buildConfigFieldMap build config read strategy
// param cfgValue reflect.Value of pointer to struct instance
func buildConfigFieldMapImpl(refValOfPtrStruct reflect.Value) map[string]reflectField {
	ret := make(map[string]reflectField)

	structEntity := refValOfPtrStruct.Elem()
	structType := structEntity.Type()
	if structEntity.Kind() != reflect.Struct {
		panic("cfg is not struct")
	}

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldEntity := structEntity.FieldByName(field.Name)

		if !field.IsExported() {
			_, ok := field.Tag.Lookup(TAG)
			if ok {
				logWarn("field %s has %s tag. but ignored. because field %s is not exported.\n",
					field.Name, TAG, field.Name)
			}

			continue // 非公開フィールドは対象外
		}

		// Structは再帰処理が必要。Structに更に潜る
		if field.Type.Kind() == reflect.Struct {
			// Structのポインタを渡さないと書き換えができないのでポインタを得る
			structCfg := buildConfigFieldMapImpl(fieldEntity.Addr())

			for k, v := range structCfg {
				preexist, ok := ret[k]
				if ok {
					panic(fmt.Sprintf("config key duplicated: %s. first is %s, second is %s",
						k, preexist.Field.Name, field.Name))
				}
				ret[k] = v
			}
			continue // Struct自体に直接設定を入れることはできないので無視
		}

		// フィールド
		tag, hasTag := field.Tag.Lookup(TAG)
		if strict && !hasTag {
			continue // strictモードではcfgタグがついていないフィールドは無視
		}

		// build options(CFG envname, type etc) from tag and field
		opt := createTagFromFieldDef(field)
		if hasTag {
			parseTag(&opt, field.Name, tag)
		}

		// log("field=%s key=%s (tag %s)\n", field.Name, opt.ConfigKey, tag)
		preexist, alreadyHasKey := ret[opt.ConfigKey]
		if alreadyHasKey {
			panic(fmt.Sprintf("config key duplicated: %s. first is %s, second is %s",
				opt.ConfigKey, preexist.Field.Name, field.Name))
		}
		ret[opt.ConfigKey] = reflectField{Field: field, RefValue: fieldEntity, Options: opt}
	}

	return ret
}

// default options from struct field definition
func createTagFromFieldDef(field reflect.StructField) options {
	opt := options{Slice: false, Merge: false}

	name := field.Name
	opt.ConfigKey = strcase.UpperSnakeCase(name)

	switch field.Type.Kind() {
	case reflect.Slice:
		opt.Slice = true
	case reflect.Map:
		opt.Map = true
		opt.Merge = true
	default:
		// no need to do
	}

	return opt
}

// parseTag parses `cfg:"xxxxx, ttttt, ooooo"` from struct
// result modifying opt
func parseTag(opt *options, fieldName, tagString string) {

	if tagString == "" {
		return
	}

	splitted := strings.Split(tagString, ",")
	switch {
	case len(splitted) == 1:
		opt.ConfigKey = splitted[0]
	case len(splitted) == 2:
		opt.ConfigKey = splitted[0]
		switch {
		case strings.EqualFold(splitted[1], "slice"):
			panic("'cfg: name, slice' is deprecated. use 'cfg: name, overwrite'. ")
		case strings.EqualFold(splitted[1], "mergeslice"):
			panic("'cfg: name, slice' is deprecated. use 'cfg: name' (default behavior is merge). ")
		case strings.EqualFold(splitted[1], "overwrite"):

		case strings.EqualFold(splitted[1], "merge"):
		}
	default:
		msg := fmt.Sprintf("Illegal cfg tag %s on %s", tagString, fieldName)
		panic(msg)
	}
}

// transformValueMap transforms valueMap to configMap
// Slice化の処理を行う
// valueMap map[string]string -> map[string]string|[]string
func transformValueMap(valueMap map[string]string, configFieldMap map[string]reflectField) map[string]interface{} {
	ret := make(map[string]interface{})

	for key, refField := range configFieldMap {
		logDebug("key %s refF=%v\n", key, refField)

		cfgKey := refField.Options.ConfigKey
		switch {
		case refField.Options.Slice:
			slice := buildSliceFromValueMap(cfgKey, valueMap)
			if len(slice) == 0 {
				continue // sliceがempty => no config
			}
			ret[cfgKey] = slice
		default:
			v, ok := valueMap[cfgKey]
			if !ok {
				continue
			}
			ret[cfgKey] = v
		}
	}

	return ret
}
