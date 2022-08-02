package envconfig

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/stoewer/go-strcase"
	"github.com/yakumo-saki/go-envconfig/internal/util"
)

var paths []string

type Options struct {
	FieldNames []string
	Slice      bool
	SliceMerge bool
}

type ReflectField struct {
	Field    reflect.StructField
	RefValue reflect.Value
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
	if len(paths) == 0 {
		return fmt.Errorf("no paths are added")
	}

	cfgType := reflect.TypeOf(cfg)
	fmt.Printf("%v\n", cfgType)
	if !strings.HasPrefix(cfgType.String(), "*") {
		return fmt.Errorf("cfg must be pass as pointer, but passed %s", cfgType.String())
	}

	// tag -> fielddesc のmapを作る
	configFieldMap := buildConfigFieldMap(cfg)
	for k, v := range configFieldMap {
		fmt.Printf("%s %v", k, v)
	}

	for _, p := range paths {
		if util.FileExists(p) {
			applyEnvfile(p, cfg)
		} else {
			// debug
			fmt.Printf("FNTF %s\n", p)
		}
	}

	return nil
}

func buildConfigFieldMap(cfg interface{}) map[string]ReflectField {
	ret := make(map[string]ReflectField)

	cfgValue := reflect.ValueOf(cfg)
	cfgElemValue := cfgValue.Elem()
	cfgType := cfgElemValue.Type()
	if cfgElemValue.Kind() != reflect.Struct {
		panic("cfg is not struct")
	}

	fmt.Println("loop")
	for i := 0; i < cfgType.NumField(); i++ {
		field := cfgType.Field(i)
		fieldVal := cfgElemValue.FieldByName(field.Name)

		fmt.Printf("cfgElemValue Name=%s (Value)=%v Type=%s CanSet=%v\n", field.Name, fieldVal, fieldVal.Type(), fieldVal.CanSet())
		switch field.Type.Kind() {
		case reflect.String:
			fieldVal.SetString(fmt.Sprintf("replaced by reflection. %s", fieldVal))
		case reflect.Int:
			fieldVal.SetInt(10000 + fieldVal.Int())
		case reflect.Slice:
			fmt.Println("slice")
			strs := []string{"add", "by", "reflection"}
			if fieldVal.IsNil() {
				fieldVal.Set(reflect.ValueOf(strs))
			} else {
				fieldVal.Set(reflect.AppendSlice(fieldVal, reflect.ValueOf(strs)))
			}
		}

		if !field.IsExported() {
			continue // 非公開フィールドは対象外
		}

		tag := field.Tag.Get("cfg")
		fmt.Printf("tag %s\n", tag)
		ret["abc"] = ReflectField{Field: field, RefValue: fieldVal}

	}

	return ret
}

func parseTag(tagString string) {
	strcase.UpperSnakeCase(tagString)
}

func applyEnvfile(path string, cfg interface{}) error {

	env, err := loadEnvfile(path)
	if err != nil {
		return err
	}

	// apply

	fmt.Println(env)
	return nil
}

func loadEnvfile(path string) (map[string]string, error) {
	ret := make(map[string]string)
	filebyte, err := ioutil.ReadFile(path)
	if err != nil {
		return ret, err
	}

	lines := strings.Split(string(filebyte), "\n")
	for _, lineStr := range lines {
		if len(lineStr) == 0 {
			continue
		}

		line := util.StringTrim(lineStr)
		if util.IsCommentLine(line) {
			continue
		}

		fmt.Println(line)
		sep := strings.SplitN(line, "=", 2)
		if len(sep) != 2 {
			fmt.Printf("separeted count %d != 2. %s", len(sep), line)
		}
		ret[sep[0]] = sep[1]
	}

	return ret, nil
}
