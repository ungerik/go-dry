package dry

import (
	"math/rand"
	"time"
)

func RandSeedWithTime() {
	rand.Seed(time.Now().UTC().UnixNano())
}
