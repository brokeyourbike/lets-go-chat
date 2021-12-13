package cache

import (
	"sync"

	"github.com/google/uuid"
)

type ActiveUsersRepo struct {
	users sync.Map
}

func NewActiveUsersRepo() *ActiveUsersRepo {
	return &ActiveUsersRepo{
		users: sync.Map{},
	}
}

func (r *ActiveUsersRepo) Add(userId uuid.UUID) error {
	r.users.Store(userId, true)
	return nil
}

func (r *ActiveUsersRepo) Delete(userId uuid.UUID) error {
	r.users.Delete(userId)
	return nil
}

func (r *ActiveUsersRepo) Count() int {
	count := 0
	r.users.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}
