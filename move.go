package main

import (
	"fmt"
	"github.com/CameronHonis/chess"
	"log"
)

type MoveType uint8

const (
	NORMAL_MOVE MoveType = iota
	CAPTURES_EN_PASSANT
	PAWN_PROMOTION
	CASTLING
)

// Move is a compressed bit representation of a move
// The compress is important for minimizing memory requirements in the hash table
// The bit layout of a move is as follows:
// - bits 1-6 represent start square
// - bits 7-12 represent end square
// - bits 13-14 represent the promote piece:
//   - 0: Knight
//   - 1: Bishop
//   - 2: Rook
//   - 3: Queen
//
// - bits 15-16 specifies the move type:
//   - 0: a normal move
//   - 1: takes en passant
//   - 2: a promotion
type Move uint16

func NewMove(startSq, endSq, epSq Square, promoType PieceType, isCastles bool) Move {
	if DEBUG {
		if startSq > 0b111111 {
			log.Fatalf("cannot create move, start square occupies more than 6 bits: %b", startSq)
		}
		if endSq > 0b111111 {
			log.Fatalf("cannot create move, end square occupies more than 6 bits: %b", endSq)
		}
		if epSq != NULL_SQ && (epSq < SQ_A3 || epSq > SQ_H3) && (epSq < SQ_A6 || epSq > SQ_H6) {
			log.Fatalf("cannot create move, invalid en passant square %s", epSq.String())
		}
		if promoType == PAWN || promoType >= KING {
			log.Fatalf("cannot create move, invalid promo type: %d", promoType)
		}
	}
	var moveTypeBits Move
	if promoType != EMPTY_PIECE_TYPE {
		moveTypeBits = 0b10
	} else if endSq == epSq {
		moveTypeBits = 0b01
	} else if isCastles {
		moveTypeBits = 0b11
	}

	var promoBits Move
	if promoType == KNIGHT {
		promoBits = 0b0000
	} else if promoType == BISHOP {
		promoBits = 0b0100
	} else if promoType == ROOK {
		promoBits = 0b1000
	} else if promoType == QUEEN {
		promoBits = 0b1100
	}

	return Move(startSq&0b111111)<<10 | Move(endSq&0b111111)<<4 | promoBits | moveTypeBits
}

func NewNormalMove(startSq, endSq Square) Move {
	return NewMove(startSq, endSq, NULL_SQ, EMPTY_PIECE_TYPE, false)
}

func NewEnPassantMove(startSq, endSq Square) Move {
	return NewMove(startSq, endSq, endSq, EMPTY_PIECE_TYPE, false)
}

func NewPromoteMoves(startSq Square, endSq Square) []Move {
	return []Move{
		NewMove(startSq, endSq, NULL_SQ, KNIGHT, false),
		NewMove(startSq, endSq, NULL_SQ, BISHOP, false),
		NewMove(startSq, endSq, NULL_SQ, ROOK, false),
		NewMove(startSq, endSq, NULL_SQ, QUEEN, false),
	}
}

func LegalMoveFromLegacyMove(legacyMove *chess.Move, pos *Position) (Move, error) {
	iter := NewLegalMoveIter(pos)
	for {
		move, done := iter.Next()
		if done {
			break
		}
		if move.StartSq().Rank() != legacyMove.StartSquare.Rank {
			continue
		}
		if move.StartSq().File() != legacyMove.StartSquare.File {
			continue
		}
		if move.EndSq().Rank() != legacyMove.EndSquare.Rank {
			continue
		}
		if move.EndSq().File() != legacyMove.EndSquare.File {
			continue
		}
		if move.PromotedTo() == EMPTY_PIECE_TYPE && legacyMove.PawnUpgradedTo != chess.EMPTY {
			continue
		}
		if move.PromotedTo() == KNIGHT && !(legacyMove.PawnUpgradedTo == chess.WHITE_KNIGHT || legacyMove.PawnUpgradedTo == chess.BLACK_KNIGHT) {
			continue
		}
		if move.PromotedTo() == BISHOP && !(legacyMove.PawnUpgradedTo == chess.WHITE_BISHOP || legacyMove.PawnUpgradedTo == chess.BLACK_BISHOP) {
			continue
		}
		if move.PromotedTo() == ROOK && !(legacyMove.PawnUpgradedTo == chess.WHITE_ROOK || legacyMove.PawnUpgradedTo == chess.BLACK_ROOK) {
			continue
		}
		if move.PromotedTo() == QUEEN && !(legacyMove.PawnUpgradedTo == chess.WHITE_QUEEN || legacyMove.PawnUpgradedTo == chess.BLACK_QUEEN) {
			continue
		}
		return move, nil
	}
	return NULL_MOVE, fmt.Errorf("could not find legal move %s in pos:\n%s", legacyMove.ToLongAlgebraic(), pos)
}

func (m Move) StartSq() Square {
	return Square(m >> 10)
}

func (m Move) EndSq() Square {
	return Square((m >> 4) & 0b111111)
}

func (m Move) PromotedTo() PieceType {
	if m.Type() != PAWN_PROMOTION {
		return EMPTY_PIECE_TYPE
	}
	return PieceType((m>>2)&0b11) + KNIGHT
}

func (m Move) Type() MoveType {
	return MoveType(m & 0b11)
}

func (m Move) IsNull() bool {
	return m == 0
}

func (m Move) String() string {
	if m.PromotedTo() == KNIGHT {
		return fmt.Sprintf("%s%sn", m.StartSq(), m.EndSq())
	} else if m.PromotedTo() == BISHOP {
		return fmt.Sprintf("%s%sb", m.StartSq(), m.EndSq())
	} else if m.PromotedTo() == ROOK {
		return fmt.Sprintf("%s%sr", m.StartSq(), m.EndSq())
	} else if m.PromotedTo() == QUEEN {
		return fmt.Sprintf("%s%sq", m.StartSq(), m.EndSq())
	} else {
		return fmt.Sprintf("%s%s", m.StartSq(), m.EndSq())
	}
}

const NULL_MOVE = Move(0)
