package main

import (
	"fmt"
	"github.com/CameronHonis/chess"
	"log"
)

func main() {
	fmt.Println("Mila v0.0.0 - a lightweight chess AI written in go by Cameron Honis")
	tt := NewTranspTable()
	go NewUci(tt).Start()
	StartSearchTest(tt)
}

func StartSearchTest(tt *TranspTable) {
	pos, posErr := chess.BoardFromFEN(" 1r1q3r/pBP2pbp/1p2p1pn/4P2k/4QP2/B4N1P/P5P1/R4RK1 w - - 1 19")
	if posErr != nil {
		log.Fatalf("pos error %s", posErr)
	}
	NewSearch(pos, &SearchConstraints{maxDepth: 6}, tt).Start()
}
