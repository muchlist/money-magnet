package service

import "github.com/google/uuid"

type AddPersonData struct {
	Owner      uuid.UUID
	Person     uuid.UUID
	PocketID   uint64
	IsReadOnly bool
}

type RemovePersonData struct {
	Owner    uuid.UUID
	Person   uuid.UUID
	PocketID uint64
}
