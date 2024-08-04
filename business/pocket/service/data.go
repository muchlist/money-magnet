package service

import "github.com/muchlist/moneymagnet/pkg/xulid"

type AddPersonData struct {
	Owner      xulid.ULID
	Person     xulid.ULID
	PocketID   xulid.ULID
	IsReadOnly bool
}

type RemovePersonData struct {
	Owner    xulid.ULID
	Person   xulid.ULID
	PocketID xulid.ULID
}
