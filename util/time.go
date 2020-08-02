package util

import (
	"fmt"
	"time"
)

func FormatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	s := m / time.Second
	if h > 0 {
		return fmt.Sprintf("%02dh%02dm", h, m)
	} else if m > 0 {
		return fmt.Sprintf("%02dm", m)
	} else {
		return fmt.Sprintf("%ds", s)
	}
}
