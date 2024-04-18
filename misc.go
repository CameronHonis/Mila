package main

import (
	"log"
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

func (p Piece) Color() Color {
	if DEBUG {
		if p == 0 {
			log.Fatalf("cannot get color for EMPTY piece")
		}
		if p >= N_PIECES {
			log.Fatalf("cannot get color for piece %d, piece out of range", p)
		}
	}
	if p >= W_PAWN && p <= W_KING {
		return WHITE
	} else {
		return BLACK
	}
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
	RESULT_DRAW_MATL
	RESULT_DRAW_RULE50
	RESULT_DRAW_REPETITION
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

type Ply uint

func PlyFromNMoves(nMoves uint, isWhiteTurn bool) Ply {
	if isWhiteTurn {
		return Ply(nMoves-1) * 2
	} else {
		return Ply(nMoves)*2 - 1
	}
}

func NMovesFromPly(ply Ply) uint {
	return uint(ply+2) / 2
}

// Material represents all the material on a current position. Both the light and
// dark square bishops must be accounted for. Therefore, material piece indices are
// as follows:
//  00. W_PAWN
//  01. W_KNIGHT
//  02. W_LIGHT_BISHOP
//  03. W_DARK_BISHOP
//  04. W_ROOK
//  05. W_QUEEN
//  06. B_PAWN
//  07. B_KNIGHT
//  08. B_LIGHT_BISHOP
//  09. B_DARK_BISHOP
//  10. B_ROOK
//  11. B_QUEEN
type Material [12]uint8

func InitMaterial() Material {
	return Material{8, 2, 1, 1, 2, 1, 8, 2, 1, 1, 2, 1}
}

func (m *Material) AddPiece(piece Piece, sq Square) {
	matIdx := m.pieceToMatIdx(piece, sq)
	if DEBUG {
		if piece.Type() == KING {
			log.Fatalf("cannot add kings to Material")
		}
		if m[matIdx] == math.MaxUint8 {
			log.Fatalf("could not remove piece %s from material, piece count already 255", piece)
		}
	}
	m[matIdx]++
}

func (m *Material) RemovePiece(piece Piece, sq Square) {
	matIdx := m.pieceToMatIdx(piece, sq)
	if DEBUG {
		if piece.Type() == KING {
			log.Fatalf("cannot remove kings from Material")
		}
		if m[matIdx] == 0 {
			log.Fatalf("could not remove piece %s from material, piece count already 0", piece)
		}
	}
	m[matIdx]--
}

func (m *Material) IsForcedDraw() bool {
	if m.nQueens() > 0 {
		return false
	}
	if m.nRooks() > 0 {
		return false
	}
	if m.nPawns() > 0 {
		return false
	}

	nWKnights := m.nWKnights()
	nWLightBishops := m.nWLightBishops()
	nWDarkBishops := m.nWDarkBishops()
	nWBishops := nWLightBishops + nWDarkBishops
	if nWKnights > 2 {
		return false
	}
	if nWKnights > 0 && nWBishops > 0 {
		return false
	}
	if nWLightBishops > 0 && nWDarkBishops > 0 {
		return false
	}

	nBKnights := m.nBKnights()
	nBLightBishops := m.nBLightBishops()
	nBDarkBishops := m.nBDarkBishops()
	nBBishops := nBLightBishops + nBDarkBishops
	if nBKnights > 2 {
		return false
	}
	if nBKnights > 0 && nBBishops > 0 {
		return false
	}
	if nBLightBishops > 0 && nBDarkBishops > 0 {
		return false
	}

	return true
}

func (m *Material) nQueens() uint8 {
	return m[5] + m[11]
}

func (m *Material) nRooks() uint8 {
	return m[4] + m[10]
}

func (m *Material) nPawns() uint8 {
	return m[0] + m[6]
}

func (m *Material) nWKnights() uint8 {
	return m[1]
}

func (m *Material) nWLightBishops() uint8 {
	return m[2]
}

func (m *Material) nWDarkBishops() uint8 {
	return m[3]
}

func (m *Material) nBKnights() uint8 {
	return m[7]
}

func (m *Material) nBLightBishops() uint8 {
	return m[8]
}

func (m *Material) nBDarkBishops() uint8 {
	return m[9]
}

func (m *Material) pieceToMatIdx(piece Piece, sq Square) uint {
	if piece == W_PAWN {
		return 0
	} else if piece == W_KNIGHT {
		return 1
	} else if piece == W_BISHOP {
		if sq.IsDark() {
			return 3
		} else {
			return 2
		}
	} else if piece == W_ROOK {
		return 4
	} else if piece == W_QUEEN {
		return 5
	} else if piece == B_PAWN {
		return 6
	} else if piece == B_KNIGHT {
		return 7
	} else if piece == B_BISHOP {
		if sq.IsDark() {
			return 9
		} else {
			return 8
		}
	} else if piece == B_ROOK {
		return 10
	} else { // B_QUEEN
		return 11
	}
}
