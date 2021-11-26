package cache

import (
	"sync"

	"github.com/google/uuid"
)

type ActiveUsersRepo struct {
	users map[uuid.UUID]bool
	mu    *sync.RWMutex
}

func NewActiveUsersRepo() *ActiveUsersRepo {
	return &ActiveUsersRepo{
		users: make(map[uuid.UUID]bool),
		mu:    &sync.RWMutex{},
	}
}

func (r *ActiveUsersRepo) Add(userId uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[userId] = true
	return nil
}

func (r *ActiveUsersRepo) Delete(userId uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.users, userId)
	return nil
}

func (r *ActiveUsersRepo) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.users)
}
