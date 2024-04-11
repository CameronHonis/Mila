package main

import (
	"fmt"
	"strings"
)

type Position struct {
	pieces         [N_SQUARES]Piece
	pieceBitboards [N_PIECES]Bitboard
	castleRights   [N_CASTLE_RIGHTS]bool
	enPassantSq    Square
	result         Result
}

func InitPos() *Position {
	return &Position{
		pieces: [N_SQUARES]Piece{
			W_ROOK, W_KNIGHT, W_BISHOP, W_QUEEN, W_KING, W_BISHOP, W_KNIGHT, W_ROOK,
			W_PAWN, W_PAWN, W_PAWN, W_PAWN, W_PAWN, W_PAWN, W_PAWN, W_PAWN,
			EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY,
			EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY,
			EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY,
			EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY,
			B_PAWN, B_PAWN, B_PAWN, B_PAWN, B_PAWN, B_PAWN, B_PAWN, B_PAWN,
			B_ROOK, B_KNIGHT, B_BISHOP, B_QUEEN, B_KING, B_BISHOP, B_KNIGHT, B_ROOK,
		},
		pieceBitboards: [N_PIECES]Bitboard{
			0b00000000_00000000_11111111_11111111_11111111_11111111_00000000_00000000,
			0b00000000_11111111_00000000_00000000_00000000_00000000_00000000_00000000,
			0b01000010_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
			0b00100100_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
			0b10000001_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
			0b00010000_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
			0b00001000_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
			0b00000000_00000000_00000000_00000000_00000000_00000000_11111111_00000000,
			0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_01000010,
			0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00100100,
			0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_10000001,
			0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00010000,
			0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00001000,
		},
		castleRights: [N_CASTLE_RIGHTS]bool{true, true, true, true},
		enPassantSq:  NULL_SQ,
		result:       RESULT_IN_PROGRESS,
	}
}

func (p *Position) String() string {
	var rtnBuilder = strings.Builder{}
	for rank := 8; rank > 0; rank-- {
		rtnBuilder.WriteString(fmt.Sprintf("%d ", rank))
		for file := 1; file < 9; file++ {
			idx := 8*(rank-1) + (file - 1)
			piece := p.pieces[idx]
			pieceStr := piece.String()
			if piece > 0 {
				rtnBuilder.WriteString(fmt.Sprintf("%s ", pieceStr))
			} else {
				var isDark bool
				if rank%2 == 0 {
					isDark = file%2 == 0
				} else {
					isDark = file%2 == 1
				}
				if isDark {
					rtnBuilder.WriteString("░░")
				} else {
					rtnBuilder.WriteString("  ")
				}
			}
		}
		rtnBuilder.WriteByte('\n')
	}
	rtnBuilder.WriteString("  ")
	for file := 1; file < 9; file++ {
		rtnBuilder.WriteByte(byte('0' + file))
		rtnBuilder.WriteByte(' ')
	}
	return rtnBuilder.String()
}
