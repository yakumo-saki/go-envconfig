package envconfig

import (
	"fmt"
	"strings"
)

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
// No dont need to call this. envconfig.New() is call this.
func (ec *EnvConfig) EnableLogWithDefaultLogger() {
	ec.userlog = func(format string, a ...interface{}) { fmt.Printf(format, a...) }
	ec.warnlog = nil
	ec.debuglog = nil
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

// returns logDebug output is set
func (ec *EnvConfig) isLogDebugEnabled() bool {
	return ec.debuglog != nil
}

// l is for internal logging
func l(format string, a ...interface{}) {

	if !debugLog {
		return
	}
	if !strings.HasSuffix(format, "\n") {
		format = format + "\n"
	}
	fmt.Printf("INTERNAL:"+format, a...)
}
