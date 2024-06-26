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

func ConcatBitboards(bbs ...Bitboard) string {
	bbStrs := make([]string, len(bbs))
	for i, bb := range bbs {
		bbStrs[i] = bb.String()
	}
	return ConcatStringsHorizontally(bbStrs...)
}

func ConcatStringsHorizontally(multiLineStrs ...string) string {
	builder := strings.Builder{}
	splitLineStrs := make([][]string, 0)
	var maxRows int
	for _, multiLineStr := range multiLineStrs {
		splitLineStr := strings.Split(multiLineStr, "\n")
		splitLineStrs = append(splitLineStrs, splitLineStr)
		if len(splitLineStr) > maxRows {
			maxRows = len(splitLineStr)
		}
	}
	for row := 0; row < maxRows; row++ {
		for _, splitLineStr := range splitLineStrs {
			if row < len(splitLineStr) {
				builder.WriteString(splitLineStr[row])
			}
		}
		builder.WriteByte('\n')
	}
	return builder.String()
}

func ReverseSlice[T any](s []T) {
	for i := len(s)/2 - 1; i >= 0; i-- {
		opp := len(s) - 1 - i
		s[i], s[opp] = s[opp], s[i]
	}
}
