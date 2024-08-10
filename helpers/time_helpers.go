package helpers

import (
	"database/sql"
	"fmt"
	"time"
)

// ParseDatetime parses a datetime string into a time.Time object.
func ParseDatetime(datetime []byte) (time.Time, error) {
	parsedTime, err := time.Parse("2006-01-02 15:04:05", string(datetime))
	if err != nil {
		return time.Time{}, fmt.Errorf("Failed to parse datetime: %v", err)
	}
	return parsedTime, nil
}

// ParseNullableDatetime parses a datetime string into a sql.NullTime object, handling NULL values.
func ParseNullableDatetime(datetime []byte) (sql.NullTime, error) {
	if datetime == nil {
		return sql.NullTime{Valid: false}, nil
	}
	parsedTime, err := time.Parse("2006-01-02 15:04:05", string(datetime))
	if err != nil {
		return sql.NullTime{Valid: false}, fmt.Errorf("Failed to parse nullable datetime: %v", err)
	}
	return sql.NullTime{Time: parsedTime, Valid: true}, nil
}

// FormatTime formats a time.Time object into a string without the timezone.
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// FormatNullableTime formats a sql.NullTime object into a string, handling NULL values.
func FormatNullableTime(nt sql.NullTime) string {
	if nt.Valid {
		return nt.Time.Format("2006-01-02 15:04:05")
	}
	return ""
}
