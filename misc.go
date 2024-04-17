package main

import "log"

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

func NewPiece(pt PieceType, color Color) Piece {
	if pt == EMPTY_PIECE_TYPE {
		return EMPTY
	}
	if color == WHITE {
		if pt == PAWN {
			return W_PAWN
		} else if pt == KNIGHT {
			return W_KNIGHT
		} else if pt == BISHOP {
			return W_BISHOP
		} else if pt == ROOK {
			return W_ROOK
		} else if pt == QUEEN {
			return W_QUEEN
		} else {
			return W_KING
		}
	} else {
		if pt == PAWN {
			return B_PAWN
		} else if pt == KNIGHT {
			return B_KNIGHT
		} else if pt == BISHOP {
			return B_BISHOP
		} else if pt == ROOK {
			return B_ROOK
		} else if pt == QUEEN {
			return B_QUEEN
		} else {
			return B_KING
		}
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

func (p Piece) Type() PieceType {
	if DEBUG {
		if p >= N_PIECES {
			log.Fatalf("cannot get piece type of unknown piece %d", p)
		}
	}
	if p == EMPTY {
		return EMPTY_PIECE_TYPE
	} else if p == W_PAWN || p == B_PAWN {
		return PAWN
	} else if p == W_KNIGHT || p == B_KNIGHT {
		return KNIGHT
	} else if p == W_BISHOP || p == B_BISHOP {
		return BISHOP
	} else if p == W_ROOK || p == B_ROOK {
		return ROOK
	} else if p == W_QUEEN || p == B_QUEEN {
		return QUEEN
	} else {
		return KING
	}
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
	W_CASTLE_KINGSIDE_RIGHT CastleRight = iota
	W_CASTLE_QUEENSIDE_RIGHT
	B_CASTLE_KINGSIDE_RIGHT
	B_CASTLE_QUEENSIDE_RIGHT
	N_CASTLE_RIGHTS
)

type Color uint8

func NewColor(isWhite bool) Color {
	if isWhite {
		return WHITE
	} else {
		return BLACK
	}
}

const (
	WHITE Color = iota
	BLACK
	N_COLORS
)

type ZHash uint64
