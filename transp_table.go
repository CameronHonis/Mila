package main

import (
	"log"
	"sync"
)

type TranspTable struct {
	entryByHash map[ZHash]TTEntry
	mu          sync.Mutex
}

func NewTranspTable() *TranspTable {
	return &TranspTable{
		entryByHash: make(map[ZHash]TTEntry),
	}
}

func (tt *TranspTable) PostResults(hash ZHash, score int16, move Move, depth uint8) {
	tt.mu.Lock()
	defer tt.mu.Unlock()
	prevEntry, ok := tt.entryByHash[hash]
	if !ok || depth > prevEntry.Depth {
		tt.entryByHash[hash] = TTEntry{
			Score: score,
			Depth: depth,
			Move:  move,
		}
	}
}

func (tt *TranspTable) GetEntry(hash ZHash) (entry TTEntry, exists bool) {
	tt.mu.Lock()
	defer tt.mu.Unlock()
	var ok bool
	if entry, ok = tt.entryByHash[hash]; ok {
		return entry, true
	}
	return
}

func (tt *TranspTable) Line(pos *Position, depth uint8) []Move {
	entry, entryExists := tt.GetEntry(pos.hash)
	if DEBUG {
		if !entryExists {
			log.Fatalf("could not get entry for line after %s at depth %d", pos.FEN(), depth)
		}
		if entry.Depth < depth {
			log.Fatalf("entry depth (%d) lower than requested depth (%d) while building line", entry.Depth, depth)
		}
	}
	if depth == 1 {
		return []Move{entry.Move}
	} else {
		captPiece, lastFrozenPos := pos.MakeMove(entry.Move)
		defer pos.UnmakeMove(entry.Move, lastFrozenPos, captPiece)

		moves := tt.Line(pos, depth-1)
		return append(moves, entry.Move)
	}
}

type TTEntry struct {
	Score int16
	Depth uint8
	Move  Move
}
