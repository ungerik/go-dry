package dry

import "os"

// GetenvDefault retrieves the value of the environment variable
// named by the key. It returns the given defaultValue if the
// variable is not present.
func GetenvDefault(key, defaultValue string) string {
	ret := os.Getenv(key)
	if ret == "" {
		return defaultValue
	}

	return ret
}
