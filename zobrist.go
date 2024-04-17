package main

import (
	"github.com/CameronHonis/chess"
	"math/rand"
)

type zobristLookupsLegacy struct {
	PieceKeys                  [N_SQUARES][N_PIECES]uint64 // 6 white pieces + 6 black pieces + 1 en passant type
	EpSqKeys                   [N_SQUARES]uint64
	WhiteTurnKey               uint64
	BlackTurnKey               uint64
	CanWhiteKingsideCastleKey  uint64
	CanBlackKingsdieCastleKey  uint64
	CanWhiteQueensideCastleKey uint64
	CanBlackQueensideCastleKey uint64
}

func newZobristLookups() *zobristLookupsLegacy {
	rand.Seed(0b101010101010101010101010101010101010101010101010101010101010101)
	zLookups := &zobristLookupsLegacy{}
	for sq := Square(0); sq < N_SQUARES; sq++ {
		for piece := Piece(0); piece < N_PIECES; piece++ {
			zLookups.PieceKeys[sq][piece] = rand.Uint64()
		}
		zLookups.EpSqKeys[sq] = rand.Uint64()
	}
	zLookups.WhiteTurnKey = rand.Uint64()
	zLookups.BlackTurnKey = rand.Uint64()
	zLookups.CanWhiteKingsideCastleKey = rand.Uint64()
	zLookups.CanBlackKingsdieCastleKey = rand.Uint64()
	zLookups.CanWhiteQueensideCastleKey = rand.Uint64()
	zLookups.CanBlackQueensideCastleKey = rand.Uint64()
	return zLookups
}

func randInt64() int64 {
	return rand.Int63() | (int64(rand.Intn(2)) << 63)
}

var lookups = newZobristLookups()

func ZobristHashOnLegacyBoard(pos *chess.Board) uint64 {
	var hash uint64
	for sq := Square(0); sq < N_SQUARES; sq++ {
		piece := pos.Pieces[sq.Rank()-1][sq.File()-1]
		hash ^= lookups.PieceKeys[sq][piece]
	}
	if pos.OptEnPassantSquare != nil {
		epSq := SqFromCoords(int(pos.OptEnPassantSquare.Rank), int(pos.OptEnPassantSquare.File))
		hash ^= lookups.EpSqKeys[epSq]
	}
	if pos.CanWhiteCastleKingside {
		hash ^= lookups.CanWhiteKingsideCastleKey
	}
	if pos.CanWhiteCastleQueenside {
		hash ^= lookups.CanWhiteQueensideCastleKey
	}
	if pos.CanBlackCastleKingside {
		hash ^= lookups.CanBlackKingsdieCastleKey
	}
	if pos.CanBlackCastleQueenside {
		hash ^= lookups.CanBlackQueensideCastleKey
	}
	if pos.IsWhiteTurn {
		hash ^= lookups.WhiteTurnKey
	} else {
		hash ^= lookups.BlackTurnKey
	}
	return hash
}

//type ZHash uint64
//
//func NewZHash(pos *Position) {
//	var hash ZHash
//	for sq := Square(0); sq < N_SQUARES; sq++ {
//
//	}
//}
