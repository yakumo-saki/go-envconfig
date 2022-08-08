package envconfig

import "fmt"

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
