package main

import (
	"github.com/CameronHonis/chess"
)

//func main() {
//	fmt.Println("Mila v0.0.0 - a lightweight chess AI written in go by Cameron Honis")
//	NewUci(NewTranspTable()).Start()
//}

func main() {
	tt := NewTranspTable()
	pos, _ := chess.BoardFromFEN("1r1q3r/pBP2pbp/1p2p1pn/4P2k/4QP2/B4N1P/P5P1/R4RK1 w - - 1 19")
	searchConstraints := &SearchConstraints{maxDepth: 6}
	NewSearch(pos, searchConstraints, tt).Start()
}
