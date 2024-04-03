package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func Bold(s string) string {
	return fmt.Sprintf("\033[1m%s\033[0m", s)
}

func Tabbed(times int, s string) string {
	return strings.Repeat("    ", times) + s
}

func PadToWidth(width int, s string) string {
	paddingNeeded := width - utf8.RuneCountInString(s)
	if paddingNeeded < 0 {
		return s
	}
	return s + strings.Repeat(" ", paddingNeeded)
}
