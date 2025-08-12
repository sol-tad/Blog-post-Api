package infrastructure

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes a plaintext password using bcrypt with a cost of 14
// Returns the hashed password as a string, or an error if hashing fails
func HashPassword(password string) (string, error) {
	// GenerateFromPassword returns the bcrypt hash of the password at the given cost
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPassword compares a plaintext password with a hashed password
// Returns true if they match, false otherwise
func CheckPassword(password, hashedPassword string) bool {
	// CompareHashAndPassword returns nil on success, or an error if passwords don't match
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
