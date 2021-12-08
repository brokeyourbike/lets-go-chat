// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	models "github.com/brokeyourbike/lets-go-chat/models"
	mock "github.com/stretchr/testify/mock"
)

// UsersRepo is an autogenerated mock type for the UsersRepo type
type UsersRepo struct {
	mock.Mock
}

// Create provides a mock function with given fields: user
func (_m *UsersRepo) Create(user models.User) error {
	ret := _m.Called(user)

	var r0 error
	if rf, ok := ret.Get(0).(func(models.User) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByUserName provides a mock function with given fields: userName
func (_m *UsersRepo) GetByUserName(userName string) (models.User, error) {
	ret := _m.Called(userName)

	var r0 models.User
	if rf, ok := ret.Get(0).(func(string) models.User); ok {
		r0 = rf(userName)
	} else {
		r0 = ret.Get(0).(models.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}