package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/muchlist/moneymagnet/pkg/utils/convert"
	"github.com/muchlist/moneymagnet/pkg/utils/slicer"
)

type SpendFilter struct {
	PocketID  uuid.NullUUID
	User      uuid.NullUUID
	Category  uuid.NullUUID
	IsIncome  *bool
	Type      []int
	DateStart *time.Time
	DateEnd   *time.Time
}

type SpendFilterRaw struct {
	User      string
	Category  string
	IsIncome  string
	Type      string
	DateStart string
	DateEnd   string
}

func (p SpendFilterRaw) ToModel() SpendFilter {
	var result SpendFilter

	// user must be uuid format
	if p.User != "" {
		userUUID, err := uuid.Parse(p.User)
		if err == nil {
			result.User.UUID = userUUID
			result.User.Valid = true
		}
	}

	// category must be uuid format
	if p.Category != "" {
		categoryUUID, err := uuid.Parse(p.Category)
		if err == nil {
			result.Category.UUID = categoryUUID
			result.Category.Valid = true
		}
	}

	result.IsIncome = convert.StringToPtrBool(p.IsIncome)

	result.Type, _ = slicer.CsvToSliceInt(p.Type)

	start, err := convert.StringEpochToTime(p.DateStart)
	if err == nil {
		result.DateStart = &start
	}

	end, err := convert.StringEpochToTime(p.DateEnd)
	if err == nil {
		result.DateEnd = &end
	}

	return result
}
