package user

import (
	"testing"

	"github.com/matryer/is"
)

func TestItCanAddUser(t *testing.T) {
	is := is.New(t)

	users := make(Users)
	got := users.AddUser("key", User{})

	is.Equal("key", got)
}

func TestItGetUserByUserName(t *testing.T) {
	is := is.New(t)

	users := make(Users)
	user := User{UserName: "foo"}
	users["key"] = user

	got, err := users.GetUserByUserName("foo")
	is.NoErr(err)
	is.Equal(user, got)

	got, err = users.GetUserByUserName("bar")
	is.True(err != nil)
	is.Equal(User{}, got)
}
