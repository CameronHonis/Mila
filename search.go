package main

import (
	"fmt"
	"github.com/CameronHonis/marker"
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
	Root        *Position
	TT          *TranspTable
	Constraints *SearchConstraints

	__controls__ marker.Marker
	isHalted     bool

	__ephemeral__    marker.Marker
	nodeCntOnDepth   int
	hashHitsOnDepth  int
	pruneCntsOnDepth []int

	__accumulated__ marker.Marker
	depth           uint8
	score           float64
	accNodeCnt      int
}

func NewSearch(pos *Position, constraints *SearchConstraints, tt *TranspTable) *Search {
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
	var line []Move
	for {
		s.ToNextDepth()

		if s.isHalted {
			break
		}

		var score int16
		score, line = s.searchToDepth(s.Root, s.depth)
		dt := time.Now().Sub(lastResultTime)
		lastResultTime = time.Now()
		var out = fmt.Sprintf("info depth %d score ", s.depth)
		if score == -MATE_VAL || score == MATE_VAL {
			out += fmt.Sprintf("mate %d", s.depth)
		} else {
			out += fmt.Sprintf("%d", score)
		}

		if PRINT_LINE {
			var lineStr string
			for moveIdx, move := range line {
				if moveIdx == 0 {
					lineStr += move.String()
				} else {
					lineStr += " " + move.String()
				}
			}
			out += fmt.Sprintf(" moves %s", lineStr)
		} else {
			out += fmt.Sprintf(" move %s", line[0].String())
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

		if score == -MATE_VAL || score == MATE_VAL {
			break
		}
	}
	if len(line) > 0 {
		fmt.Printf("bestmove %s\n", line[0].String())
	}

}

func (s *Search) searchToDepth(pos *Position, depth uint8) (score int16, line []Move) {
	if !pos.HasLegalMoves() {
		if pos.IsMate() {
			return -MATE_VAL, make([]Move, 0)
		} else {
			return DRAW_VAL, make([]Move, 0)
		}
	}
	score, line = s._searchToDepth(pos, depth, -MATE_VAL, MATE_VAL)
	ReverseSlice(line)
	return score, line
}

func (s *Search) _searchToDepth(pos *Position, depth uint8, alpha int16, beta int16) (score int16, line []Move) {
	if depth == 0 || pos.result != RESULT_IN_PROGRESS {
		if pos.IsMate() {
			return -MATE_VAL, nil
		}
		s.IncrNode()
		return EvalPos(pos), nil
	}

	var anticipated = NULL_MOVE
	if TRANSP_TABLE_LOOKUPS_ENABLED {
		entry, exists := s.TT.GetEntry(pos.hash)
		if exists {
			s.hashHitsOnDepth++
			if entry.Depth >= depth {
				return entry.Score, s.TT.Line(pos, depth)
			} else {
				anticipated = entry.Move
			}
		}
	}

	iter := NewLegalMoveIter(pos)
	if MOVE_SORT_ENABLED {
		iter.pMoves = SortMoves(pos, iter.pMoves, anticipated)
	}

	score = alpha
	for {
		move, done := iter.Next()
		if done {
			break
		}

		captPiece, lastFrozenPos := pos.MakeMove(move)
		moveScore, moveLine := s._searchToDepth(pos, depth-1, -beta, -alpha)
		moveScore = -moveScore
		pos.UnmakeMove(move, lastFrozenPos, captPiece)

		if moveScore > alpha {
			alpha = moveScore
			score = moveScore
			if moveLine == nil {
				line = []Move{move}
			} else {
				line = append(moveLine, move)
			}

			if TRANSP_TABLE_LOOKUPS_ENABLED {
				s.TT.PostResults(pos.hash, score, move, depth)
			}
		}

		if ALPHA_BETA_PRUNING_ENABLED && moveScore >= beta {
			s.TallyPrune(int(depth), len(iter.pMoves)-1-iter.idx)
			return
		}
	}

	return
}

func (s *Search) MaxSearchMs() int {
	var msForSearch = func(pos *Position, bankMs int, incrMs int) int {
		expMoves := ExpMoves(pos)
		return incrMs + bankMs/expMoves
	}
	maxSearchMs := s.Constraints.MaxSearchMs()
	if s.Root.isWhiteTurn && s.Constraints.whiteMs > 0 {
		maxSearchMs = MinInt(maxSearchMs, msForSearch(s.Root, s.Constraints.whiteMs, s.Constraints.whiteIncrMs))
	} else if !s.Root.isWhiteTurn && s.Constraints.blackMs > 0 {
		maxSearchMs = MinInt(maxSearchMs, msForSearch(s.Root, s.Constraints.blackMs, s.Constraints.blackIncrMs))
	}
	return maxSearchMs
}
