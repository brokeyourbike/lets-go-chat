package cache

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddAndCount(t *testing.T) {
	cases := map[string]struct {
		setup     func(t *testing.T, repo *ActiveUsersRepo)
		wantCount int
	}{
		"it can Add item": {
			setup: func(t *testing.T, repo *ActiveUsersRepo) {
				key := uuid.New()
				require.NoError(t, repo.Add(key))
			},
			wantCount: 1,
		},
		"Add two items with the same key": {
			setup: func(t *testing.T, repo *ActiveUsersRepo) {
				key := uuid.New()
				require.NoError(t, repo.Add(key))
				require.NoError(t, repo.Add(key))
			},
			wantCount: 1,
		},
		"Add two items with different keys": {
			setup: func(t *testing.T, repo *ActiveUsersRepo) {
				key1 := uuid.New()
				key2 := uuid.New()

				require.NoError(t, repo.Add(key1))
				require.NoError(t, repo.Add(key2))
			},
			wantCount: 2,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repo := NewActiveUsersRepo()
			c.setup(t, repo)

			assert.Equal(t, c.wantCount, repo.Count())
		})
	}
}

func TestDeleteAndCount(t *testing.T) {
	cases := map[string]struct {
		setup     func(t *testing.T, repo *ActiveUsersRepo)
		wantCount int
	}{
		"Delete can remove item": {
			setup: func(t *testing.T, repo *ActiveUsersRepo) {
				key := uuid.New()

				require.NoError(t, repo.Add(key))
				repo.Delete(key)
			},
			wantCount: 0,
		},
		"Delete will remove only items with specified key": {
			setup: func(t *testing.T, repo *ActiveUsersRepo) {
				key1 := uuid.New()
				key2 := uuid.New()

				require.NoError(t, repo.Add(key1))
				require.NoError(t, repo.Add(key2))

				repo.Delete(key1)
			},
			wantCount: 1,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repo := NewActiveUsersRepo()
			c.setup(t, repo)

			assert.Equal(t, c.wantCount, repo.Count())
		})
	}
}

func TestCount(t *testing.T) {
	cases := map[string]struct {
		setup     func(t *testing.T, repo *ActiveUsersRepo)
		wantCount int
	}{
		"Count without items": {
			setup:     func(t *testing.T, repo *ActiveUsersRepo) {},
			wantCount: 0,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repo := NewActiveUsersRepo()
			c.setup(t, repo)

			assert.Equal(t, c.wantCount, repo.Count())
		})
	}
}
