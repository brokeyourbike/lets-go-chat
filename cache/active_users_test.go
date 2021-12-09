package cache

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ActiveUsersSuite struct {
	suite.Suite
	repository *ActiveUsersRepo
}

func (s *ActiveUsersSuite) SetupTest() {
	s.repository = NewActiveUsersRepo()
}

func TestActiveUsers(t *testing.T) {
	suite.Run(t, new(ActiveUsersSuite))
}

func (s *ActiveUsersSuite) Test_repository_Add() {
	key := uuid.New()

	err := s.repository.Add(key)
	require.NoError(s.T(), err)

	require.Equal(s.T(), 1, s.repository.Count())
}

func (s *ActiveUsersSuite) Test_repository_Delete() {
	key := uuid.New()

	require.NoError(s.T(), s.repository.Add(key))
	require.Equal(s.T(), 1, s.repository.Count())

	s.repository.Delete(key)
	require.Equal(s.T(), 0, s.repository.Count())
}

func (s *ActiveUsersSuite) Test_repository_Count() {
	key := uuid.New()

	require.NoError(s.T(), s.repository.Add(key))
	require.NoError(s.T(), s.repository.Add(key))

	require.Equal(s.T(), 1, s.repository.Count())
}
