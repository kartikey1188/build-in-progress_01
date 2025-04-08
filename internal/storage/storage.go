package storage

import "github.com/kartikey1188/build-in-progress_01/internal/types"

type Storage interface {
	LoginAndRegister
}

type LoginAndRegister interface {
	CreateUser(user types.User) (int64, error)
	GetUserByEmail(email string) (types.User, error)
	UpdateLastLogin(userID int64, lastLogin types.DateTime) error
}
