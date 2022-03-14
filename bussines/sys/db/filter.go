package db

import (
	"fmt"
	"strings"

	"github.com/muchlist/moneymagnet/bussines/sys/validate"
	"github.com/muchlist/moneymagnet/foundation/tools"
)

type Filters struct {
	Page         int      `json:"page" validate:"gt=0,lte=1000000"`
	PageSize     int      `json:"page_size" validate:"gt=0,lte=100"`
	Sort         string   `json:"sort" validate:"required"`
	SortSafelist []string `json:"sort_safe_list"`
}

func (f *Filters) setDefault() {
	if f.Page == 0 {
		f.Page = 1
	}

	if f.PageSize == 0 {
		f.PageSize = 50
	}
}

func (f *Filters) Validate() error {

	f.setDefault()

	msg, err := validate.Struct(f)
	if err != nil {
		return fmt.Errorf(msg)
	}

	if len(f.SortSafelist) != 0 {
		if !tools.In(f.Sort, f.SortSafelist...) {
			return fmt.Errorf("invalid sort value")
		}
	}

	return nil
}

func (f *Filters) SortColumn() string {
	for _, safeValue := range f.SortSafelist {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}
	panic("unsafe sort parameter: " + f.Sort)
}

func (f *Filters) SortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func (f *Filters) SortColumnDirection() string {
	return fmt.Sprintf("%s %s", f.SortColumn(), f.SortDirection())
}

func (f *Filters) Limit() int {
	return f.PageSize
}

func (f *Filters) Offset() int {
	return (f.Page - 1) * f.PageSize
}
