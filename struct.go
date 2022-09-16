package envconfig

import "reflect"

// LogFunc is log output function
// parameters are same as fmt.Printf
type LogFunc func(string, ...interface{})

type options struct {
	ConfigKey    string
	Slice        bool
	SliceMerge   bool // slice only. true=merge false=destroy and create
	Map          bool
	MapMergeType mapMergeType
}

type reflectField struct {
	Field    reflect.StructField
	RefValue reflect.Value
	Options  options
}

type mapMergeType int

const (
	OverwriteAll mapMergeType = iota
	KeyMerge
	ValueMerge
)
