package stdgomods

import (
	"os"
	"strconv"
)

// If an ENV variable is set, return that value. Otherwise, return the
// default
func GetEnvVar(envName string, defaultValue string) string {
	ENV_VAL := os.Getenv(envName)
	if ENV_VAL == "" {
		return defaultValue
	}
	return ENV_VAL
}

// Sometimes we want to use boolean ENV variables.
func GetBoolEnvVar(envName string, defaultValue bool) bool {
	ENV_VAL := os.Getenv(envName)
	if ENV_VAL == "" {
		return defaultValue
	}
	bool_val, err := strconv.ParseBool(ENV_VAL)
	PanicIfError(err, "Could not parse env value into boolean")
	return bool_val
}

// Sometimes we want to use numeric ENV variables.
func GetNumericEnvVar(envName string, defaultValue int) int {
	ENV_VAL := os.Getenv(envName)
	if ENV_VAL == "" {
		return defaultValue
	}
	retInt, err := strconv.Atoi(ENV_VAL)
	if err != nil {
		fmt.Printf("WARNING: The environment variable for '%s' is '%s'. I could not convert that to a number. Using default value of %d\n", envName, ENV_VAL, defaultValue)
		return defaultValue
	}
	return retInt
}
