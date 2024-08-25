package db

import (
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
)

var (
	ErrDBNotFound          = errors.New("not found")
	ErrDBDuplicatedEntry   = errors.New("duplicated entry")
	ErrDBRelationNotFound  = errors.New("invalid relation")
	ErrDBInvalidInput      = errors.New("invalid input syntax")
	ErrDBBuildQuery        = errors.New("query not valid")
	ErrDBSortFilter        = errors.New("invalid filter or sort value")
	ErrDBInvalidCursorType = errors.New("invalid cursor type")
)

func ParseError(err error) error {
	if err == pgx.ErrNoRows {
		return ErrDBNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return ErrDBDuplicatedEntry
		case pgerrcode.ForeignKeyViolation:
			return ErrDBRelationNotFound
		case pgerrcode.InvalidTextRepresentation:
			return ErrDBInvalidInput
		case pgerrcode.UndefinedColumn:
			return ErrDBBuildQuery
		}
	}

	return err
}
