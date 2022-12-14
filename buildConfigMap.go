package envconfig

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/stoewer/go-strcase"
)

func (ec *EnvConfig) buildConfigFieldMap(cfg interface{}) map[string]reflectField {
	cfgValue := reflect.ValueOf(cfg)
	return ec.buildConfigFieldMapImpl(cfgValue)
}

// buildConfigFieldMap build config read strategy
// param cfgValue reflect.Value of pointer to struct instance
func (ec *EnvConfig) buildConfigFieldMapImpl(refValOfPtrStruct reflect.Value) map[string]reflectField {
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
			_, ok := field.Tag.Lookup(tag)
			if ok {
				ec.logWarn("field %s has %s tag. but ignored. because field %s is not exported.\n",
					field.Name, tag, field.Name)
			}

			continue // 非公開フィールドは対象外
		}

		// Structは再帰処理が必要。Structに更に潜る
		if field.Type.Kind() == reflect.Struct {
			// Structのポインタを渡さないと書き換えができないのでポインタを得る
			structCfg := ec.buildConfigFieldMapImpl(fieldEntity.Addr())

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
		tag, hasTag := field.Tag.Lookup(tag)
		if ec.strict && !hasTag {
			continue // strictモードではcfgタグがついていないフィールドは無視
		}

		// build options(CFG envname, type etc) from tag and field
		opt := createTagFromFieldDef(field)
		if hasTag {
			processfield := parseTag(&opt, field, tag)
			if !processfield {
				continue // cfg:"" が指定されているのでフィールドを無視
			}
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
	opt := options{Slice: false, Map: false}

	name := field.Name
	opt.ConfigKey = strcase.UpperSnakeCase(name)

	switch field.Type.Kind() {
	case reflect.Slice:
		opt.Slice = true
		opt.SliceMerge = true
	case reflect.Map:
		opt.Map = true
		opt.MapMergeType = KeyMerge
	default:
		// no need to do
	}

	return opt
}

// parseTag parses `cfg:"xxxxx, ttttt, ooooo"` from struct
// result modifying opt
// @return false=ignore this field, true=process this field
func parseTag(opt *options, field reflect.StructField, tagString string) bool {

	if tagString == "" {
		return false
	}

	splitted := strings.Split(tagString, ",")
	switch {
	case len(splitted) == 1:
		opt.ConfigKey = splitted[0]
	case len(splitted) == 2:
		opt.ConfigKey = splitted[0]

		switch field.Type.Kind() {
		case reflect.Slice:
			switch strings.ToLower(splitted[1]) {
			case "overwrite":
				opt.SliceMerge = false
			case "merge":
				opt.SliceMerge = true
			default:
				msg := fmt.Sprintf("'cfg: %s, %s' is not reconized or not applicable to slice.", splitted[0], splitted[1])
				panic(msg)
			}
		case reflect.Map:
			switch strings.ToLower(splitted[1]) {
			case "overwrite":
				opt.MapMergeType = OverwriteAll
			case "keymerge":
				opt.MapMergeType = KeyMerge
			case "valuemerge":
				opt.MapMergeType = ValueMerge
			default:
				msg := fmt.Sprintf("'cfg: %s, %s' is not reconized or not applicable to map.", splitted[0], splitted[1])
				panic(msg)
			}
		default:
			msg := fmt.Sprintf("'cfg: %s, %s' is only suitable for slice or map", splitted[0], splitted[1])
			panic(msg)
		}
	default:
		msg := fmt.Sprintf("Illegal cfg tag %s on %s", tagString, field.Name)
		panic(msg)
	}
	return true
}

// transformValueMap transforms valueMap to configMap
// Slice化の処理を行う
// valueMap map[string]string -> map[string]string|[]string
func (ec *EnvConfig) transformValueMap(valueMap map[string]string, configFieldMap map[string]reflectField) (map[string]interface{}, error) {
	ret := make(map[string]interface{})

	for key, refField := range configFieldMap {
		ec.logDebug("key %s refF=%v\n", key, refField)

		cfgKey := refField.Options.ConfigKey
		switch {
		case refField.Options.Slice:
			slice := buildSliceFromValueMap(cfgKey, valueMap)
			if len(slice) == 0 {
				continue // sliceがempty => no config
			}
			ret[cfgKey] = slice
		case refField.Options.Map:
			// determine value type of map (key is always string)
			m, err := buildMapFromValueMap(cfgKey, refField.Field.Type, valueMap)
			if err != nil {
				return ret, err
			}
			ret[cfgKey] = m
		default:
			v, ok := valueMap[cfgKey]
			if !ok {
				continue
			}
			ret[cfgKey] = v
		}
	}

	return ret, nil
}
