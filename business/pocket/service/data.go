package service

import "github.com/google/uuid"

type AddPersonData struct {
	Owner      uuid.UUID
	Person     uuid.UUID
	PocketID   uuid.UUID
	IsReadOnly bool
}

type RemovePersonData struct {
	Owner    uuid.UUID
	Person   uuid.UUID
	PocketID uuid.UUID
}
