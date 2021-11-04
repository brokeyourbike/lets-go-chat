package hasher

import "testing"

const PASSWORD = "super-secret-password"

func TestHashPassword(t *testing.T) {
	if _, err := HashPassword(PASSWORD); err != nil {
		t.Errorf("fizzbuzz() = %q, want %v", err, nil)
	}
}

func TestCheckPasswordHash(t *testing.T) {
	want := "$2a$04$09IzQ3oawFacAKHjG7QFneYFIaxV2fCNy7RG63RlFKQd.1ChHU6Xa"
	if got := CheckPasswordHash(PASSWORD, want); got == false {
		t.Errorf("fizzbuzz() = %v, want %q", got, want)
	}
}
