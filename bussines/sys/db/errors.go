package db

import (
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
)

var (
	ErrDBNotFound        = errors.New("not found")
	ErrDBDuplicatedEntry = errors.New("duplicated entry")
	ErrDBParentNotFound  = errors.New("invalid parent")
	ErrDBInvalidTextEnum = errors.New("invalid enum text input")
	ErrDBBuildQuery      = errors.New("query not valid")
	ErrDBSortFilter      = errors.New("invalid filter or sort value")
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
			return ErrDBParentNotFound
		case pgerrcode.InvalidTextRepresentation:
			return ErrDBInvalidTextEnum
		case pgerrcode.UndefinedColumn:
			return ErrDBBuildQuery
		}
	}

	return err
}
