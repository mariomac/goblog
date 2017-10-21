package env

import "os"

func Get(propertyName string) string {
	return GetDef(propertyName, "")
}

func GetDef(propertyName string, defaultValue string) string {
	value, found := os.LookupEnv(propertyName)
	if found == false {
		return defaultValue
	} else {
		return value
	}
}
