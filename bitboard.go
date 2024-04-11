package main

import (
	"fmt"
	"strings"
)

type Bitboard uint64

func (b Bitboard) String() string {
	var rtnBuilder = strings.Builder{}
	for rank := 8; rank > 0; rank-- {
		rtnBuilder.WriteString(fmt.Sprintf("%d ", rank))
		var mask Bitboard
		for file := 1; file < 9; file++ {
			idx := 8*(rank-1) + (file - 1)
			mask = 1 << idx
			if mask&b > 0 {
				rtnBuilder.WriteString("██")
			} else {
				var isDark bool
				if rank%2 == 0 {
					isDark = file%2 == 0
				} else {
					isDark = file%2 == 1
				}
				if isDark {
					rtnBuilder.WriteString("░░")
				} else {
					rtnBuilder.WriteString("  ")
				}
			}
		}
		rtnBuilder.WriteByte('\n')
	}
	rtnBuilder.WriteString("  ")
	for file := 1; file < 9; file++ {
		rtnBuilder.WriteByte(byte('0' + file))
		rtnBuilder.WriteByte(' ')
	}
	return rtnBuilder.String()
}
