package model

import (
	"time"

	"github.com/muchlist/moneymagnet/pkg/convert"
	"github.com/muchlist/moneymagnet/pkg/slicer"
	"github.com/muchlist/moneymagnet/pkg/xulid"
)

type SpendFilter struct {
	PocketID  xulid.NullULID
	User      xulid.NullULID
	Category  xulid.NullULID
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
		userULID, err := xulid.Parse(p.User)
		if err == nil {
			result.User.ULID = userULID
			result.User.Valid = true
		}
	}

	// category must be uuid format
	if p.Category != "" {
		categoryULID, err := xulid.Parse(p.Category)
		if err == nil {
			result.Category.ULID = categoryULID
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
