package main

import (
	"github.com/CameronHonis/chess"
	"math/rand"
)

type zobristLookups struct {
	PieceNums                  [8][8][13]int64 // 6 white pieces + 6 black pieces + 1 en passant type
	WhiteTurnNum               int64
	BlackTurnNum               int64
	CanWhiteKingsideCastleNum  int64
	CanBlackKingsdieCastleNum  int64
	CanWhiteQueensideCastleNum int64
	CanBlackQueensideCastleNum int64
}

func newZobristLookups() *zobristLookups {
	rand.Seed(0b101010101010101010101010101010101010101010101010101010101010101)
	zLookups := &zobristLookups{}
	for r := 0; r < 8; r++ {
		for c := 0; c < 8; c++ {
			for t := 0; t < 13; t++ {
				zLookups.PieceNums[r][c][t] = randInt64()
			}
		}
	}
	zLookups.WhiteTurnNum = randInt64()
	zLookups.BlackTurnNum = randInt64()
	zLookups.CanWhiteKingsideCastleNum = randInt64()
	zLookups.CanBlackKingsdieCastleNum = randInt64()
	zLookups.CanWhiteQueensideCastleNum = randInt64()
	zLookups.CanBlackQueensideCastleNum = randInt64()
	return zLookups
}

func randInt64() int64 {
	return rand.Int63() | (int64(rand.Intn(2)) << 63)
}

var lookups = newZobristLookups()

func ZobristHash(pos *chess.Board) int64 {
	var hash int64
	for r := 0; r < 8; r++ {
		for c := 0; c < 8; c++ {
			pieceType := pos.Pieces[r][c]
			if pieceType == chess.EMPTY {
				continue
			}
			pieceIdx := pos.Pieces[r][c] - 1
			hash ^= lookups.PieceNums[r][c][pieceIdx]
		}
	}
	if pos.OptEnPassantSquare != nil {
		r := pos.OptEnPassantSquare.Rank - 1
		c := pos.OptEnPassantSquare.File - 1
		hash ^= lookups.PieceNums[r][c][12]
	}
	if pos.CanWhiteCastleKingside {
		hash ^= lookups.CanWhiteKingsideCastleNum
	}
	if pos.CanWhiteCastleQueenside {
		hash ^= lookups.CanWhiteQueensideCastleNum
	}
	if pos.CanBlackCastleKingside {
		hash ^= lookups.CanBlackKingsdieCastleNum
	}
	if pos.CanBlackCastleQueenside {
		hash ^= lookups.CanBlackQueensideCastleNum
	}
	if pos.IsWhiteTurn {
		hash ^= lookups.WhiteTurnNum
	} else {
		hash ^= lookups.BlackTurnNum
	}
	return hash
}
