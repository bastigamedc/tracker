package store

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func encode(v any) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

func decode(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

var durationRe = regexp.MustCompile(`^(?:(\d+)h)?\s*(?:(\d+)m)?$`)

func ParseDuration(s string) (int, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return 0, fmt.Errorf("empty duration")
	}
	m := durationRe.FindStringSubmatch(s)
	if m == nil {
		return 0, fmt.Errorf("invalid duration %q (expected e.g. 2h30m)", s)
	}
	var minutes int
	if m[1] != "" {
		h, _ := strconv.Atoi(m[1])
		minutes += h * 60
	}
	if m[2] != "" {
		mm, _ := strconv.Atoi(m[2])
		minutes += mm
	}
	if minutes == 0 {
		return 0, fmt.Errorf("duration is zero")
	}
	return minutes, nil
}

func FormatDuration(minutes int) string {
	h := minutes / 60
	m := minutes % 60
	if h == 0 {
		return fmt.Sprintf("%dm", m)
	}
	if m == 0 {
		return fmt.Sprintf("%dh", h)
	}
	return fmt.Sprintf("%dh%dm", h, m)
}

func MinutesFromDuration(s string) int {
	m, err := ParseDuration(s)
	if err != nil {
		return 0
	}
	return m
}
