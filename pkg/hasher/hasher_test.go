package hasher

import (
	"fmt"
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

func BenchmarkHashPassword(b *testing.B) {
	password := "super-secret-password"
	for n := 0; n < b.N; n++ {
		HashPassword(password)
	}
}

func BenchmarkCheckPasswordHash(b *testing.B) {
	password := "super-secret-password"
	hash := "$2a$04$09IzQ3oawFacAKHjG7QFneYFIaxV2fCNy7RG63RlFKQd.1ChHU6Xa"

	for n := 0; n < b.N; n++ {
		CheckPasswordHash(password, hash)
	}
}

func ExampleHashPassword() {
	hashed, _ := HashPassword("super")
	fmt.Println(hashed)
}

func ExampleCheckPasswordHash() {
	fmt.Println(CheckPasswordHash("super-secret-password", "$2a$04$09IzQ3oawFacAKHjG7QFneYFIaxV2fCNy7RG63RlFKQd.1ChHU6Xa"))
	fmt.Println(CheckPasswordHash("super-secret-password", "not-a-hash"))
}
