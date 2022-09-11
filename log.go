package envconfig

import "fmt"

// EnableUserLog enables log message for users.
func (ec *EnvConfig) EnableUserLog(user func(format string, a ...interface{})) {
	ec.userlog = user
}

// EnableDebugLog enables logging with specified log output function.
// Not output log if function is Nil
func (ec *EnvConfig) EnableDebugLog(debug, warn func(format string, a ...interface{})) {
	ec.debuglog = debug
	ec.warnlog = warn
}

// EnableLogWithDefaultLogger enables logging with fmt.Printf output.
func (ec *EnvConfig) EnableLogWithDefaultLogger() {
	ec.userlog = func(format string, a ...interface{}) { fmt.Printf(format, a...) }
	ec.warnlog = func(format string, a ...interface{}) { fmt.Printf("WARN :"+format, a...) }
	ec.debuglog = func(format string, a ...interface{}) { fmt.Printf("DEBUG:"+format, a...) }
}

func (ec *EnvConfig) logUser(format string, a ...interface{}) {
	ec.userlog(format, a)
}

func (ec *EnvConfig) logWarn(format string, a ...interface{}) {
	if ec.debuglog == nil {
		return
	}
	ec.warnlog(format, a...)
}

func (ec *EnvConfig) logDebug(format string, a ...interface{}) {
	if ec.debuglog == nil {
		return
	}
	ec.debuglog(format, a...)
}
