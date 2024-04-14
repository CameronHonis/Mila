package main

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
