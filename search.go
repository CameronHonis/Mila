package main

import (
	"context"
	"fmt"
	"github.com/CameronHonis/chess"
	"github.com/CameronHonis/marker"
	"log"
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

type SearchNode struct {
	__static__  marker.Marker
	Parent      *SearchNode
	Pos         *chess.Board
	Depth       uint8
	SiblingRank uint8

	__dynamic__ marker.Marker // plan is to track threads "on" or "downstream" this node and only require mutex locks
	// when two or more threads could potentially mutate the same node. This should be avoided to the highest degree
	// though with a smart ordering of nodes in the Edge
	mu           sync.Mutex
	alpha        float64
	beta         float64
	move         *chess.Move
	cpuCnt       uint8
	childCnt     uint8
	childDoneCnt uint8
	didExpand    bool
}

func (se *SearchNode) ComesBefore(_other Sortable) bool {
	other, ok := _other.(*SearchNode)
	if !ok {
		log.Fatalf("cannot compare a SearchNode to a %s", _other)
	}
	if se.SiblingRank != other.SiblingRank {
		return se.SiblingRank < other.SiblingRank
	}
	if se.Depth != other.Depth {
		return se.Depth < other.Depth
	}
	return ZobristHash(se.Pos) < ZobristHash(other.Pos)
}

type Edge struct {
	PriorityQueue[*SearchNode]
	edge []*SearchNode
}

func (e *Edge) Push(node *SearchNode) {
	e.edge = append(e.edge, node)
}

func (e *Edge) Pop() *SearchNode {
	if len(e.edge) == 0 {
		return nil
	}
	last := e.edge[len(e.edge)-1]
	e.edge = e.edge[:len(e.edge)-1]
	return last
}

func (e *Edge) Len() uint {
	return uint(len(e.edge))
}

type Search struct {
	__static__  marker.Marker
	Root        *chess.Board
	TT          *TranspTable
	Constraints *SearchConstraints
	Depth       uint8
	Ctx         context.Context // root context
	CancelCtx   context.CancelFunc
	Edge        Edge

	__dynamic__ marker.Marker
	nodeCnt     int
	score       float64
	mu          sync.Mutex
}

func NewSearch(pos *chess.Board, constraints *SearchConstraints, tt *TranspTable) *Search {
	return &Search{
		Root:        pos,
		TT:          tt,
		Constraints: constraints,
	}
}

func (s *Search) IncrNode() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nodeCnt++
	if s.nodeCnt >= s.Constraints.NodeCntLmt() {
		fmt.Println("halting search, max node count reached")
		s.CancelCtx()
	}
}

func (s *Search) IncrDepth() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Depth++
	if s.Depth > s.Constraints.DepthLmt() {
		fmt.Println("halting search, max depth reached")
		s.CancelCtx()
	}
}

func (s *Search) NodeCnt() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.nodeCnt
}

func (s *Search) Start() {
	maxSearchMs := s.MaxSearchMs()
	var bestMove *chess.Move
	var lastResultTime = time.Now()
	s.Ctx, s.CancelCtx = context.WithTimeout(context.Background(), time.Duration(maxSearchMs)*time.Millisecond)
	for {
		s.IncrDepth()

		if err := s.Ctx.Err(); err != nil {
			fmt.Println("halting search, context error: ", err)
			break
		}

		results := s.searchToDepth(s.Root, s.Depth)
		bestMove = results.BestMove
		dt := time.Now().Sub(lastResultTime)
		lastResultTime = time.Now()
		var lineStr = ""
		for moveIdx, move := range s.BestLine() {
			if moveIdx == 0 {
				lineStr += move.ToLongAlgebraic()
			} else {
				lineStr += " " + move.ToLongAlgebraic()
			}
		}
		fmt.Printf("info depth %d score %f line %s nodes %d time %dms\n", s.Depth, results.Score, lineStr, s.NodeCnt(), dt.Milliseconds())
	}
	if bestMove != nil {
		fmt.Printf("bestmove %s\n", bestMove.ToLongAlgebraic())
	}

}

func (s *Search) searchToDepth(pos *chess.Board, depth uint8) *SearchResult {
	s.Edge.Push(&SearchNode{Pos: pos, Depth: depth, SiblingRank: 0, alpha: -math.MaxFloat64, beta: math.MaxFloat64})
	for s.Edge.Len() > 0 {
		ele := s.Edge.Pop()
		s.searchNode(ele)
	}
}

func (s *Search) searchNode(node *SearchNode) {
	posHash := ZobristHash(node.Pos)

	var pvMove *chess.Move
	if entry, _ := s.TT.GetEntry(posHash); entry != nil {
		if entry.Depth >= node.Depth {
			propagateSearchResults(node, &SearchResult{BestMove: entry.Move, Score: entry.Score})
			return
		} else {
			pvMove = entry.Move
		}
	}
	if node.Depth == 0 {

	}
}

func propagateSearchResults(node *SearchNode, results *SearchResult) {
	if results.Score > node.alpha {
		node.alpha = results.Score
		node.move = results.BestMove
		if results.Score > node.beta {

		}
	}
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

func (s *Search) BestLine() []*chess.Move {
	s.mu.Lock()
	defer s.mu.Unlock()
	pos := s.Root
	ttEntry, _ := s.TT.GetEntry(ZobristHash(pos))
	rtn := make([]*chess.Move, 0)
	for {
		if ttEntry == nil {
			break
		}
		if ttEntry.Move == nil {
			break
		}
		rtn = append(rtn, ttEntry.Move)
		pos = chess.GetBoardFromMove(pos, ttEntry.Move)
		ttEntry, _ = s.TT.GetEntry(ZobristHash(pos))
	}
	return rtn
}
