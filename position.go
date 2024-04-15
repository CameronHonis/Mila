package main

import (
	"fmt"
	"strings"
)

type Position struct {
	pieces         [N_SQUARES]Piece
	pieceBitboards [N_PIECES]Bitboard
	colorBitboards [N_COLORS]Bitboard
	isWhiteTurn    bool
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
			0b00000000_00000000_00000000_00000000_00000000_00000000_11111111_00000000,
			0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_01000010,
			0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00100100,
			0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_10000001,
			0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00001000,
			0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00010000,
			0b00000000_11111111_00000000_00000000_00000000_00000000_00000000_00000000,
			0b01000010_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
			0b00100100_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
			0b10000001_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
			0b00001000_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
			0b00010000_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
		},
		colorBitboards: [N_COLORS]Bitboard{
			0b00000000_00000000_00000000_00000000_00000000_00000000_11111111_11111111,
			0b11111111_11111111_00000000_00000000_00000000_00000000_00000000_00000000,
		},
		isWhiteTurn: true,
	}
}

func PosFromFEN(fen string) (*Position, error) {
	var pos = &Position{}
	fenSegs := strings.Split(fen, " ")
	if len(fenSegs) != 6 {
		return nil, fmt.Errorf("invalid number of fen segments %d, expected 6", len(fenSegs))
	}

	piecesFen := fenSegs[0]
	fenPiecesRows := strings.Split(piecesFen, "/")
	if len(fenPiecesRows) != 8 {
		return nil, fmt.Errorf("invalid number of rows in FEN pieces %d, expected 8", len(fenPiecesRows))
	}
	for fenRowIdx, fenPiecesRow := range fenPiecesRows {
		rank := 8 - fenRowIdx
		var file = 1
		for _, fenPiece := range []byte(fenPiecesRow) {
			if file > 8 {
				return nil, fmt.Errorf("too many pieces on rank %d in fen %s", rank, fen)
			}
			if fenPiece >= '1' && fenPiece <= '8' {
				for i := 0; i < int(fenPiece-'0'); i++ {
					sq := SqFromCoords(rank, file)
					pos.pieceBitboards[EMPTY] |= BBWithHighBitsAt(int(sq))
					file++
				}
				continue
			}

			sq := SqFromCoords(rank, file)
			piece := PieceFromChar(fenPiece)
			pos.pieces[sq] = piece
			pos.pieceBitboards[piece] |= BBWithHighBitsAt(int(sq))
			if piece.IsWhite() {
				pos.colorBitboards[WHITE] |= BBWithHighBitsAt(int(sq))
			} else {
				pos.colorBitboards[BLACK] |= BBWithHighBitsAt(int(sq))
			}
			file++
		}
		if file < 9 {
			return nil, fmt.Errorf("not enough pieces on rank %d in fen %s", rank, fen)
		}
	}

	turnSpecifier := fenSegs[1]
	if turnSpecifier == "w" {
		pos.isWhiteTurn = true
	} else if turnSpecifier == "b" {
		pos.isWhiteTurn = false
	} else {
		return nil, fmt.Errorf("unexpected length for turn specifier %s in fen %s", fenSegs[1], fen)
	}

	return pos, nil
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

func (p *Position) FEN() string {
	var rtnBuilder strings.Builder
	for rank := 8; rank > 0; rank-- {
		var consecSpaces int
		for file := 1; file < 9; file++ {
			sq := SqFromCoords(rank, file)
			if sq == 64 {
				fmt.Println(rank, file)
			}
			piece := p.pieces[sq]
			if piece == EMPTY {
				consecSpaces++
			} else {
				if consecSpaces > 0 {
					rtnBuilder.WriteByte(byte('0' + consecSpaces))
				}
				rtnBuilder.WriteByte(piece.Char())
			}
		}
		if consecSpaces > 0 {
			rtnBuilder.WriteByte(byte('0' + consecSpaces))
		}
		if rank != 1 {
			rtnBuilder.WriteByte('/')
		}
	}
	rtnBuilder.WriteByte(' ')
	if p.isWhiteTurn {
		rtnBuilder.WriteByte('w')
	} else {
		rtnBuilder.WriteByte('b')
	}
	return rtnBuilder.String()
}

func (p *Position) OccupiedBB() Bitboard {
	// TODO: evaluate the speed of this method in the broader context of search speed
	//		 relative to maintaining an updated Position.occupied field on make/unmake
	//		 move.

	var rtn Bitboard
	for pieceIdx, pieceBB := range p.pieceBitboards {
		if pieceIdx == 0 {
			continue
		}
		rtn ^= pieceBB
	}
	return rtn
}
