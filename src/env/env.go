// Package env holds some helper functions about environment variables
package env

import (
	"fmt"
	"os"
)

// GetDef returns the environment variable named as the argument, or the
// stated default value if the variable does not exist
func GetDef[T any](propertyName string, defaultValue T) T {
	value, found := os.LookupEnv(propertyName)
	if found == false {
		return defaultValue
	}
	var val T
	if _, err := fmt.Sscan(value, &val); err != nil {
		panic("error parsing " + propertyName + "=" + value + " -> " + err.Error())
	}
	return val
}
