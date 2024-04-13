package main

import "log"

// Move is a compressed bit representation of a move
// The compress is important for minimizing memory requirements in the hash table
// The bit layout of a move is as follows:
// - first 6 bits represent start square
// - next 6 bits represent end square
// - last 4 bits represent the captured piece
type Move uint16

func NewMove(startSq, endSq Square, state *State) Move {
	if SAFE_MOVE_PARSING {
		if startSq > 0b111111 {
			log.Fatalf("start square occupies more than 6 bits: %b", startSq)
		}
		if endSq > 0b111111 {
			log.Fatalf("end square occupies more than 6 bits: %b", endSq)
		}
	}
	capturedPiece := state.Pos.pieces[endSq]
	if SAFE_MOVE_PARSING {
		if capturedPiece > 0b1111 {
			log.Fatalf("captured piece bits out of range: %b", capturedPiece)
		}
	}
	return Move(startSq&0b111111)<<10 | Move(endSq&0b111111)<<4 | Move(capturedPiece&0b1111)
}

func NullMove() Move {
	return 0
}

func (m Move) StartSq() Square {
	return Square(m >> 10)
}

func (m Move) EndSq() Square {
	return Square((m >> 4) & 0b111111)
}

func (m Move) CapturedPiece() Piece {
	return Piece(m & 0b1111)
}

func (m Move) IsNull() bool {
	return m == 0
}
