package paging

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/muchlist/moneymagnet/pkg/slicer"
)

// Cursor digunakan untuk pagination dan sorting pada database
type Cursor struct {
	cursor     string
	pageSize   int
	cursorType string
	cursorList []string
}

// Getter dan Setter untuk Cursor

func (f *Cursor) GetCursor() string {
	return f.cursor
}
func (f *Cursor) SetCursor(cursor string) {
	f.cursor = cursor
}
func (f *Cursor) GetPageSize() int {
	return f.pageSize
}
func (f *Cursor) GetPageSizePlusOne() int {
	return f.pageSize + 1 // need additional data for next cursor value
}
func (f *Cursor) SetPageSize(pageSize int) {
	f.pageSize = pageSize
}
func (f *Cursor) GetCursorType() string {
	return f.cursorType
}
func (f *Cursor) SetCursorType(cursorType string) {
	f.cursorType = cursorType
}
func (f *Cursor) GetCursorList() []string {
	return f.cursorList
}
func (f *Cursor) SetCursorList(cursorList []string) {
	f.cursorList = cursorList
}

// Validate melakukan validasi dan menetapkan nilai default untuk filter
func (f *Cursor) Validate() error {
	if len(f.cursorList) == 0 {
		return errors.New("list of fields that are safe to be cursor is required")
	}

	f.ensureDefault()

	if f.pageSize < 1 || f.pageSize > 101 {
		return errors.New("invalid page_size value")
	}

	if !slicer.In(f.cursorType, f.cursorList) {
		return errors.New("invalid cursor type")
	}

	return nil
}

// GetCursorColumn mengembalikan nama kolom yang aman berdasarkan cursor type
func (f *Cursor) GetCursorColumn() (string, error) {
	for _, safeValue := range f.cursorList {
		if f.cursorType == safeValue {
			return strings.TrimPrefix(f.cursorType, "-"), nil
		}
	}
	return "", errors.New("unsafe cursor parameter: " + f.cursorType)
}

// GetSortColumnDirection mengembalikan format sorting SQL seperti "id ASC"
func (f *Cursor) GetSortColumnDirection() (string, error) {
	column, err := f.GetCursorColumn()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s %s", column, f.getSortDirection()), nil
}

// Direction mengembalikan arah untuk where cursor
func (f *Cursor) GetDirection() string {
	if strings.HasPrefix(f.cursorType, "-") {
		return "<" // if sort by desc, cursor next direction is <
	}
	return ">"
}

// EnsureDefault() menetapkan nilai default untuk filter
func (f *Cursor) ensureDefault() {
	if f.cursorType == "" {
		f.cursorType = f.cursorList[0] // cursor type pertama adalah default
	}

	if f.pageSize == 0 {
		f.pageSize = 50 // default page size adalah 50
	}
}

// GetSortDirection mengembalikan arah sorting ("ASC" atau "DESC") berdasarkan prefix dari cursorType
func (f *Cursor) getSortDirection() string {
	if strings.HasPrefix(f.cursorType, "-") {
		return "DESC"
	}
	return "ASC"
}

// CursorMetadata menyimpan metadata untuk pagination
type CursorMetadata struct {
	CurrentCursor string `json:"current_cursor" example:"01ARZ3NDEKTSV4RRFFQ69G5FAW"`
	CursorType    string `json:"cursor_type" example:"id"`
	PageSize      int    `json:"page_size" example:"50"`
	NextCursor    string `json:"next_cursor" example:"01ARZ3NDEKTSV5RRFFQ69G5AAA"`
	NextPage      string `json:"next_page" example:"/users?limit=50&cursor=01ARZ3NDEKTSV5RRFFQ69G5AAA&cursor_type=id"`
	ReverseCursor string `json:"reverse_cursor" example:"01ARZ3NDEKTSV5RRFFQ69G5AAA"`
	ReversePage   string `json:"reverse_page" example:"/users?limit=50&cursor=01ARZ3NDEKTSV5RRFFQ69G5AAA&cursor_type=-id"`
}

func (c *CursorMetadata) GenerateAndApplyPageUri(basePath string, queryParams url.Values) {
	if c.NextCursor != "" {
		copyParamsForNextCursor := cloneQueryParams(queryParams)
		copyParamsForNextCursor.Set("cursor", c.NextCursor)
		copyParamsForNextCursor.Set("cursor_type", c.CursorType)
		copyParamsForNextCursor.Set("page_size", strconv.Itoa(c.PageSize))

		c.NextPage = basePath + "?" + copyParamsForNextCursor.Encode()
	}

	if c.ReverseCursor != "" {
		copyParamsForNextCursor := cloneQueryParams(queryParams)
		copyParamsForNextCursor.Set("cursor", c.ReverseCursor)
		reverseCursorType := ""
		if strings.Contains(c.CursorType, "-") {
			reverseCursorType = strings.TrimPrefix(c.CursorType, "-")
		} else {
			reverseCursorType = fmt.Sprintf("-%s", c.CursorType)
		}
		copyParamsForNextCursor.Set("cursor_type", reverseCursorType)
		copyParamsForNextCursor.Set("page_size", strconv.Itoa(c.PageSize))

		c.ReversePage = basePath + "?" + copyParamsForNextCursor.Encode()
	}
}

// cloneQueryParams makes a copy of url.Values to avoid modifying the original map.
func cloneQueryParams(params url.Values) url.Values {
	copyParams := make(url.Values)
	for key, values := range params {
		copyParams[key] = append([]string(nil), values...)
	}
	return copyParams
}
