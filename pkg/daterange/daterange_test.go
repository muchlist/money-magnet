// daterange_test.go
package daterange

import (
	"testing"
	"time"
)

func TestParseDateRange(t *testing.T) {
	// Test for "last-7-days"
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	expectedStart := time.Date(now.Year(), now.Month(), now.Day()-6, 0, 0, 0, 0, loc)
	expectedEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, loc)

	dr, err := ParseDateRange("last-7-days", "Asia/Jakarta")
	if err != nil {
		t.Errorf("ParseDateRange returned an error: %v", err)
	}
	if !dr.StartDate.Equal(expectedStart) || !dr.EndDate.Equal(expectedEnd) {
		t.Errorf("ParseDateRange returned incorrect dates: got %v to %v, want %v to %v",
			dr.StartDate, dr.EndDate, expectedStart, expectedEnd)
	}

	// Test for invalid range type
	_, err = ParseDateRange("invalid-type", "Asia/Jakarta")
	if err == nil {
		t.Error("ParseDateRange should return an error for invalid range types")
	}
}

func TestParseMonth(t *testing.T) {
	year, monthNum, err := parseYearMonth("2024-2")
	if err != nil {
		t.Errorf("parseYearMonth returned an error: %v", err)
	}
	if monthNum != 2 {
		t.Errorf("parseYearMonth returned %d, want %d", monthNum, 2)
	}
	if year != 2024 {
		t.Errorf("parseYearMonth returned %d, want %d", year, 2024)
	}

	_, _, err = parseYearMonth("invalid")
	if err == nil {
		t.Error("parseMonth should return an error for invalid input")
	}
}

func TestCalculateMonthRange(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	year := 2023
	month := time.February

	expectedStart := time.Date(year, month, 1, 0, 0, 0, 0, loc)
	expectedEnd := expectedStart.AddDate(0, 1, 0).Add(-time.Second)

	dr := calculateMonthRange(year, month, loc)
	if !dr.StartDate.Equal(expectedStart) || !dr.EndDate.Equal(expectedEnd) {
		t.Errorf("calculateMonthRange returned incorrect dates: got %v to %v, want %v to %v",
			dr.StartDate, dr.EndDate, expectedStart, expectedEnd)
	}
}
