package main_test

//
//import (
//	"bufio"
//	"fmt"
//	"github.com/CameronHonis/Mila"
//	"github.com/CameronHonis/chess"
//	. "github.com/onsi/ginkgo/v2"
//	"log"
//	"os"
//	"strconv"
//	"strings"
//	"time"
//)
//
//const PRINT_LEAFS = false
//const FOCUS_TEST_IDX = -1
//const MAX_DEPTH = 3
//
//var scanner *bufio.Scanner
//
//func perft(state *main.State, depth int) int {
//	return _perft(state, depth, make([]*chess.Move, 0))
//}
//
//func _perft(state *main.State, depth int, hist []*chess.Move) int {
//	leaf := depth == 1
//	nodeCnt := 0
//	iter := main.NewLegalMoveIter(state)
//	var done bool
//	for !done {
//		var move main.Move
//		move, done = iter.Next()
//		nextState :=
//	}
//	moves, movesErr := chess.GetLegalMoves(board)
//	if movesErr != nil {
//		log.Fatalf("error getting moves for board %s: %s", board, movesErr)
//	}
//	for _, move := range moves {
//		nextBoard := chess.GetBoardFromMove(board, move)
//		newHist := hist
//		newHist = append(newHist, move)
//		if leaf && PRINT_LEAFS {
//			printLeaf(nextBoard, newHist)
//		}
//		if nextBoard.IsCheckmate() {
//			if leaf {
//				nodeCnt++
//			}
//			continue
//		}
//		if leaf {
//			nodeCnt++
//		} else {
//			nodeCnt += _perft(nextBoard, depth-1, newHist)
//		}
//	}
//	return nodeCnt
//}
//
//func printLeaf(board *chess.Board, moves []*chess.Move) {
//	out := strings.Builder{}
//	out.WriteString("leaf ")
//	out.WriteString(board.ToFEN())
//	out.WriteString(" | ")
//	for _, move := range moves {
//		out.WriteString(fmt.Sprintf("%s%s ", move.StartSquare.ToAlgebraicCoords(), move.EndSquare.ToAlgebraicCoords()))
//	}
//	fmt.Println(out.String())
//}
//
//func parsePerftLine(line string) (fen string, depthNodeCntPairs [][2]int) {
//	splitTxt := strings.Split(line, ";")
//	fen = splitTxt[0]
//	depthNodeCntPairs = make([][2]int, 0)
//	for _, perftStr := range splitTxt[1:] {
//		perftStrSplit := strings.Split(perftStr, " ")
//		depthStr := perftStrSplit[0][1:]
//		depth, parseDepthErr := strconv.Atoi(depthStr)
//		if parseDepthErr != nil {
//			log.Fatalf("could not parse depth from %s:\n\t%s", depthStr, parseDepthErr)
//		}
//
//		expNodeCntStr := perftStrSplit[1]
//		expNodeCnt, parseNodeCntErr := strconv.Atoi(expNodeCntStr)
//		if parseNodeCntErr != nil {
//			log.Fatalf("could not parse expNodeCnt from %s:\n\t%s", expNodeCntStr, parseNodeCntErr)
//		}
//		depthNodeCntPairs = append(depthNodeCntPairs, [2]int{depth, expNodeCnt})
//	}
//	return fen, depthNodeCntPairs
//}
//
//func perftFromFile() {
//	file, err := os.Open("./perft")
//	if err != nil {
//		log.Fatalf("failed to open file: %s", err)
//	}
//
//	scanner = bufio.NewScanner(file)
//	scanner.Split(bufio.ScanLines)
//
//	testIdx := 0
//	for scanner.Scan() {
//		currTestIdx := testIdx
//		testIdx++
//
//		shouldSkipTest := FOCUS_TEST_IDX >= 0 && FOCUS_TEST_IDX != currTestIdx
//		if shouldSkipTest {
//			continue
//		}
//		line := scanner.Text()
//		fmt.Printf("[TEST %d] %s\n", currTestIdx, line)
//		fen, depthNodeCntPairs := parsePerftLine(line)
//
//		board, boardErr := chess.BoardFromFEN(fen)
//		if boardErr != nil {
//			log.Fatalf("could not construct board from FEN %s:\n\t%s", fen, boardErr)
//		}
//
//		for _, depthNodeCntPair := range depthNodeCntPairs {
//			depth := depthNodeCntPair[0]
//			expNodeCnt := depthNodeCntPair[1]
//
//			if depth > MAX_DEPTH {
//				continue
//			}
//
//			start := time.Now()
//			actNodeCnt := perft(board, depth)
//			if actNodeCnt != expNodeCnt {
//				log.Fatalf("node count mismatch at depth %d, actual %d vs exp %d", depth, actNodeCnt, expNodeCnt)
//			} else {
//				elapsed := time.Since(start)
//				fmt.Printf("depth %d passed in %s\n", depth, elapsed)
//			}
//		}
//	}
//
//	_ = file.Close()
//}
//
//var _ = It("perft", func() {
//	//f, _ := os.Create("cpu.prof")
//	//_ = pprof.StartCPUProfile(f)
//	//defer pprof.StopCPUProfile()
//	perftFromFile()
//	//board, _ := chess.BoardFromFEN("8/8/8/8/8/8/6k1/4K2R w K - 0 1")
//	//perft(board, 3)
//})
