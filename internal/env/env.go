package env

import "os"

func GetEnv(name, defValue string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		value = defValue
	}
	return value
}
