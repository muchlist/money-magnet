package xulid

import (
	"strings"

	"github.com/oklog/ulid/v2"
)

func Parse(ulidString string) (ULID, error) {
	u, err := ulid.Parse(ulidString)
	if err != nil {
		return ULID{}, err
	}
	return ConvertULIDToXULID(u), nil
}

func ParseULIDs(input string) []ULID {
	var ulids []ULID
	if input == "" {
		return ulids
	}

	elements := strings.Split(input, ",")
	for _, elem := range elements {
		parsedULID, err := Parse(elem)
		if err == nil {
			ulids = append(ulids, parsedULID)
		}
	}
	return ulids
}

func MustParse(ulidString string) ULID {
	u := ulid.MustParse(ulidString)
	return ConvertULIDToXULID(u)
}

// ConvertULIDToXULID converts xulid.ULID to XULID
func ConvertULIDToXULID(id ulid.ULID) ULID {
	return ULID(id)
}
