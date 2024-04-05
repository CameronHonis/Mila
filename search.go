package main

import (
	"fmt"
	"github.com/CameronHonis/chess"
	"math"
	"time"
)

const ALPHA_BETA_PRUNING_ENABLED = true
const MOVE_SORT_ENABLED = true

type SearchResults struct {
	BestMove *chess.Move
	Score    float64
}

type Search struct {
	Root        *chess.Board
	Constraints *SearchConstraints
	Depth       int

	isHalted bool
	nodeCnt  int
	line     []*chess.Move
	score    float64
}

func NewSearch(pos *chess.Board, constraints *SearchConstraints) *Search {
	return &Search{
		Root:        pos,
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
	var bestMove *chess.Move
	for {
		s.IncrDepth()

		if s.IsHalted() {
			break
		}

		var score float64
		bestMove, score = s.searchToDepth(s.Root, s.Depth, -math.MaxFloat64, math.MaxFloat64)
		fmt.Printf("info depth %d score %f move %s nodes %d\n", s.Depth, score, bestMove.ToLongAlgebraic(), s.NodeCnt())
	}
	if bestMove != nil {
		fmt.Printf("bestmove %s\n", bestMove.ToLongAlgebraic())
	}

}

func (s *Search) searchToDepth(pos *chess.Board, depth int, alpha float64, beta float64) (bestMove *chess.Move, score float64) {
	if depth == 0 || pos.Result != chess.BOARD_RESULT_IN_PROGRESS {
		s.IncrNode()
		return nil, EvalPos(pos)
	}

	moves, err := chess.GetLegalMoves(pos)
	if err != nil || len(moves) == 0 {
		panic(fmt.Sprintf("could not get legal moves from pos %s: %s", pos, err))
	}
	if MOVE_SORT_ENABLED {
		moves = SortMoves(pos, moves)
	}
	var bestScore float64
	for _, move := range moves {
		if s.IsHalted() {
			break
		}
		newPos := chess.GetBoardFromMove(pos, move)
		_, subScore := s.searchToDepth(newPos, depth-1, alpha, beta)
		if pos.IsWhiteTurn {
			if subScore > alpha {
				alpha = subScore
				bestScore = subScore
				bestMove = move
				if ALPHA_BETA_PRUNING_ENABLED && subScore > beta {
					// black would not get in this line
					break
				}
			}
		} else { // black turn
			if subScore < beta {
				beta = subScore
				bestScore = subScore
				bestMove = move
				if ALPHA_BETA_PRUNING_ENABLED && subScore < alpha {
					// white would not get in this line
					break
				}
			}
		}
	}

	return bestMove, bestScore
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
