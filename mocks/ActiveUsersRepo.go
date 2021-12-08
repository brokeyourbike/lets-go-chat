// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	uuid "github.com/google/uuid"
	mock "github.com/stretchr/testify/mock"
)

// ActiveUsersRepo is an autogenerated mock type for the ActiveUsersRepo type
type ActiveUsersRepo struct {
	mock.Mock
}

// Add provides a mock function with given fields: userId
func (_m *ActiveUsersRepo) Add(userId uuid.UUID) error {
	ret := _m.Called(userId)

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) error); ok {
		r0 = rf(userId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Count provides a mock function with given fields:
func (_m *ActiveUsersRepo) Count() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// Delete provides a mock function with given fields: userId
func (_m *ActiveUsersRepo) Delete(userId uuid.UUID) error {
	ret := _m.Called(userId)

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) error); ok {
		r0 = rf(userId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}