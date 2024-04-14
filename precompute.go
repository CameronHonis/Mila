package main

import (
	"github.com/CameronHonis/set"
)

// N_ROW_PERMUS is the number of possible permutations of occupied squares
// for any row (2^8). This number happens to be the permutations of occupied
// squares for columns or diagonals as well.
const N_ROW_PERMUS = 1 << 8

// Splitting the precomputed rook moves into rows and columns as separate look tables
// reduces the overall amount of entries in the set. This works out since a rooks moves
// on a given row have no dependence to its moves on a file, and vice versa.

var RankAttacks map[Bitboard]Bitboard
var FileAttacks map[Bitboard]Bitboard
var PosDiagAttacks map[Bitboard]Bitboard
var NegDiagAttacks map[Bitboard]Bitboard

func Init() {

}

func genAllRankOccupiedBBs() *set.Set[Bitboard] {
	rtn := set.EmptySet[Bitboard]()
	for rank := uint8(1); rank <= N_RANKS; rank++ {
		for occupied := 1; occupied < N_ROW_PERMUS; occupied++ {
			bb := BBWithRank(rank, uint8(occupied))
			rtn.Add(bb)
		}
	}
	return rtn
}

func genAllFileOccupiedBBs() *set.Set[Bitboard] {
	rtn := set.EmptySet[Bitboard]()
	for file := uint8(1); file <= N_FILES; file++ {
		for occupied := 1; occupied < N_ROW_PERMUS; occupied++ {
			bb := BBWithFile(file, uint8(occupied))
			rtn.Add(bb)
		}
	}
	return rtn
}

func genAllPosDiagOccupiedBBs() *set.Set[Bitboard] {
	rtn := set.EmptySet[Bitboard]()
	for posDiagIdx := uint8(0); posDiagIdx < N_DIAGS; posDiagIdx++ {
		var diagSize uint8
		if posDiagIdx <= 7 {
			diagSize = 8 - (7 - posDiagIdx)
		} else {
			diagSize = 8 - (posDiagIdx - 7)
		}
		nDiagPerms := 1 << diagSize
		for occupied := 1; occupied < nDiagPerms; occupied++ {
			bb := BBWithPosDiag(posDiagIdx, uint8(occupied))
			rtn.Add(bb)
		}
	}
	return rtn
}

func genAllNegDiagOccupiedBBs() *set.Set[Bitboard] {
	rtn := set.EmptySet[Bitboard]()
	for negDiagIdx := uint8(0); negDiagIdx < N_DIAGS; negDiagIdx++ {
		var diagSize uint8
		if negDiagIdx <= 7 {
			diagSize = 8 - (7 - negDiagIdx)
		} else {
			diagSize = 8 - (negDiagIdx - 7)
		}
		nDiagPerms := 1 << diagSize
		for occupied := 1; occupied < nDiagPerms; occupied++ {
			bb := BBWithNegDiag(negDiagIdx, uint8(occupied))
			rtn.Add(bb)
		}
	}
	return rtn
}
