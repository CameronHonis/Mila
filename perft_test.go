package main_test

import (
	"bufio"
	"fmt"
	"github.com/CameronHonis/Mila"
	. "github.com/onsi/ginkgo/v2"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const QUIET = true
const PRINT_ROOT_MOVE_NODES = false
const FOCUS_TEST_IDX = -1
const MAX_DEPTH = 4

var scanner *bufio.Scanner

func perft(pos *main.Position, depth int) int {
	return _perft(pos, depth, true)
}

func _perft(pos *main.Position, depth int, isRoot bool) int {
	isLeaf := depth == 1
	nodeCnt := 0
	iter := main.NewLegalMoveIter(pos)
	var done bool
	for {
		var move main.Move
		move, done = iter.Next()
		if done {
			break
		}
		var moveNodeCnt int
		if isLeaf {
			moveNodeCnt = 1
		} else {
			capturedPiece, lastFrozenPos := pos.MakeMove(move)
			moveNodeCnt = _perft(pos, depth-1, false)
			pos.UnmakeMove(move, lastFrozenPos, capturedPiece)
		}
		if !QUIET && PRINT_ROOT_MOVE_NODES && isRoot {
			fmt.Printf("%s: %d\n", move, moveNodeCnt)
		}
		nodeCnt += moveNodeCnt
	}
	return nodeCnt
}

func parsePerftLine(line string) (fen string, depthNodeCntPairs [][2]int) {
	splitTxt := strings.Split(line, ";")
	fen = strings.TrimSpace(splitTxt[0])
	depthNodeCntPairs = make([][2]int, 0)
	for _, perftStr := range splitTxt[1:] {
		perftStrSplit := strings.Split(perftStr, " ")
		depthStr := perftStrSplit[0][1:]
		depth, parseDepthErr := strconv.Atoi(depthStr)
		if parseDepthErr != nil {
			log.Fatalf("could not parse depth from %s:\n\t%s", depthStr, parseDepthErr)
		}

		expNodeCntStr := perftStrSplit[1]
		expNodeCnt, parseNodeCntErr := strconv.Atoi(expNodeCntStr)
		if parseNodeCntErr != nil {
			log.Fatalf("could not parse expNodeCnt from %s:\n\t%s", expNodeCntStr, parseNodeCntErr)
		}
		depthNodeCntPairs = append(depthNodeCntPairs, [2]int{depth, expNodeCnt})
	}
	return fen, depthNodeCntPairs
}

func perftFromFile() {
	file, err := os.Open("./perft")
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}

	scanner = bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	testIdx := 0
	for scanner.Scan() {
		currTestIdx := testIdx
		testIdx++

		shouldSkipTest := FOCUS_TEST_IDX >= 0 && FOCUS_TEST_IDX != currTestIdx
		if shouldSkipTest {
			continue
		}
		line := scanner.Text()
		if !QUIET {
			fmt.Printf("[TEST %d] %s\n", currTestIdx, line)
		}
		fen, depthNodeCntPairs := parsePerftLine(line)

		pos, posErr := main.FromFEN(fen)
		if posErr != nil {
			log.Fatalf("could not build pos from FEN %s:\n\t%s", fen, posErr)
		}

		for _, depthNodeCntPair := range depthNodeCntPairs {
			depth := depthNodeCntPair[0]
			expNodeCnt := depthNodeCntPair[1]

			if depth > MAX_DEPTH {
				continue
			}

			start := time.Now()
			actNodeCnt := perft(pos, depth)
			if actNodeCnt != expNodeCnt {
				log.Fatalf("node count mismatch at depth %d, actual %d vs exp %d", depth, actNodeCnt, expNodeCnt)
			} else {
				elapsed := time.Since(start)
				if !QUIET {
					fmt.Printf("depth %d passed in %s\n", depth, elapsed)
				}
			}
		}
	}

	_ = file.Close()
}

var _ = It("perft", func() {
	//f, _ := os.Create("cpu.prof")
	//_ = pprof.StartCPUProfile(f)
	//defer pprof.StopCPUProfile()

	perftFromFile()
})
