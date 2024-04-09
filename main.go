package main

import (
	"github.com/CameronHonis/chess"
	"os"
	"runtime/pprof"
)

//func main() {
//	fmt.Println("Mila v0.0.0 - a lightweight chess AI written in go by Cameron Honis")
//	NewUci(NewTranspTable()).Start()
//}

func main() {
	// QUESTIONS:
	// 1. What is taking so long to search 50k nodes in 1.7s?
	// creating contexts seems to be very expensive

	// 2: Why is the best move different in the context-closing version? (g2g4 vs. b7a6)
	// 2a. Why is the depth 4 line only 3 moves?

	file, _ := os.Create("cpu_prof_v0.3.1")
	_ = pprof.StartCPUProfile(file)
	tt := NewTranspTable()
	pos, _ := chess.BoardFromFEN("1r1q3r/pBP2pbp/1p2p1pn/4P2k/4QP2/B4N1P/P5P1/R4RK1 w - - 1 19")
	searchConstraints := &SearchConstraints{maxDepth: 4}
	NewSearch(pos, searchConstraints, tt).Start()
	pprof.StopCPUProfile()
}
