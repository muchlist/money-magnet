package db

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
)

var (
	ErrDBNotFound        = errors.New("not found")
	ErrDBDuplicatedEntry = errors.New("duplicated entry")
	ErrDBParentNotFound  = errors.New("invalid parent")
	ErrDBInvalidTextEnum = errors.New("invalid enum text input")
	ErrDBBuildQuery      = errors.New("query not valid")
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

// Returning return for helping build query
// input : id, name, category
// output : "RETURNING id, name, category"
func Returning(columns ...string) string {
	sb := strings.Builder{}
	sb.WriteString("RETURNING ")
	for _, key := range columns {
		sb.WriteString(key + ", ")
	}
	return strings.TrimSuffix(sb.String(), ", ")
}

// A return A.text for helping join query
// input : updated_at
// output : "A.updated_at"
func A(text string) string {
	return fmt.Sprintf("A.%s", text)
}

// B return B.text for helping join query
// input : updated_at
// output : "B.updated_at"
func B(text string) string {
	return fmt.Sprintf("B.%s", text)
}

// C return C.text for helping join query
// input : updated_at
// output : "C.updated_at"
func C(text string) string {
	return fmt.Sprintf("C.%s", text)
}

// Dot return table.column for helping join query
// input : (user , updated_at)
// output : "user.updated_at"
func Dot(table, column string) string {
	return fmt.Sprintf("%s.%s", table, column)
}

// CoalesceInt Coalesce(null,default)
func CoalesceInt(text string, def int) string {
	return fmt.Sprintf("Coalesce(%s,%d)", text, def)
}
