package envconfig

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/yakumo-saki/go-envconfig/internal/envfile"
	"github.com/yakumo-saki/go-envconfig/internal/util"
)

// EnvConfig is main struct of EnvConfig.
// Use envconfig.New() to get instance.
type EnvConfig struct {
	userlog LogFunc // Log for user. eg. wrong value, cant convert string to int

	warnlog  LogFunc // Log for developer.
	debuglog LogFunc // Log for developer.

	strict bool     // strictモード（明示的にcfgタグを書かない限りインジェクトしない）
	paths  []string // 設定ファイル検索パス
}

// struct tag. "cfg"
const TAG = "cfg"

// LogFunc is log output function
// parameters are same as fmt.Printf
type LogFunc func(string, ...interface{})

type options struct {
	ConfigKey string
	Slice     bool
	Map       bool
	Merge     bool // slice or map only. true=merge false=destroy and create
}

type reflectField struct {
	Field    reflect.StructField
	RefValue reflect.Value
	Options  options
}

// Always use this method to get instance of EnvConfig
func New() *EnvConfig {
	ec := EnvConfig{}
	return &ec
}

// UseStrict enables strict mode.
// in strict mode, all config field must be tagged with "cfg:"
func (ec *EnvConfig) UseStrict() {
	ec.strict = true
}

// AddPath adds config read path.
//
// every path not needs to be existed.
// If file not found on path, its simply ignored.
//
// Same config are overwrite by last added config file.
// returns path count added
func (ec *EnvConfig) AddPath(path string) int {
	if ec.paths == nil {
		ec.ClearPath()
	}
	ec.paths = append(ec.paths, path)
	return len(ec.paths)
}

// ClearPath clear all paths added by AddPath()
func (ec *EnvConfig) ClearPath() {
	ec.paths = make([]string, 0)
}

// LoadConfig starts config process
// parameter must be pointer of struct entity
func (ec *EnvConfig) LoadConfig(cfg interface{}) error {

	cfgType := reflect.TypeOf(cfg)
	if !strings.HasPrefix(cfgType.String(), "*") {
		return fmt.Errorf("cfg must be pass as pointer, but passed %s", cfgType.String())
	}

	// tag -> fielddesc のmapを作る
	configFieldMap := ec.buildConfigFieldMap(cfg)
	for k, v := range configFieldMap {
		ec.logDebug("configField: key=%s value=%v\n", k, v)
	}

	// load config files
	for _, p := range ec.paths {
		if util.FileExists(p) {
			env, err := envfile.ReadEnvfile(p)
			if err != nil {
				return err
			}
			err = ec.applyEnvMap(env, configFieldMap, cfg)
			if err != nil {
				return err
			}
		} else {
			// debug
			ec.logDebug("FileNotFound %s\n", p)
		}
	}

	// load environment values
	envMap := getOSEnv()
	err := ec.applyEnvMap(envMap, configFieldMap, cfg)
	if err != nil {
		return err
	}

	return nil
}

// buildSliceFromValueMap make slice[string] from config suffix by _00 ~ _99
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

// applyEnvMap apply valueMap to cfg struct, using configFieldMap
func (ec *EnvConfig) applyEnvMap(valueMap map[string]string, configFieldMap map[string]reflectField, cfg interface{}) error {

	// Slice対応があるので configFieldMapから該当する設定
	transformedMap := ec.transformValueMap(valueMap, configFieldMap)

	for key, val := range transformedMap {
		refField, ok := configFieldMap[key]
		if !ok {
			panic("transformedMap key is not found in configFieldMap")
		}

		field := refField.Field
		fieldVal := refField.RefValue
		option := refField.Options

		ec.logDebug("cfgElemValue Name=%s (Value)=%v Type=%s(%s) CanSet=%v\n",
			field.Name, fieldVal, fieldVal.Type(), field.Type.Kind(), fieldVal.CanSet())

		switch fieldVal.Kind() {
		case reflect.Slice:
			// slice needs special care
			sliceStr := val.([]string) // value slice is always string

			// slice of string uses fast path.
			sliceType := field.Type.Elem() // type of slice

			if sliceType.Kind() == reflect.String {
				ec.logDebug("slice of string. use fast path\n")
				if fieldVal.IsNil() || !option.Merge {
					ec.logDebug("slice overwrite\n")
					fieldVal.Set(reflect.ValueOf(sliceStr))
				} else {
					fieldVal.Set(reflect.AppendSlice(fieldVal, reflect.ValueOf(sliceStr)))
				}

				break
			}

			// normal path
			// logDebug("slice of %s. use normal path\n", sliceType)
			newSlice, err := convertSlice(sliceStr, sliceType)
			if err != nil {
				return fmt.Errorf("error on field %s: %w", field.Name, err)
			}
			if fieldVal.IsNil() || !option.Merge {
				ec.logDebug("slice overwrite\n")
				fieldVal.Set(newSlice)
			} else {
				fieldVal.Set(reflect.AppendSlice(fieldVal, newSlice))
			}
		default:
			v, err := convertTo(val.(string), field.Type)
			if err != nil {
				ec.logUser("Can't parse %s as %s (field %s): %w", val.(string), field.Type, field.Name, err)
				return fmt.Errorf("error on field %s: %w", field.Name, err)
			}
			fieldVal.Set(v)
		}
	}

	return nil
}
