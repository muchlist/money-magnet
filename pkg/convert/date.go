package convert

import (
	"errors"
	"strconv"
	"time"
)

func StringEpochToTime(str string) (time.Time, error) {
	if str == "" {
		return time.Time{}, errors.New("cannot be empty string")
	}

	number, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return time.Time{}, errors.New("must be epoch time format")
	}
	return time.Unix(number, 0), nil
}
