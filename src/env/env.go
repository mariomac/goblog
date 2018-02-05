// Package env holds some helper functions about environment variables
package env

import "os"

// GetDef returns the environment variable named as the argument, or the
// stated default value if the variable does not exist
func GetDef(propertyName string, defaultValue string) string {
	value, found := os.LookupEnv(propertyName)
	if found == false {
		return defaultValue
	}
	return value
}
