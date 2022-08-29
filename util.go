package envconfig

import (
	"os"
	"strings"
)

func getOSEnv() map[string]string {
	ret := make(map[string]string)

	for _, v := range os.Environ() {
		splitted := strings.SplitN(v, "=", 2)
		ret[splitted[0]] = splitted[1]
	}
	return ret
}
