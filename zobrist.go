package main

import (
	"github.com/CameronHonis/chess"
	"math/rand"
)

type zobristLookupsLegacy struct {
	PieceKeys                  [N_SQUARES][N_PIECES - 1]ZHash // 6 white pieces + 6 black pieces + 1 en passant type
	EpSqKeys                   [N_SQUARES]ZHash
	WhiteTurnKey               ZHash
	BlackTurnKey               ZHash
	CanWhiteKingsideCastleKey  ZHash
	CanBlackKingsideCastleKey  ZHash
	CanWhiteQueensideCastleKey ZHash
	CanBlackQueensideCastleKey ZHash
}

func newZobristLookups() *zobristLookupsLegacy {
	rand.Seed(0b101010101010101010101010101010101010101010101010101010101010101)
	zLookups := &zobristLookupsLegacy{}
	for sq := Square(0); sq < N_SQUARES; sq++ {
		for piece := Piece(0); piece < N_PIECES-1; piece++ {
			zLookups.PieceKeys[sq][piece] = ZHash(rand.Uint64())
		}
		zLookups.EpSqKeys[sq] = ZHash(rand.Uint64())
	}
	zLookups.WhiteTurnKey = ZHash(rand.Uint64())
	zLookups.BlackTurnKey = ZHash(rand.Uint64())
	zLookups.CanWhiteKingsideCastleKey = ZHash(rand.Uint64())
	zLookups.CanBlackKingsideCastleKey = ZHash(rand.Uint64())
	zLookups.CanWhiteQueensideCastleKey = ZHash(rand.Uint64())
	zLookups.CanBlackQueensideCastleKey = ZHash(rand.Uint64())
	return zLookups
}

var lookups = newZobristLookups()

func ZobristHashOnLegacyBoard(pos *chess.Board) uint64 {
	var hash ZHash
	for sq := SQ_A1; sq < N_SQUARES; sq++ {
		piece := pos.Pieces[sq.Rank()-1][sq.File()-1]
		if piece == chess.EMPTY {
			continue
		}
		hash ^= lookups.PieceKeys[sq][piece-1]
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
		hash ^= lookups.CanBlackKingsideCastleKey
	}
	if pos.CanBlackCastleQueenside {
		hash ^= lookups.CanBlackQueensideCastleKey
	}
	if pos.IsWhiteTurn {
		hash ^= lookups.WhiteTurnKey
	} else {
		hash ^= lookups.BlackTurnKey
	}
	return uint64(hash)
}

type ZHash uint64

func NewZHash(pos *Position) ZHash {
	var hash ZHash
	for sq := Square(0); sq < N_SQUARES; sq++ {
		piece := pos.pieces[sq]
		if piece == EMPTY {
			continue
		}
		hash ^= lookups.PieceKeys[sq][piece-1]
	}
	fPos := pos.frozenPos
	if fPos.EnPassantSq != NULL_SQ {
		hash ^= lookups.EpSqKeys[fPos.EnPassantSq]
	}
	if fPos.CastleRights[W_CASTLE_KINGSIDE_RIGHT] {
		hash ^= lookups.CanWhiteKingsideCastleKey
	}
	if fPos.CastleRights[W_CASTLE_QUEENSIDE_RIGHT] {
		hash ^= lookups.CanWhiteQueensideCastleKey
	}
	if fPos.CastleRights[B_CASTLE_KINGSIDE_RIGHT] {
		hash ^= lookups.CanBlackKingsideCastleKey
	}
	if fPos.CastleRights[B_CASTLE_QUEENSIDE_RIGHT] {
		hash ^= lookups.CanBlackQueensideCastleKey
	}
	if pos.isWhiteTurn {
		hash ^= lookups.WhiteTurnKey
	} else {
		hash ^= lookups.BlackTurnKey
	}
	return hash
}

func (zh ZHash) UpdatePieceOnSq(prevPiece, piece Piece, sq Square) ZHash {
	if prevPiece != EMPTY {
		zh ^= lookups.PieceKeys[sq][prevPiece]
	}
	if piece != EMPTY {
		zh ^= lookups.PieceKeys[sq][piece]
	}
	return zh
}

func (zh ZHash) UpdateEnPassantSq(prevEpSq, epSq Square) ZHash {
	if prevEpSq != NULL_SQ {
		zh ^= lookups.EpSqKeys[prevEpSq]
	}
	if epSq != NULL_SQ {
		zh ^= lookups.EpSqKeys[epSq]
	}
	return zh
}

// ToggleCastleRight is a generalized version of RemoveCastleRight, as it's important
// to add castle rights back when updating the Position.hash while unmaking a move.
func (zh ZHash) ToggleCastleRight(castleRight CastleRight) ZHash {
	if castleRight == W_CASTLE_KINGSIDE_RIGHT {
		zh ^= lookups.CanWhiteKingsideCastleKey
	} else if castleRight == W_CASTLE_QUEENSIDE_RIGHT {
		zh ^= lookups.CanWhiteQueensideCastleKey
	} else if castleRight == B_CASTLE_KINGSIDE_RIGHT {
		zh ^= lookups.CanBlackKingsideCastleKey
	} else if castleRight == B_CASTLE_QUEENSIDE_RIGHT {
		zh ^= lookups.CanBlackQueensideCastleKey
	}
	return zh
}

func (zh ZHash) ToggleTurn() ZHash {
	zh ^= lookups.WhiteTurnKey
	zh ^= lookups.BlackTurnKey
	return zh
}
