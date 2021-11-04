// `hasher` package provides common functions for working with passwords.
package hasher

import "golang.org/x/crypto/bcrypt"

// HashPassword returns the bcrypt hash of the password at the minimum cost.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hash), err
}

// CheckPasswordHash compares a plaintext password with it's possible bcrypt hash equivalent.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
