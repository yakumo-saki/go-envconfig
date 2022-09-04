package envconfig_test

import "os"

func SetEnv(key, value string) {
	err := os.Setenv(key, value)
	if err != nil {
		panic(err)
	}
}
