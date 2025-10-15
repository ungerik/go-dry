package dry

import (
	"os"
	"strings"
)

// EnvironMap returns the current environment variables as a map.
// Each environment variable is expected to be in "KEY=value" format.
// Variables without an "=" separator are skipped.
func EnvironMap() map[string]string {
	return environToMap(os.Environ())
}

// environToMap converts environment variable strings to a map.
// Each string is split at the first "=" into key and value.
// Strings without "=" are skipped (returns found=false from strings.Cut).
// For variables like "KEY=", the value will be an empty string.
// For variables like "KEY=value1=value2", only splits at the first "=".
func environToMap(environ []string) map[string]string {
	ret := make(map[string]string)

	for _, v := range environ {
		key, value, found := strings.Cut(v, "=")
		if !found {
			continue
		}
		ret[key] = value
	}

	return ret
}

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
