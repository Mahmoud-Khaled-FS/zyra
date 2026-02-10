package utils

import (
	"fmt"
	"strings"
	"time"
)

func PrettyDuration(d time.Duration) string {
	if d < 0 {
		d = -d
	}

	h := d / time.Hour
	d -= h * time.Hour

	m := d / time.Minute
	d -= m * time.Minute

	s := d / time.Second
	d -= s * time.Second

	ms := d / time.Millisecond
	d -= ms * time.Millisecond

	us := d / time.Microsecond
	ns := d - us*time.Microsecond

	parts := []string{}
	if h > 0 {
		parts = append(parts, fmt.Sprintf("%dh", h))
	}
	if m > 0 {
		parts = append(parts, fmt.Sprintf("%dm", m))
	}
	if s > 0 {
		parts = append(parts, fmt.Sprintf("%ds", s))
	}
	if ms > 0 {
		parts = append(parts, fmt.Sprintf("%dms", ms))
	}
	if us > 0 {
		parts = append(parts, fmt.Sprintf("%dus", us))
	}
	if ns > 0 && h == 0 && m == 0 && s == 0 && ms == 0 && us == 0 {
		parts = append(parts, fmt.Sprintf("%dns", ns))
	}

	if len(parts) == 0 {
		return "0s"
	}

	return fmt.Sprintf("%s", join(parts, " "))
}

func join(elems []string, sep string) string {
	var res strings.Builder
	for i, s := range elems {
		res.WriteString(s)
		if i != len(elems)-1 {
			res.WriteString(sep)
		}
	}
	return res.String()
}
