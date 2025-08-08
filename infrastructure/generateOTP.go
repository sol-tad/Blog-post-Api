package infrastructure

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateOTP generates a 6-digit random one-time password (OTP) as a string
func GenerateOTP() string {
	// Seed the random number generator with current time in nanoseconds
	// to ensure different results on each call
	rand.Seed(time.Now().UnixNano())
	
	// Generate a random integer between 0 and 999999 (inclusive)
	// and format it as a zero-padded 6-digit string (e.g., "004321")
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
