package ctype

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

// UppercaseString is a custom type that automatically converts strings to uppercase.
type UppercaseString string

func ToUppercaseString(s string) UppercaseString {
	return UppercaseString(strings.ToUpper(s))
}

func FromUppercaseString(s UppercaseString) string {
	return string(s)
}

func (u UppercaseString) Value() (driver.Value, error) {
	return strings.ToUpper(string(u)), nil
}

func (u *UppercaseString) Scan(value interface{}) error {
	if value == nil {
		*u = UppercaseString("")
		return nil
	}

	switch v := value.(type) {
	case string:
		*u = UppercaseString(strings.ToUpper(v))
	case []byte:
		*u = UppercaseString(strings.ToUpper(string(v)))
	default:
		return fmt.Errorf("cannot scan type %T into UppercaseString", value)
	}
	return nil
}
