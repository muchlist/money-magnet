package model

import (
	"strings"
	"time"

	"github.com/muchlist/moneymagnet/pkg/convert"
	"github.com/muchlist/moneymagnet/pkg/slicer"
	"github.com/muchlist/moneymagnet/pkg/xulid"
)

type SpendFilterMultiPocket struct {
	Pockets    []xulid.ULID
	Users      []xulid.ULID
	Categories []xulid.ULID
	Name       string
	IsIncome   *bool
	Types      []int
	DateStart  *time.Time
	DateEnd    *time.Time
}

type SpendFilterMultiPocketRaw struct {
	Pockets    string
	Users      string
	Categories string
	Name       string
	IsIncome   string
	Types      string
	DateStart  string
	DateEnd    string
}

func (p SpendFilterMultiPocketRaw) ToModel() SpendFilterMultiPocket {
	var result SpendFilterMultiPocket

	result.Pockets = xulid.ParseULIDs(p.Pockets)
	result.Users = xulid.ParseULIDs(p.Users)
	result.Categories = xulid.ParseULIDs(p.Categories)

	result.Name = strings.ToUpper(p.Name)
	result.IsIncome = convert.StringToPtrBool(p.IsIncome)
	result.Types, _ = slicer.CsvToSliceInt(p.Types)

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
