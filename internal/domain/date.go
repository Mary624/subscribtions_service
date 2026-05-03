package domain

import (
	"fmt"
	"strconv"
	"strings"
)

// Date model info
// @Description The date
type Date struct {
	Year  int
	Month int
}

func (d *Date) IsValid() bool {
	return d.Year > 1900 && d.Year < 3000 && d.Month > 0 && d.Month < 13
}

func (d *Date) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, d.String())), nil
}

func (d *Date) UnmarshalJSON(data []byte) error {
	dataStr := strings.ReplaceAll(string(data), "\"", "")
	if dataStr == "" {
		return nil
	}

	parts := strings.Split(dataStr, "-")
	if len(parts) != 2 {
		return fmt.Errorf("invalid date")
	}

	year, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("failed to parse year: %w", err)
	}

	if len(parts[0]) > 1 && parts[0][0] == 0 {
		parts[0] = parts[0][1:]
	}
	month, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("failed to parse month: %w", err)
	}

	d.Year = year
	d.Month = month

	return nil
}

func (d *Date) After(date Date) bool {
	if d.Year > date.Year {
		return true
	}
	if d.Year < date.Year {
		return false
	}

	return d.Month > date.Month
}

func (d *Date) String() string {
	if d.Year == 0 && d.Month == 0 {
		return ""
	}

	month := strconv.Itoa(d.Month)
	if d.Month < 10 {
		month = "0" + month
	}

	return fmt.Sprintf("%s-%d", month, d.Year)
}

func (d *Date) DateString() string {
	if d.Year == 0 && d.Month == 0 {
		return ""
	}

	month := strconv.Itoa(d.Month)
	if d.Month < 10 {
		month = "0" + month
	}

	return fmt.Sprintf("%d-%s-01", d.Year, month)
}
