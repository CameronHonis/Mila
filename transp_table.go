package main

import (
	"github.com/CameronHonis/chess"
	"sync"
)

type TranspTable struct {
	entryByHash map[int64]*TTEntry
	mu          sync.Mutex
}

func NewTranspTable() *TranspTable {
	return &TranspTable{
		entryByHash: make(map[int64]*TTEntry),
	}
}

func (tt *TranspTable) PostResults(hash int64, score float64, move *chess.Move, depth int) {
	tt.mu.Lock()
	defer tt.mu.Unlock()
	prevEntry := tt.entryByHash[hash]
	if prevEntry == nil || depth > prevEntry.Depth {
		tt.entryByHash[hash] = &TTEntry{
			Score: score,
			Depth: depth,
			Move:  move,
		}
	}
}

func (tt *TranspTable) GetEntry(hash int64) (entry *TTEntry, exists bool) {
	tt.mu.Lock()
	defer tt.mu.Unlock()
	var ok bool
	if entry, ok = tt.entryByHash[hash]; ok {
		return entry, true
	}
	return
}

type TTEntry struct {
	Score float64
	Depth int
	Move  *chess.Move
	// bestMove *chess.Move
	mu sync.Mutex
}

func (e *TTEntry) Update(score float64, depth int, move *chess.Move) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.Score = score
	e.Depth = depth
	e.Move = move
}
