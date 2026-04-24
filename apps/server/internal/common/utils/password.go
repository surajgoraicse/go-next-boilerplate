package utils

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	// DefaultCost is the default bcrypt cost for password hashing
	DefaultCost = bcrypt.DefaultCost
)

// HashPassword generates a bcrypt hash of the password
// Example usage:
//
//	hashedPassword, err := HashPassword("mypassword")
//	if err != nil {
//	    return err
//	}
//	// Store hashedPassword in database
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPassword compares a bcrypt hashed password with its possible plaintext equivalent
// Returns true if the password matches the hash
// Example usage:
//
//	isValid := CheckPassword("mypassword", hashedPasswordFromDB)
//	if !isValid {
//	    return errors.New("invalid password")
//	}
func CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
