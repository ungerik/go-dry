package dry

import (
	cryptoRand "crypto/rand"
	"fmt"
	mathRand "math/rand"
	"time"
)

// RandSeedWithTime calls rand.Seed() with the current time.
func RandSeedWithTime() {
	mathRand.Seed(time.Now().UTC().UnixNano())
}

// RandomHexString returns a random lower case hex string with length.
func RandomHexString(length int) string {
	buffer := make([]byte, length/2)
	_, err := cryptoRand.Read(buffer)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", buffer)
}

// RandomHEXString returns a random upper case hex string with length.
func RandomHEXString(length int) string {
	buffer := make([]byte, length/2)
	_, err := cryptoRand.Read(buffer)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%X", buffer)
}
