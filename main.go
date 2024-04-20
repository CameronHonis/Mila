package main

import (
	"fmt"
	"os"
	"runtime/pprof"
)

const DEBUG = false
const PROFILE = false

func main() {
	initAttackPrecomputes()

	if PROFILE {
		f, _ := os.Create("cpu.prof")
		_ = pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	fmt.Println("Mila v0.0.0 - a lightweight chess AI written in go by Cameron Honis")
	tt := NewTranspTable()
	go NewUci(tt).Start()
}

func StartSearchTest(tt *TranspTable) {
	pos, _ := FromFEN("1r1q3r/pBP2pbp/1p2p1pn/4P2k/4QP2/B4N1P/P5P1/R4RK1 w - - 1 19")
	//pos, _ := FromFEN("5n1k/5p1p/4p1pP/6P1/1p6/pP2R2B/P7/K7 w - - 0 1")
	//pos, _ := FromFEN("7k/8/8/5p1p/5rQ1/8/3K2R1/8 w - - 0 1")
	NewSearch(pos, &SearchConstraints{maxDepth: 10}, tt).Start()
}
