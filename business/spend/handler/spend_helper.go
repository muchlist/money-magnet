package handler

import (
	"net/url"

	"github.com/muchlist/moneymagnet/business/spend/model"
)

func extractSpendFilter(values url.Values) model.SpendFilter {
	rawFilter := model.SpendFilterRaw{
		User:      values.Get("user"),
		Category:  values.Get("category"),
		Name:      values.Get("name"),
		IsIncome:  values.Get("is_income"),
		Type:      values.Get("type"),
		DateStart: values.Get("date_start"),
		DateEnd:   values.Get("date_end"),
	}
	return rawFilter.ToModel()
}

func extractSpendMultiPocketFilter(values url.Values) model.SpendFilterMultiPocket {
	rawFilter := model.SpendFilterMultiPocketRaw{
		Pockets:    values.Get("pockets"),
		Users:      values.Get("users"),
		Categories: values.Get("categories"),
		Name:       values.Get("name"),
		IsIncome:   values.Get("is_income"),
		Types:      values.Get("types"),
		DateStart:  values.Get("date_start"),
		DateEnd:    values.Get("date_end"),
	}
	return rawFilter.ToModel()
}
