package envconfig

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/stoewer/go-strcase"

	"github.com/yakumo-saki/go-envconfig/internal/envfile"
	"github.com/yakumo-saki/go-envconfig/internal/util"
)

const TAG = "cfg"

type LogFunc func(string, ...interface{})

var warnlog LogFunc
var debuglog LogFunc

var paths []string

// strictモード（明示的にcfgタグを書かない限りインジェクトしない）
var strict bool = false

type Options struct {
	ConfigKey  string
	Slice      bool
	SliceMerge bool
}

type ReflectField struct {
	Field    reflect.StructField
	RefValue reflect.Value
	Options  Options
}

// UseStrict enables strict mode.
// in strict mode, all config field must be tagged with "cfg:"
func UseStrict() {
	strict = true
}

// AddPath adds config read path.
//
// every path not needs to be existed.
// If file not found on path, its simply ignored.
//
// Same config are overwrite by last added config file.
// returns path count added
func AddPath(path string) int {
	if paths == nil {
		ClearPath()
	}
	paths = append(paths, path)
	return len(paths)
}

func ClearPath() {
	paths = make([]string, 0)
}

func LoadConfig(cfg interface{}) error {

	cfgType := reflect.TypeOf(cfg)
	if !strings.HasPrefix(cfgType.String(), "*") {
		return fmt.Errorf("cfg must be pass as pointer, but passed %s", cfgType.String())
	}

	// tag -> fielddesc のmapを作る
	configFieldMap := buildConfigFieldMap(cfg)
	for k, v := range configFieldMap {
		logDebug("configField: key=%s value=%v\n", k, v)
	}

	// load config files
	for _, p := range paths {
		if util.FileExists(p) {
			env, err := envfile.ReadEnvfile(p)
			if err != nil {
				return err
			}
			err = applyEnvfile(env, configFieldMap, cfg)
			if err != nil {
				return err
			}
		} else {
			// debug
			logDebug("FileNotFound %s\n", p)
		}
	}

	// load environment values

	return nil
}

func buildConfigFieldMap(cfg interface{}) map[string]ReflectField {
	cfgValue := reflect.ValueOf(cfg)
	return buildConfigFieldMapImpl(cfgValue)
}

// buildConfigFieldMap build config read strategy
// param cfgValue reflect.Value of pointer to struct instance
func buildConfigFieldMapImpl(refValOfPtrStruct reflect.Value) map[string]ReflectField {
	ret := make(map[string]ReflectField)

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
		tag, ok := field.Tag.Lookup(TAG)
		if strict && !ok {
			continue // strictモードではcfgタグがついていないフィールドは無視
		}

		opt := parseTag(field.Name, tag)
		// log("field=%s key=%s (tag %s)\n", field.Name, opt.ConfigKey, tag)
		preexist, ok := ret[opt.ConfigKey]
		if ok {
			panic(fmt.Sprintf("config key duplicated: %s. first is %s, second is %s",
				opt.ConfigKey, preexist.Field.Name, field.Name))
		}
		ret[opt.ConfigKey] = ReflectField{Field: field, RefValue: fieldEntity, Options: opt}
		if opt.Slice && field.Type.Kind() != reflect.Slice {
			panic(fmt.Sprintf("struct field %s is defined as slice(by tag), but not slice. check field definition", field.Name))
		}
	}

	return ret
}

// parseTag parses `cfg:""` from struct
func parseTag(fieldName, tagString string) Options {
	opt := Options{Slice: false, SliceMerge: false}
	opt.ConfigKey = strcase.UpperSnakeCase(fieldName)

	if tagString == "" {
		return opt
	}

	splitted := strings.Split(tagString, ",")
	switch {
	case len(splitted) == 1:
		opt.ConfigKey = splitted[0]
	case len(splitted) == 2:
		opt.ConfigKey = splitted[0]
		if strings.EqualFold(splitted[1], "slice") {
			opt.Slice = true
			opt.SliceMerge = false
		} else if strings.EqualFold(splitted[1], "mergeslice") {
			opt.Slice = true
			opt.SliceMerge = true
		}
	default:
		msg := fmt.Sprintf("Illegal cfg tag %s on %s", tagString, fieldName)
		panic(msg)
	}

	return opt
}

// transformValueMap transforms valueMap to configMap
// Slice化の処理を行う
func transformValueMap(valueMap map[string]string, configFieldMap map[string]ReflectField) map[string]interface{} {
	ret := make(map[string]interface{})

	for key, refField := range configFieldMap {
		logDebug("key %s refF=%v\n", key, refField)

		cfgKey := refField.Options.ConfigKey
		if refField.Options.Slice {
			slice := buildSliceFromValueMap(cfgKey, valueMap)
			ret[cfgKey] = slice
		} else {
			v, ok := valueMap[cfgKey]
			if ok {
				ret[cfgKey] = v
			}
		}

	}

	return ret
}

func buildSliceFromValueMap(prefix string, valueMap map[string]string) []string {
	ret := make([]string, 0)
	for i := 0; i < 100; i++ {
		v, ok := valueMap[prefix+fmt.Sprint(i)]
		if ok {
			ret = append(ret, v)
		}

		// 0〜9 の場合 00~09 もチェック
		if i < 10 {
			v, ok := valueMap[fmt.Sprintf("%s%02d", prefix, i)]
			if ok {
				ret = append(ret, v)
			}
		}
	}
	return ret
}

func applyEnvfile(valueMap map[string]string, configFieldMap map[string]ReflectField, cfg interface{}) error {

	// Slice対応があるので configFieldMapから該当する設定
	transformedMap := transformValueMap(valueMap, configFieldMap)

	for key, val := range transformedMap {
		refField, ok := configFieldMap[key]
		if !ok {
			panic("transformedMap key is not found in configFieldMap")
		}

		field := refField.Field
		fieldVal := refField.RefValue
		option := refField.Options

		logDebug("cfgElemValue Name=%s (Value)=%v Type=%s(%s) CanSet=%v\n",
			field.Name, fieldVal, fieldVal.Type(), field.Type.Kind(), fieldVal.CanSet())
		// ここ、stringのvalをそれぞれの型に変換しないといけないんだけども、全部の型についてなにか書かないと駄目？
		// せっかくリフレクションしてるのでどうにか楽したい

		switch fieldVal.Kind() {
		case reflect.Slice:
			// slice needs special care
			sliceStr := val.([]string) // value slice is always string

			// slice of string uses fast path.
			sliceType := field.Type.Elem() // type of slice

			if sliceType.Kind() == reflect.String {
				logDebug("slice of string. use fast path\n")
				if fieldVal.IsNil() || !option.SliceMerge {
					logDebug("slice overwrite\n")
					fieldVal.Set(reflect.ValueOf(sliceStr))
				} else {
					fieldVal.Set(reflect.AppendSlice(fieldVal, reflect.ValueOf(sliceStr)))
				}

				break
			}

			// normal path
			logDebug("slice of %s. use normal path\n", sliceType)
			newSlice, err := convertSlice(sliceStr, sliceType)
			if err != nil {
				return fmt.Errorf("error on field %s: %w", field.Name, err)
			}
			if fieldVal.IsNil() || !option.SliceMerge {
				logDebug("slice overwrite\n")
				fieldVal.Set(newSlice)
			} else {
				fieldVal.Set(reflect.AppendSlice(fieldVal, newSlice))
			}
		default:
			v, err := convertTo(val.(string), field.Type)
			if err != nil {
				return fmt.Errorf("error on field %s: %w", field.Name, err)
			}
			fieldVal.Set(v)
		}
	}

	return nil
}

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
		return reflect.ValueOf(nil), fmt.Errorf("unsupported type %s", typ.Kind().String())
	}
}

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

// EnableLog enables logging with specified log output function.
// Not output log if function is Nil
func EnableLog(debug, warn func(format string, a ...interface{})) {
	debuglog = debug
	warnlog = warn
}

// EnableLogWithDefaultLogger enables logging with fmt.Printf output.
func EnableLogWithDefaultLogger() {
	warnlog = func(format string, a ...interface{}) { fmt.Printf("WARN :"+format, a...) }
	debuglog = func(format string, a ...interface{}) { fmt.Printf("DEBUG:"+format, a...) }
}

func logWarn(format string, a ...interface{}) {
	if debuglog == nil {
		return
	}
	warnlog(format, a...)
}

func logDebug(format string, a ...interface{}) {
	if debuglog == nil {
		return
	}
	debuglog(format, a...)
}
