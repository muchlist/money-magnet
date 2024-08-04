package xulid

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type NullULID struct {
	ULID  ULID
	Valid bool
}

// Scan implements the Scanner interface.
func (n *NullULID) Scan(value interface{}) error {
	if value == nil {
		n.ULID, n.Valid = ULID{}, false
		return nil
	}

	switch v := value.(type) {
	case []byte:
		id, err := Parse(string(v))
		if err != nil {
			return err
		}
		n.ULID, n.Valid = id, true
	case string:
		id, err := Parse(v)
		if err != nil {
			return err
		}
		n.ULID, n.Valid = id, true
	default:
		return fmt.Errorf("invalid type for NullULID: %T", value)
	}

	return nil
}

// Value implements the driver Valuer interface.
func (n NullULID) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.ULID.String(), nil
}

// MarshalJSON implements the json.Marshaler interface.
func (n NullULID) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(n.ULID.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (n *NullULID) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Valid = false
		return nil
	}

	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	id, err := Parse(str)
	if err != nil {
		return err
	}

	n.ULID = id
	n.Valid = true
	return nil
}
