package xulid

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/oklog/ulid/v2"
)

// ULID is custom ulid.ULID who implement sql.Scanner, driver.Valuer, JSON and String
type ULID ulid.ULID

// Scan implements the Scanner interface.
func (u *ULID) Scan(value interface{}) error {
	if value == nil {
		*u = ULID{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		id, err := ulid.Parse(string(v))
		if err != nil {
			return err
		}
		*u = ULID(id)
	case string:
		id, err := ulid.Parse(v)
		if err != nil {
			return err
		}
		*u = ULID(id)
	default:
		return fmt.Errorf("invalid type for ULID: %T", value)
	}

	return nil
}

// Value implements the driver Valuer interface.
func (u ULID) Value() (driver.Value, error) {
	return ulid.ULID(u).String(), nil
}

func (u ULID) String() string {
	return ulid.ULID(u).String()
}

// MarshalJSON implements the json.Marshaler interface.
func (u ULID) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (u *ULID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	id, err := ulid.Parse(s)
	if err != nil {
		return err
	}
	*u = ULID(id)
	return nil
}
