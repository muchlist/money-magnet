package daterange

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// DateRange holds the start and end time
type DateRange struct {
	StartDate time.Time
	EndDate   time.Time
}

// ParseDateRange parses the date range type and calculates the start and end dates.
// zone example : Asia/Makassar , Asia/Jakarta
// rangeType example : last-7-days. 2024-1, 2024-2
func ParseDateRange(rangeType string, zone string) (DateRange, error) {
	loc, err := time.LoadLocation(zone)
	if err != nil {
		return DateRange{}, fmt.Errorf("invalid timezone: %v", err)
	}

	now := time.Now().In(loc)

	switch {
	case rangeType == "last-7-days":
		year, month, day := now.Date()
		return calculateLast7Days(year, month, day, loc), nil
	case strings.Contains(rangeType, "-"): // Assume format is "yyyy-mm"
		year, month, err := parseYearMonth(rangeType)
		if err != nil {
			return DateRange{}, err
		}
		return calculateMonthRange(year, month, loc), nil
	default:
		return DateRange{}, fmt.Errorf("invalid range type")
	}
}

// parseYearMonth converts "yyyy-mm" string to year and month.
func parseYearMonth(yearMonthStr string) (int, time.Month, error) {
	parts := strings.Split(yearMonthStr, "-")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid year-month format")
	}
	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid year format: %v", err)
	}
	monthNum, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid month format: %v", err)
	}
	if monthNum < 1 || monthNum > 12 {
		return 0, 0, fmt.Errorf("month must be between 1 and 12")
	}
	return year, time.Month(monthNum), nil
}

// calculateMonthRange calculates the start and end dates for a given month.
func calculateMonthRange(year int, month time.Month, loc *time.Location) DateRange {
	start := time.Date(year, month, 1, 0, 0, 0, 0, loc)
	end := start.AddDate(0, 1, 0).Add(-time.Second)
	return DateRange{StartDate: start, EndDate: end}
}

// calculateLast7Days calculates the start and end dates for the last 7 days.
func calculateLast7Days(year int, month time.Month, day int, loc *time.Location) DateRange {
	start := time.Date(year, month, day-6, 0, 0, 0, 0, loc)
	end := time.Date(year, month, day, 23, 59, 59, 0, loc)
	return DateRange{StartDate: start, EndDate: end}
}
