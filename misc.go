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

func Yellow(s string) string {
	return fmt.Sprintf("\033[33m%s\033[0m", s)
}

func MaxInt(int1, int2 int) int {
	if int1 > int2 {
		return int1
	}
	return int2
}

func MinInt(int1, int2 int) int {
	if int1 < int2 {
		return int1
	}
	return int2
}

func Swap[T interface{}](slice []T, i, j int) {
	temp := slice[i]
	slice[i] = slice[j]
	slice[j] = temp
}
