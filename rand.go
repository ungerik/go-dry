package dry

import (
	cryptorand "crypto/rand"
	"fmt"
)

// RandSeedWithTime is deprecated since Go 1.20.
// The math/rand package now automatically seeds itself.
// This function is kept for backward compatibility but does nothing.
//
// Deprecated: No longer needed. Go 1.20+ automatically seeds math/rand.
func RandSeedWithTime() {
	// No-op: Go 1.20+ automatically seeds math/rand
}

func getRandomHexString(length int, formatStr string) string {
	var buffer []byte
	if length%2 == 0 {
		buffer = make([]byte, length/2)
	} else {
		buffer = make([]byte, (length+1)/2)
	}
	_, err := cryptorand.Read(buffer)
	if err != nil {
		return ""
	}
	hexString := fmt.Sprintf(formatStr, buffer)
	return hexString[:length]
}

// RandomHexString returns a random lower case hex string with length.
func RandomHexString(length int) string {
	return getRandomHexString(length, "%x")
}

// RandomHEXString returns a random upper case hex string with length.
func RandomHEXString(length int) string {
	return getRandomHexString(length, "%X")
}
