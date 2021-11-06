package main

import (
	"fmt"

	"github.com/brokeyourbike/lets-go-chat/pkg/hasher"
)

func main() {
	hash, _ := hasher.HashPassword("super-secret-password")
	fmt.Println(hash)
}
