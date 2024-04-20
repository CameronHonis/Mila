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
	if !ok || depth > prevEntry.Depth || (depth == prevEntry.Depth && score > prevEntry.Score) {
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
	nMoves := depth
	frozenPoss := make([]*FrozenPos, nMoves)
	capturedPieces := make([]Piece, nMoves)
	line := make([]Move, nMoves)
	for moveIdx := uint8(0); moveIdx < nMoves; moveIdx++ {
		entry, entryExists := tt.GetEntry(pos.hash)
		if DEBUG {
			if !entryExists {
				log.Fatalf("could not get entry for line after %s at depth %d", pos.FEN(), depth)
			}
			if entry.Depth < depth {
				log.Fatalf("entry depth (%d) lower than requested depth (%d) while building line", entry.Depth, depth)
			}
		}
		line[moveIdx] = entry.Move
		isLastMove := moveIdx == nMoves-1
		if !isLastMove {
			capturedPieces[moveIdx], frozenPoss[moveIdx] = pos.MakeMove(entry.Move)
		}
	}

	//undo moves on pos
	for moveIdx := int(nMoves) - 2; moveIdx >= 0; moveIdx-- {
		move := line[moveIdx]
		captPiece := capturedPieces[moveIdx]
		frozenPos := frozenPoss[moveIdx]
		pos.UnmakeMove(move, frozenPos, captPiece)
	}
	return line
}

type TTEntry struct {
	Score int16
	Depth uint8
	Move  Move
}
