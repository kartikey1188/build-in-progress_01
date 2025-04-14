package storage

import "github.com/kartikey1188/build-in-progress_01/internal/types"

type Storage interface {
	LoginAndRegister
}

type LoginAndRegister interface {
	CreateCollectorUser(user types.Collector) (int64, error)
	CreateBusinessUser(user types.Business) (int64, error)
	GetUserByEmail(email string) (types.User, error)
	UpdateLastLogin(userID int64, lastLogin types.DateTime) error
	GetCollectorByEmail(email string) (types.Collector, error)
	GetBusinessByEmail(email string) (types.Business, error)
}
