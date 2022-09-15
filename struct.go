package envconfig

import "reflect"

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
