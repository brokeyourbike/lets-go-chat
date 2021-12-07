package hasher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	password := "super-secret-password"
	hashed, err := HashPassword(password)

	require.NoError(t, err)
	require.NotEqual(t, password, hashed)
}

func TestCheckPasswordHash(t *testing.T) {
	hash := "$2a$04$09IzQ3oawFacAKHjG7QFneYFIaxV2fCNy7RG63RlFKQd.1ChHU6Xa"
	isValid := CheckPasswordHash("super-secret-password", hash)

	require.True(t, isValid)
}
