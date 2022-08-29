package envconfig

import "fmt"

// EnableDebugLog enables logging with specified log output function.
// Not output log if function is Nil
func EnableDebugLog(debug, warn, user func(format string, a ...interface{})) {
	debuglog = debug
	warnlog = warn
}

// EnableLogWithDefaultLogger enables logging with fmt.Printf output.
func EnableLogWithDefaultLogger() {
	userlog = func(format string, a ...interface{}) { fmt.Printf(format, a...) }
	warnlog = func(format string, a ...interface{}) { fmt.Printf("WARN :"+format, a...) }
	debuglog = func(format string, a ...interface{}) { fmt.Printf("DEBUG:"+format, a...) }
}

func logUser(format string, a ...interface{}) {
	userlog(format, a)
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
