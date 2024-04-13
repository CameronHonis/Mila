package main

import (
	"fmt"
	"math"
)

type PieceType uint8

const (
	EMPTY_PIECE_TYPE PieceType = iota
	PAWN
	KNIGHT
	BISHOP
	ROOK
	QUEEN
	KING
	N_PIECE_TYPES
)

type Piece uint8

const (
	EMPTY Piece = iota
	W_PAWN
	W_KNIGHT
	W_BISHOP
	W_ROOK
	W_QUEEN
	W_KING
	B_PAWN
	B_KNIGHT
	B_BISHOP
	B_ROOK
	B_QUEEN
	B_KING
	N_PIECES
)

func PieceFromChar(char byte) Piece {
	if char == 'p' {
		return B_PAWN
	} else if char == 'n' {
		return B_KNIGHT
	} else if char == 'b' {
		return B_BISHOP
	} else if char == 'r' {
		return B_ROOK
	} else if char == 'q' {
		return B_QUEEN
	} else if char == 'k' {
		return B_KING
	} else if char == 'P' {
		return W_PAWN
	} else if char == 'N' {
		return W_KNIGHT
	} else if char == 'B' {
		return W_BISHOP
	} else if char == 'R' {
		return W_ROOK
	} else if char == 'Q' {
		return W_QUEEN
	} else if char == 'K' {
		return W_KING
	} else {
		return EMPTY
	}
}

func (p Piece) IsWhite() bool {
	return p >= W_PAWN && p <= W_KING
}

func (p Piece) String() string {
	switch p {
	case EMPTY:
		return " "
	case W_PAWN:
		return "♟︎"
	case W_KNIGHT:
		return "♞"
	case W_BISHOP:
		return "♝"
	case W_ROOK:
		return "♜"
	case W_QUEEN:
		return "♛"
	case W_KING:
		return "♚"
	case B_PAWN:
		return "♙"
	case B_KNIGHT:
		return "♘"
	case B_BISHOP:
		return "♗"
	case B_ROOK:
		return "♖"
	case B_QUEEN:
		return "♕"
	case B_KING:
		return "♔"
	default:
		return "?"
	}
}

func (p Piece) Char() byte {
	if p == EMPTY {
		return ' '
	} else if p == W_PAWN {
		return 'P'
	} else if p == W_KNIGHT {
		return 'N'
	} else if p == W_BISHOP {
		return 'B'
	} else if p == W_ROOK {
		return 'R'
	} else if p == W_QUEEN {
		return 'Q'
	} else if p == W_KING {
		return 'K'
	} else if p == B_PAWN {
		return 'p'
	} else if p == B_KNIGHT {
		return 'n'
	} else if p == B_BISHOP {
		return 'b'
	} else if p == B_ROOK {
		return 'r'
	} else if p == B_QUEEN {
		return 'q'
	} else if p == B_KING {
		return 'k'
	}
	return ' '
}

type Square uint8

const (
	SQ_A1 Square = iota
	SQ_B1
	SQ_C1
	SQ_D1
	SQ_E1
	SQ_F1
	SQ_G1
	SQ_H1
	SQ_A2
	SQ_B2
	SQ_C2
	SQ_D2
	SQ_E2
	SQ_F2
	SQ_G2
	SQ_H2
	SQ_A3
	SQ_B3
	SQ_C3
	SQ_D3
	SQ_E3
	SQ_F3
	SQ_G3
	SQ_H3
	SQ_A4
	SQ_B4
	SQ_C4
	SQ_D4
	SQ_E4
	SQ_F4
	SQ_G4
	SQ_H4
	SQ_A5
	SQ_B5
	SQ_C5
	SQ_D5
	SQ_E5
	SQ_F5
	SQ_G5
	SQ_H5
	SQ_A6
	SQ_B6
	SQ_C6
	SQ_D6
	SQ_E6
	SQ_F6
	SQ_G6
	SQ_H6
	SQ_A7
	SQ_B7
	SQ_C7
	SQ_D7
	SQ_E7
	SQ_F7
	SQ_G7
	SQ_H7
	SQ_A8
	SQ_B8
	SQ_C8
	SQ_D8
	SQ_E8
	SQ_F8
	SQ_G8
	SQ_H8
	N_SQUARES
	NULL_SQ Square = math.MaxUint8
)

func SqFromAlg(algCoords string) (Square, error) {
	if len(algCoords) != 2 {
		return 0, fmt.Errorf("unexpected len of algebraic notation for square from %s, expected 2", algCoords)
	}
	if algCoords[0] < 'a' || algCoords[0] > 'h' {
		return 0, fmt.Errorf("file specifier in algebraic notation for square not in expected range ['a', 'h'], got '%c'", algCoords[0])
	}
	file := algCoords[0] - 'a' + 1
	if algCoords[1] < '1' || algCoords[1] > '8' {
		return 0, fmt.Errorf("rank specifier in algebraic notation for square not in expected range ['1', '8'], got '%c'", algCoords[1])
	}
	rank := algCoords[1] - '1' + 1
	return Square(8*(rank-1) + file - 1), nil
}

func SqFromCoords(rank, file int) Square {
	return Square(8*(rank-1) + file - 1)
}

func (s Square) IsNull() bool {
	return s == NULL_SQ
}

func (s Square) String() string {
	if s == NULL_SQ {
		return "0"
	}
	rank := (s / 8) + 1
	file := (s % 8) + 1
	fileChar := byte('0' + file)
	return fmt.Sprintf("%c%d", fileChar, rank)
}

type Result uint8

const (
	RESULT_IN_PROGRESS Result = iota
	RESULT_W_CHECKMATE
	RESULT_B_CHECKMATE
	RESULT_DRAW_MATL
	RESULT_DRAW_50MOVE
	RESULT_DRAW_STALEMATE
)

type CastleRight uint8

const (
	W_CAN_CASTLE_KINGSIDE CastleRight = iota
	W_CAN_CASTLE_QUEENSIDE
	B_CAN_CASTLE_KINGSIDE
	B_CAN_CASTLE_QUEENSIDE
	N_CASTLE_RIGHTS
)

type Color uint8

const (
	WHITE Color = iota
	BLACK
	N_COLORS
)

type ZHash uint64
