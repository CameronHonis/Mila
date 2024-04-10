package main

import (
	"fmt"
	"github.com/CameronHonis/chess"
	"math"
	"time"
)

const ALPHA_BETA_PRUNING_ENABLED = true
const MOVE_SORT_ENABLED = true
const TRANSP_TABLE_LOOKUPS_ENABLED = true

type SearchResults struct {
	BestMove *chess.Move
	Score    float64
}

type Search struct {
	Root        *chess.Board
	TT          *TranspTable
	Constraints *SearchConstraints
	Depth       int

	isHalted bool
	nodeCnt  int
	line     []*chess.Move
	score    float64
}

func NewSearch(pos *chess.Board, constraints *SearchConstraints, tt *TranspTable) *Search {
	return &Search{
		Root:        pos,
		TT:          tt,
		Constraints: constraints,
		isHalted:    true,
		line:        make([]*chess.Move, 0),
	}
}

func (s *Search) UpdateFromResults(depth int, results *SearchResults) {
	s.Depth = depth
	s.line = []*chess.Move{results.BestMove}
	s.score = results.Score
}

func (s *Search) IncrNode() {
	s.nodeCnt++
	if s.nodeCnt >= s.Constraints.NodeCntLmt() {
		fmt.Println("halting search, max node count reached")
		s.isHalted = true
	}
}

func (s *Search) IncrDepth() {
	s.Depth++
	if s.Depth > s.Constraints.DepthLmt() {
		fmt.Println("halting search, max depth reached")
		s.isHalted = true
	}
}

func (s *Search) Halt() {
	s.isHalted = true
}

func (s *Search) IsHalted() bool {
	return s.isHalted
}

func (s *Search) NodeCnt() int {
	return s.nodeCnt
}

func (s *Search) Start() {
	s.isHalted = false
	maxSearchMs := s.MaxSearchMs()
	go func() {
		time.Sleep(time.Duration(maxSearchMs) * time.Millisecond)
		fmt.Println("halting search, search time allowance reached")
		s.Halt()
	}()
	var lastResultTime = time.Now()
	var line []*chess.Move
	for {
		s.IncrDepth()

		if s.IsHalted() {
			break
		}

		score := s.searchToDepth(s.Root, s.Depth, -math.MaxFloat64, math.MaxFloat64)
		dt := time.Now().Sub(lastResultTime)
		lastResultTime = time.Now()
		line = s.BestLine(s.Root)
		var lineStr string
		for moveIdx, move := range line {
			if moveIdx == 0 {
				lineStr += move.ToLongAlgebraic()
			} else {
				lineStr += " " + move.ToLongAlgebraic()
			}
		}
		fmt.Printf("info depth %d score %f move %s nodes %d time %dms\n", s.Depth, score, lineStr, s.NodeCnt(), dt.Milliseconds())
	}
	if len(line) > 0 {
		fmt.Printf("bestmove %s\n", line[0].ToLongAlgebraic())
	}

}

func (s *Search) searchToDepth(pos *chess.Board, depth int, alpha float64, beta float64) float64 {
	if depth == 0 || pos.Result != chess.BOARD_RESULT_IN_PROGRESS {
		s.IncrNode()
		return EvalPos(pos)
	}

	posHash := ZobristHash(pos)
	var anticipatedMove *chess.Move
	if TRANSP_TABLE_LOOKUPS_ENABLED {
		if ttEntry, _ := s.TT.GetEntry(posHash); ttEntry != nil {
			if ttEntry.Depth >= depth {
				return ttEntry.Score
			} else {
				anticipatedMove = ttEntry.Move
			}
		}
	}

	moves, err := chess.GetLegalMoves(pos)
	if err != nil || len(moves) == 0 {
		panic(fmt.Sprintf("could not get legal moves from pos %s: %s", pos, err))
	}
	if MOVE_SORT_ENABLED {
		moves = SortMoves(pos, moves, anticipatedMove)
	}
	var bestMove *chess.Move
	for _, move := range moves {
		newPos := chess.GetBoardFromMove(pos, move)
		score := -s.searchToDepth(newPos, depth-1, -beta, -alpha)
		if score > alpha {
			alpha = score
			bestMove = move
			if ALPHA_BETA_PRUNING_ENABLED && score > beta {
				break
			}
		}
	}
	results := &SearchResults{bestMove, alpha}
	s.TT.PostResults(posHash, results, depth)

	return alpha
}

func (s *Search) MaxSearchMs() int {
	var msForSearch = func(pos *chess.Board, bankMs int, incrMs int) int {
		expMoves := ExpMoves(pos)
		return incrMs + bankMs/expMoves
	}
	maxSearchMs := s.Constraints.MaxSearchMs()
	if s.Root.IsWhiteTurn && s.Constraints.whiteMs > 0 {
		maxSearchMs = MinInt(maxSearchMs, msForSearch(s.Root, s.Constraints.whiteMs, s.Constraints.whiteIncrMs))
	} else if !s.Root.IsWhiteTurn && s.Constraints.blackMs > 0 {
		maxSearchMs = MinInt(maxSearchMs, msForSearch(s.Root, s.Constraints.blackMs, s.Constraints.blackIncrMs))
	}
	return maxSearchMs
}

func (s *Search) BestLine(pos *chess.Board) []*chess.Move {
	rtn := make([]*chess.Move, 0)
	var ttEntry *TTEntry
	for {
		ttEntry, _ = s.TT.GetEntry(ZobristHash(pos))
		if ttEntry == nil {
			break
		}
		if ttEntry.Move == nil {
			break
		}
		rtn = append(rtn, ttEntry.Move)
		pos = chess.GetBoardFromMove(pos, ttEntry.Move)
	}
	return rtn
}
