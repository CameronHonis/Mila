package main

import (
	"fmt"
	"github.com/CameronHonis/chess"
	"github.com/CameronHonis/marker"
	"math"
	"time"
)

const ALPHA_BETA_PRUNING_ENABLED = true
const MOVE_SORT_ENABLED = true
const TRANSP_TABLE_LOOKUPS_ENABLED = true

const PRINT_LINE = true
const PRINT_HASH_HITS = true
const PRINT_PRUNE_CNTS = true
const PRINT_NODE_CNT = true
const PRINT_TIME = true

type Search struct {
	__static__  marker.Marker
	Root        *chess.Board
	TT          *TranspTable
	Constraints *SearchConstraints

	__controls__ marker.Marker
	isHalted     bool

	__ephemeral__    marker.Marker
	nodeCntOnDepth   int
	hashHitsOnDepth  int
	pruneCntsOnDepth []int

	__accumulated__ marker.Marker
	depth           int
	score           float64
	accNodeCnt      int
}

func NewSearch(pos *chess.Board, constraints *SearchConstraints, tt *TranspTable) *Search {
	return &Search{
		Root:             pos,
		TT:               tt,
		Constraints:      constraints,
		isHalted:         true,
		pruneCntsOnDepth: make([]int, 0),
	}
}

func (s *Search) IncrNode() {
	s.accNodeCnt++
	s.nodeCntOnDepth++
	if s.accNodeCnt >= s.Constraints.NodeCntLmt() {
		fmt.Println("halting search, max node count reached")
		s.isHalted = true
	}
}

func (s *Search) ToNextDepth() {
	s.depth++
	s.pruneCntsOnDepth = make([]int, s.depth)
	s.nodeCntOnDepth = 0
	s.hashHitsOnDepth = 0
	if s.depth > s.Constraints.DepthLmt() {
		fmt.Println("halting search, max depth reached")
		s.isHalted = true
	}
}

func (s *Search) TallyPrune(parDepth, cnt int) {
	idx := len(s.pruneCntsOnDepth) - parDepth
	s.pruneCntsOnDepth[idx] += cnt
}

func (s *Search) Start() {
	s.isHalted = false
	maxSearchMs := s.MaxSearchMs()
	go func() {
		time.Sleep(time.Duration(maxSearchMs) * time.Millisecond)
		fmt.Println("halting search, search time allowance reached")
		s.isHalted = true
	}()
	var lastResultTime = time.Now()
	var line []*chess.Move
	for {
		s.ToNextDepth()

		if s.isHalted {
			break
		}

		score := s.searchToDepth(s.Root, s.depth, -math.MaxFloat64, math.MaxFloat64)
		dt := time.Now().Sub(lastResultTime)
		lastResultTime = time.Now()
		out := fmt.Sprintf("info depth %d score %f", s.depth, score)

		line = s.BestLine(s.Root)
		if PRINT_LINE {
			var lineStr string
			for moveIdx, move := range line {
				if moveIdx == 0 {
					lineStr += move.ToLongAlgebraic()
				} else {
					lineStr += " " + move.ToLongAlgebraic()
				}
			}
			out += fmt.Sprintf(" moves %s", lineStr)
		} else {
			out += fmt.Sprintf(" move %s", line[0].ToLongAlgebraic())
		}

		if PRINT_NODE_CNT {
			out += fmt.Sprintf(" nodes %d", s.nodeCntOnDepth)
		}

		if PRINT_HASH_HITS {
			out += fmt.Sprintf(" hits %d", s.hashHitsOnDepth)
		}

		if PRINT_PRUNE_CNTS {
			out += " pruned"
			for _, pruneCnt := range s.pruneCntsOnDepth {
				out += fmt.Sprintf(" %d", pruneCnt)
			}
		}

		if PRINT_TIME {
			out += fmt.Sprintf(" time %dms", dt.Milliseconds())
		}
		fmt.Println(out)
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

	posHash := ZobristHashOnLegacyBoard(pos)
	var anticipatedMove *chess.Move
	if TRANSP_TABLE_LOOKUPS_ENABLED {
		if ttEntry, _ := s.TT.GetEntry(posHash); ttEntry != nil {
			s.hashHitsOnDepth++
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
	var bestScore float64
	for moveIdx, move := range moves {
		if s.isHalted {
			break
		}
		newPos := chess.GetBoardFromMove(pos, move)
		score := s.searchToDepth(newPos, depth-1, alpha, beta)
		if pos.IsWhiteTurn {
			if score > alpha {
				alpha = score
				bestScore = alpha
				bestMove = move
				if ALPHA_BETA_PRUNING_ENABLED && score >= beta {
					// black would not allow `pos`
					prunedCnt := len(moves) - moveIdx - 1
					s.TallyPrune(depth, prunedCnt)
					break
				}
			}
		} else { // black turn
			if score < beta {
				beta = score
				bestScore = beta
				bestMove = move
				if ALPHA_BETA_PRUNING_ENABLED && score <= alpha {
					// white would not allow `pos`
					prunedCnt := len(moves) - moveIdx - 1
					s.TallyPrune(depth, prunedCnt)
					break
				}
			}
		}
	}
	s.TT.PostResults(posHash, bestScore, bestMove, depth)

	return bestScore
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
		ttEntry, _ = s.TT.GetEntry(ZobristHashOnLegacyBoard(pos))
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
