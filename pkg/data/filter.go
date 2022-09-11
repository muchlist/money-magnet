package data

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/muchlist/moneymagnet/pkg/utils/slicer"
)

// ========================================================== Filter Pagination

// Filters used for pagination and sort on database
type Filters struct {
	Page         int      `json:"page" validate:"gt=0,lte=1000000"`
	PageSize     int      `json:"page_size" validate:"gt=0,lte=100"`
	Sort         string   `json:"sort" validate:"required"`
	SortSafelist []string `json:"sort_safe_list"`
}

func (f *Filters) setDefault() {
	if f.Sort == "" {
		f.Sort = f.SortSafelist[0]
	}

	if f.Page == 0 {
		f.Page = 1
	}

	if f.PageSize == 0 {
		f.PageSize = 50
	}
}

// Validate do set default value for filter field and validate when
// user use the field
func (f *Filters) Validate() error {
	f.setDefault()

	if f.Page < 1 || f.Page > 1000 {
		return errors.New("invalid page value")
	}

	if f.PageSize < 1 || f.Page > 100 {
		return errors.New("invalid page_size value")
	}

	if len(f.SortSafelist) != 0 {
		if !slicer.In(f.Sort, f.SortSafelist) {
			return errors.New("invalid sort value")
		}
	}

	return nil
}

// Check that the client-provided Sort field matches one of the entries in our safelist
// and if it does, extract the column name from the Sort field by stripping the leading
// hyphen character (if one exists).
func (f *Filters) sortColumn() string {
	for _, safeValue := range f.SortSafelist {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}
	panic("unsafe sort parameter: " + f.Sort)
}

// return the sort direction ("ASC" or "DESC") depending on the prefix character of the
// Sort field.
func (f *Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

// SortColumnDirection return sort sql format
// ex: "id ASC"
func (f *Filters) SortColumnDirection() string {
	return fmt.Sprintf("%s %s", f.sortColumn(), f.sortDirection())
}

// Limit return limit (size per page)
func (f *Filters) Limit() int {
	return f.PageSize
}

// Offset return calculated offset
func (f *Filters) Offset() int {
	return (f.Page - 1) * f.PageSize
}

// ===========================================================Metadata Pagination

// Metadata for metadata pagination
type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty" example:"1"`
	PageSize     int `json:"page_size,omitempty" example:"50"`
	FirstPage    int `json:"first_page,omitempty" example:"1"`
	LastPage     int `json:"last_page,omitempty" example:"1"`
	TotalRecords int `json:"total_records,omitempty" example:"1"`
}

func CalculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}
