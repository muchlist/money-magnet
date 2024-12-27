package service

import (
	"github.com/muchlist/moneymagnet/business/notification/port"
	"github.com/muchlist/moneymagnet/pkg/mlogger"
)

// Core manages the set of APIs for user access.
type Core struct {
	log      mlogger.Logger
	userRepo port.UserStorer
}

// NewCore constructs a core for user api access.
func NewCore(
	log mlogger.Logger,
	userRepo port.UserStorer,
) *Core {
	return &Core{
		log:      log,
		userRepo: userRepo,
	}
}
