package env

import "os"

// GetEnv return env value or default if empty.
func GetEnv(name, defValue string) string {
	value, ok := os.LookupEnv(name) //nolint: forbidigo
	if !ok {
		value = defValue
	}
	return value
}
