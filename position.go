package main

import (
	"fmt"
	"log"
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
					pos.pieceBitboards[EMPTY] |= BBWithSquares(sq)
					file++
				}
				continue
			}

			sq := SqFromCoords(rank, file)
			piece := PieceFromChar(fenPiece)
			pos.pieces[sq] = piece
			pos.pieceBitboards[piece] |= BBWithSquares(sq)
			if piece.IsWhite() {
				pos.colorBitboards[WHITE] |= BBWithSquares(sq)
			} else {
				pos.colorBitboards[BLACK] |= BBWithSquares(sq)
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
	for colorIdx, colorBB := range p.colorBitboards {
		if colorIdx == 0 {
			continue
		}
		rtn ^= colorBB
	}
	return rtn
}

// IsLegalMove is intended to filter out only valid pseudo-legal moves.
func (p *Position) IsLegalMove(pMove Move) bool {
	piece := p.pieces[pMove.StartSq()]
	pt := piece.Type()
	isWhite := piece.IsWhite()

	var selfColor Color
	var oppColor Color
	if isWhite {
		selfColor = WHITE
		oppColor = BLACK
	} else {
		selfColor = WHITE
		oppColor = WHITE
	}

	if pt == KING {
		if pMove.Type() == CASTLING {
			start := pMove.StartSq()
			end := pMove.EndSq()
			sq0 := start
			var sq1 Square
			if end > start {
				sq1 = sq0 + 1
			} else {
				sq1 = sq0 - 1
			}
			sq2 := end
			return p.isSquareAttacked(oppColor, sq0) ||
				p.isSquareAttacked(oppColor, sq1) ||
				p.isSquareAttacked(oppColor, sq2)
		} else {
			return p.isSquareAttacked(oppColor, pMove.EndSq())
		}
	} else {
		kingSq := p.pieceBitboards[NewPiece(KING, selfColor)].FirstSq()
		if p.isSquareAttacked(oppColor, kingSq) {
			return true
		}

		capturedPiece := p.makeMove(pMove)
		defer p.unmakeMove(pMove, capturedPiece)

		kingSq = p.pieceBitboards[NewPiece(KING, selfColor)].FirstSq()
		return p.isSquareAttacked(oppColor, kingSq)
	}
}

// MakeMove is a cheap way to execute a move on the piece arrangement only.
// To make a move during search, State.MakeMove should instead be used.
func (p *Position) makeMove(move Move) (capturedPiece Piece) {
	start := move.StartSq()
	end := move.EndSq()
	startMask := BBWithSquares(start)
	endMask := BBWithSquares(end)

	piece := p.pieces[start]
	capturedPiece = p.pieces[end]

	p.pieces[start] = EMPTY
	p.pieces[end] = piece

	p.pieceBitboards[piece] ^= startMask | endMask
	p.pieceBitboards[capturedPiece] ^= endMask
	p.pieceBitboards[EMPTY] ^= startMask

	p.colorBitboards[NewColor(p.isWhiteTurn)] ^= startMask | endMask
	if capturedPiece != EMPTY {
		p.colorBitboards[NewColor(!p.isWhiteTurn)] ^= endMask
	}

	if move.Type() == CASTLING {
		endFile := end.File()
		var rookStartSq Square
		var rookEndSq Square
		if endFile == 7 {
			rookStartSq = end + 1
			rookEndSq = end - 1
		} else {
			rookStartSq = end - 2
			rookEndSq = end + 1
		}
		rookStartMask := BBWithSquares(rookStartSq)
		rookEndMask := BBWithSquares(rookEndSq)

		rookPiece := p.pieces[rookStartSq]
		p.pieces[rookStartSq] = EMPTY
		p.pieces[rookEndSq] = rookPiece
		p.pieceBitboards[EMPTY] ^= rookStartMask | rookEndMask
		p.pieceBitboards[rookPiece] ^= rookStartMask | rookEndMask
		p.colorBitboards[NewColor(p.isWhiteTurn)] ^= rookStartMask | rookEndMask
	} else if move.Type() == CAPTURES_EN_PASSANT {
		epSq := SqFromCoords(int(start.Rank()), int(end.File()))
		capturedPiece = NewPiece(PAWN, NewColor(!p.isWhiteTurn))
		captureMask := BBWithSquares(epSq)
		p.pieces[epSq] = EMPTY
		p.pieceBitboards[capturedPiece] ^= captureMask
		p.pieceBitboards[EMPTY] ^= captureMask
		p.colorBitboards[NewColor(!p.isWhiteTurn)] ^= captureMask
	}

	p.isWhiteTurn = !p.isWhiteTurn
	return
}

func (p *Position) unmakeMove(move Move, capturedPiece Piece) {
	start := move.StartSq()
	end := move.EndSq()
	startMask := BBWithSquares(start)
	endMask := BBWithSquares(end)

	piece := p.pieces[end]

	p.pieces[start] = piece
	p.pieces[end] = capturedPiece

	p.pieceBitboards[piece] ^= startMask | endMask
	p.pieceBitboards[EMPTY] ^= startMask
	p.colorBitboards[NewColor(!p.isWhiteTurn)] ^= startMask | endMask

	if move.Type() != CAPTURES_EN_PASSANT {
		p.pieceBitboards[capturedPiece] ^= endMask
		if capturedPiece != EMPTY {
			p.colorBitboards[NewColor(p.isWhiteTurn)] ^= endMask
		}
	}

	if move.Type() == CASTLING {
		endFile := end.File()
		var rookStartSq Square
		var rookEndSq Square
		if endFile == 7 {
			rookStartSq = end + 1
			rookEndSq = end - 1
		} else {
			rookStartSq = end - 2
			rookEndSq = end + 1
		}
		rookStartMask := BBWithSquares(rookStartSq)
		rookEndMask := BBWithSquares(rookEndSq)

		rookPiece := p.pieces[rookEndSq]
		p.pieces[rookEndSq] = EMPTY
		p.pieces[rookStartSq] = rookPiece
		p.pieceBitboards[EMPTY] ^= rookStartMask | rookEndMask
		p.pieceBitboards[rookPiece] ^= rookStartMask | rookEndMask
		p.colorBitboards[NewColor(!p.isWhiteTurn)] ^= rookStartMask | rookEndMask
	} else if move.Type() == CAPTURES_EN_PASSANT {
		if DEBUG {
			if capturedPiece.Type() != PAWN {
				log.Fatalf("unexpected capture piece %s unmaking en passant move in position", capturedPiece)
			}
		}
		epSq := SqFromCoords(int(start.Rank()), int(end.File()))
		captureMask := BBWithSquares(epSq)
		p.pieces[epSq] = capturedPiece
		p.pieceBitboards[capturedPiece] ^= captureMask
		p.pieceBitboards[EMPTY] ^= captureMask
		p.colorBitboards[NewColor(p.isWhiteTurn)] ^= captureMask
	}

	p.isWhiteTurn = !p.isWhiteTurn
}

func (p *Position) isSquareAttacked(attackColor Color, sq Square) bool {
	knightAttacksBB := KnightAttacksBB(sq)
	if attackColor == WHITE {
		if knightAttacksBB&p.pieceBitboards[W_KNIGHT] != 0 {
			return true
		}
	} else {
		if knightAttacksBB&p.pieceBitboards[B_KNIGHT] != 0 {
			return true
		}
	}
	occupiedBB := p.OccupiedBB()
	diagAttacksBB := SlidingAttacksBB(occupiedBB, sq, BISHOP)
	if attackColor == WHITE {
		if diagAttacksBB&p.pieceBitboards[W_BISHOP] != 0 {
			return true
		}
		if diagAttacksBB&p.pieceBitboards[W_QUEEN] != 0 {
			return true
		}
	} else {
		if diagAttacksBB&p.pieceBitboards[B_BISHOP] != 0 {
			return true
		}
		if diagAttacksBB&p.pieceBitboards[B_QUEEN] != 0 {
			return true
		}
	}
	straightAttacks := SlidingAttacksBB(occupiedBB, sq, ROOK)
	if attackColor == WHITE {
		if straightAttacks&p.pieceBitboards[W_ROOK] != 0 {
			return true
		}
		if straightAttacks&p.pieceBitboards[W_QUEEN] != 0 {
			return true
		}
	} else {
		if straightAttacks&p.pieceBitboards[B_ROOK] != 0 {
			return true
		}
		if straightAttacks&p.pieceBitboards[B_QUEEN] != 0 {
			return true
		}
	}

	file := sq.File()
	if file > 1 {
		if attackColor == WHITE {
			if p.pieces[sq-9] == W_PAWN {
				return true
			}
		} else {
			if p.pieces[sq+7] == B_PAWN {
				return true
			}
		}
	}
	if file < 8 {
		if attackColor == WHITE {
			if p.pieces[sq-7] == W_PAWN {
				return true
			}
		} else {
			if p.pieces[sq+9] == B_PAWN {
				return true
			}
		}
	}
	return false
}
