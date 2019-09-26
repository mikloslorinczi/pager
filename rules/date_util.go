package rules

import (
	"fmt"
	"time"
)

// IsWorkHour checks if a given timestamp is in a workday's working hour
func IsWorkHour(stamp time.Time, location *time.Location) bool {
	// Check if stamp is a weekday or weekend
	isWorkDay := stamp.In(location).Weekday() > 0 && stamp.In(location).Weekday() < 6
	// Check if stamp is on an exception day
	for _, exception := range dayExceptions {
		if exception.year == stamp.In(location).Year() &&
			exception.month == int(stamp.In(location).Month()) &&
			exception.day == stamp.In(location).Day() {
			isWorkDay = exception.isWorkDay
		}
	}
	// Return true if this is a working day and is a working hour
	return isWorkDay && workHours[stamp.In(location).Hour()]
}

// BeginningOfMonth returns the begining time of a month
func BeginningOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth returns the end time of a month
func EndOfMonth(t time.Time) time.Time {
	return BeginningOfMonth(t).AddDate(0, 1, 0).Add(-time.Second)
}

// FormatDate formats a timestamp in the simple 1990.10.13 format
func FormatDate(t time.Time) string {
	return fmt.Sprintf("%d.%d.%d", t.Year(), t.Month(), t.Day())
}
