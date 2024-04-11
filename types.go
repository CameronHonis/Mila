package main

import "math"

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

type Square uint8

const (
	A1 Square = iota
	B1
	C1
	D1
	E1
	F1
	G1
	H1
	A2
	B2
	C2
	D2
	E2
	F2
	G2
	H2
	A3
	B3
	C3
	D3
	E3
	F3
	G3
	H3
	A4
	B4
	C4
	D4
	E4
	F4
	G4
	H4
	A5
	B5
	C5
	D5
	E5
	F5
	G5
	H5
	A6
	B6
	C6
	D6
	E6
	F6
	G6
	H6
	A7
	B7
	C7
	D7
	E7
	F7
	G7
	H7
	A8
	B8
	C8
	D8
	E8
	F8
	G8
	H8
	N_SQUARES
	NULL_SQ = math.MaxUint8
)

func (s Square) IsNull() bool {
	return s == NULL_SQ
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
