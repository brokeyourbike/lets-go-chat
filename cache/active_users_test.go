package cache

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	repo := NewActiveUsersRepo()
	key := uuid.New()

	err := repo.Add(key)
	require.NoError(t, err)

	require.Equal(t, 1, repo.Count())
}

func TestDelete(t *testing.T) {
	repo := NewActiveUsersRepo()
	key := uuid.New()

	require.NoError(t, repo.Add(key))
	require.Equal(t, 1, repo.Count())

	repo.Delete(key)
	require.Equal(t, 0, repo.Count())
}

func TestCount(t *testing.T) {
	repo := NewActiveUsersRepo()
	key := uuid.New()

	require.NoError(t, repo.Add(key))
	require.NoError(t, repo.Add(key))

	require.Equal(t, 1, repo.Count())
}
