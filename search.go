package main

import (
	"fmt"
	"github.com/CameronHonis/chess"
	"math"
	"sync"
	"time"
)

const ALPHA_BETA_PRUNING_ENABLED = true
const MOVE_SORT_ENABLED = true
const TRANSP_TABLE_LOOKUPS_ENABLED = true
const MULTITHREADING_ENABLED = true

type SearchResult struct {
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
	mu       sync.Mutex
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

func (s *Search) IncrNode() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nodeCnt++
	if s.nodeCnt >= s.Constraints.NodeCntLmt() {
		fmt.Println("halting search, max node count reached")
		s.isHalted = true
	}
}

func (s *Search) IncrDepth() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Depth++
	if s.Depth > s.Constraints.DepthLmt() {
		fmt.Println("halting search, max depth reached")
		s.isHalted = true
	}
}

func (s *Search) Halt() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.isHalted = true
}

func (s *Search) IsHalted() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.isHalted
}

func (s *Search) NodeCnt() int {
	s.mu.Lock()
	defer s.mu.Unlock()
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
	var lastResultTime = time.Now()
	for {
		s.IncrDepth()

		if s.IsHalted() {
			break
		}

		results := s.searchToDepth(s.Root, s.Depth, -math.MaxFloat64, math.MaxFloat64)
		bestMove = results.BestMove
		dt := time.Now().Sub(lastResultTime)
		lastResultTime = time.Now()
		fmt.Printf("info depth %d score %f move %s nodes %d time %dms\n", s.Depth, results.Score, bestMove.ToLongAlgebraic(), s.NodeCnt(), dt.Milliseconds())
	}
	if bestMove != nil {
		fmt.Printf("bestmove %s\n", bestMove.ToLongAlgebraic())
	}

}

func (s *Search) searchToDepth(pos *chess.Board, depth int, alpha float64, beta float64) *SearchResult {
	if depth == 0 || pos.Result != chess.BOARD_RESULT_IN_PROGRESS {
		s.IncrNode()
		res := &SearchResult{BestMove: nil, Score: EvalPos(pos)}
		return res
	}

	posHash := ZobristHash(pos)
	var anticipatedMove *chess.Move
	if TRANSP_TABLE_LOOKUPS_ENABLED {
		if ttEntry, _ := s.TT.GetEntry(posHash); ttEntry != nil {
			if ttEntry.Depth >= depth {
				return &SearchResult{
					BestMove: ttEntry.Move,
					Score:    ttEntry.Score,
				}
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

	var rtn = &SearchResult{}
	var handleResult = func(res *SearchResult, move *chess.Move) (toPruneSibs bool) {
		if pos.IsWhiteTurn {
			if res.Score > alpha {
				alpha = res.Score
				rtn.BestMove = move
				rtn.Score = res.Score
				if ALPHA_BETA_PRUNING_ENABLED && res.Score > beta {
					// black would not allow `pos`
					return true
				}
			}
		} else { // black turn
			if res.Score < beta {
				beta = res.Score
				rtn.BestMove = move
				rtn.Score = res.Score
				if ALPHA_BETA_PRUNING_ENABLED && res.Score < alpha {
					// white would not allow `pos`
					return true
				}
			}
		}
		return false
	}
	if MULTITHREADING_ENABLED {
		type searchResultMove struct {
			SearchResult *SearchResult
			Move         *chess.Move
		}
		results := make(chan *searchResultMove, len(moves))
		// Young Brothers Wait (YBW) scout
		handleResult(s.scoutToDepth(pos, depth), moves[0])
		for _, move := range moves {
			newPos := chess.GetBoardFromMove(pos, move)
			go func() {
				result := s.searchToDepth(newPos, depth-1, alpha, beta)
				results <- &searchResultMove{result, move}
			}()
		}
		var resultsCnt = 0
		for {
			var toBreak = false
			select {
			case res, ok := <-results:
				if !ok {
					toBreak = true
					break
				}
				handleResult(res.SearchResult, res.Move)
				resultsCnt++
				if resultsCnt == len(moves) {
					close(results)
				}
			}
			if toBreak {
				break
			}
		}
	} else { // multithreading not enabled
		for _, move := range moves {
			if s.IsHalted() {
				break
			}
			newPos := chess.GetBoardFromMove(pos, move)
			res := s.searchToDepth(newPos, depth-1, alpha, beta)
			toPruneSibs := handleResult(res, move)
			if toPruneSibs {
				break
			}
		}
	}

	results := rtn
	s.TT.PostResults(posHash, results, depth)

	return results
}

func (s *Search) scoutToDepth(pos *chess.Board, depth int) *SearchResult {
	if depth == 0 || pos.Result != chess.BOARD_RESULT_IN_PROGRESS {
		score := EvalPos(pos)
		return &SearchResult{BestMove: nil, Score: score}
	}
	posHash := ZobristHash(pos)
	var anticipatedMove *chess.Move
	ttEntry, _ := s.TT.GetEntry(posHash)
	if ttEntry != nil {
		anticipatedMove = ttEntry.Move
		if ttEntry.Depth >= depth {
			return &SearchResult{
				BestMove: ttEntry.Move,
				Score:    ttEntry.Score,
			}
		}
	}

	moves, movesErr := chess.GetLegalMoves(pos)
	if movesErr != nil {
		panic(fmt.Sprintf("could not get legal moves from pos %s: %s", pos, movesErr))
	}
	if MOVE_SORT_ENABLED {
		moves = SortMoves(pos, moves, anticipatedMove)
	}

	// TODO: take into account terminal boards here and in the main search
	newPos := chess.GetBoardFromMove(pos, moves[0])
	return s.scoutToDepth(newPos, depth-1)
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
