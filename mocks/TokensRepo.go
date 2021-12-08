// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	models "github.com/brokeyourbike/lets-go-chat/models"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// TokensRepo is an autogenerated mock type for the TokensRepo type
type TokensRepo struct {
	mock.Mock
}

// Create provides a mock function with given fields: token
func (_m *TokensRepo) Create(token models.Token) error {
	ret := _m.Called(token)

	var r0 error
	if rf, ok := ret.Get(0).(func(models.Token) error); ok {
		r0 = rf(token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: id
func (_m *TokensRepo) Get(id uuid.UUID) (models.Token, error) {
	ret := _m.Called(id)

	var r0 models.Token
	if rf, ok := ret.Get(0).(func(uuid.UUID) models.Token); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(models.Token)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InvalidateByUserId provides a mock function with given fields: userId
func (_m *TokensRepo) InvalidateByUserId(userId uuid.UUID) error {
	ret := _m.Called(userId)

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) error); ok {
		r0 = rf(userId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}