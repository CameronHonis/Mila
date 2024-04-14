package main

import (
	"fmt"
	"log"
	"math/bits"
	"strings"
)

type Bitboard uint64

const (
	RANK_1      = Bitboard(0b11111111)
	RANK_2      = Bitboard(0b11111111 << 8)
	RANK_3      = Bitboard(0b11111111 << 16)
	RANK_4      = Bitboard(0b11111111 << 24)
	RANK_5      = Bitboard(0b11111111 << 32)
	RANK_6      = Bitboard(0b11111111 << 40)
	RANK_7      = Bitboard(0b11111111 << 48)
	RANK_8      = Bitboard(0b11111111 << 56)
	FILE_1      = Bitboard(0b00000001_00000001_00000001_00000001_00000001_00000001_00000001_00000001)
	FILE_2      = Bitboard(0b00000010_00000010_00000010_00000010_00000010_00000010_00000010_00000010)
	FILE_3      = Bitboard(0b00000100_00000100_00000100_00000100_00000100_00000100_00000100_00000100)
	FILE_4      = Bitboard(0b00001000_00001000_00001000_00001000_00001000_00001000_00001000_00001000)
	FILE_5      = Bitboard(0b00010000_00010000_00010000_00010000_00010000_00010000_00010000_00010000)
	FILE_6      = Bitboard(0b00100000_00100000_00100000_00100000_00100000_00100000_00100000_00100000)
	FILE_7      = Bitboard(0b01000000_01000000_01000000_01000000_01000000_01000000_01000000_01000000)
	FILE_8      = Bitboard(0b10000000_10000000_10000000_10000000_10000000_10000000_10000000_10000000)
	POS_DIAG_0  = Bitboard(0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_10000000)
	POS_DIAG_1  = Bitboard(0b00000000_00000000_00000000_00000000_00000000_00000000_10000000_01000000)
	POS_DIAG_2  = Bitboard(0b00000000_00000000_00000000_00000000_00000000_10000000_01000000_00100000)
	POS_DIAG_3  = Bitboard(0b00000000_00000000_00000000_00000000_10000000_01000000_00100000_00010000)
	POS_DIAG_4  = Bitboard(0b00000000_00000000_00000000_10000000_01000000_00100000_00010000_00001000)
	POS_DIAG_5  = Bitboard(0b00000000_00000000_10000000_01000000_00100000_00010000_00001000_00000100)
	POS_DIAG_6  = Bitboard(0b00000000_10000000_01000000_00100000_00010000_00001000_00000100_00000010)
	POS_DIAG_7  = Bitboard(0b10000000_01000000_00100000_00010000_00001000_00000100_00000010_00000001)
	POS_DIAG_8  = Bitboard(0b01000000_00100000_00010000_00001000_00000100_00000010_00000001_00000000)
	POS_DIAG_9  = Bitboard(0b00100000_00010000_00001000_00000100_00000010_00000001_00000000_00000000)
	POS_DIAG_10 = Bitboard(0b00010000_00001000_00000100_00000010_00000001_00000000_00000000_00000000)
	POS_DIAG_11 = Bitboard(0b00001000_00000100_00000010_00000001_00000000_00000000_00000000_00000000)
	POS_DIAG_12 = Bitboard(0b00000100_00000010_00000001_00000000_00000000_00000000_00000000_00000000)
	POS_DIAG_13 = Bitboard(0b00000010_00000001_00000000_00000000_00000000_00000000_00000000_00000000)
	POS_DIAG_14 = Bitboard(0b00000001_00000000_00000000_00000000_00000000_00000000_00000000_00000000)
	NEG_DIAG_0  = Bitboard(0b1)
	NEG_DIAG_1  = Bitboard(0b100000010)
	NEG_DIAG_2  = Bitboard(0b10000001000000100)
	NEG_DIAG_3  = Bitboard(0b1000000100000010000001000)
	NEG_DIAG_4  = Bitboard(0b100000010000001000000100000010000)
	NEG_DIAG_5  = Bitboard(0b10000001000000100000010000001000000100000)
	NEG_DIAG_6  = Bitboard(0b1000000100000010000001000000100000010000001000000)
	NEG_DIAG_7  = Bitboard(0b100000010000001000000100000010000001000000100000010000000)
	NEG_DIAG_8  = Bitboard(0b1000000100000010000001000000100000010000001000000000000000)
	NEG_DIAG_9  = Bitboard(0b10000001000000100000010000001000000100000000000000000000000)
	NEG_DIAG_10 = Bitboard(0b100000010000001000000100000010000000000000000000000000000000)
	NEG_DIAG_11 = Bitboard(0b1000000100000010000001000000000000000000000000000000000000000)
	NEG_DIAG_12 = Bitboard(0b10000001000000100000000000000000000000000000000000000000000000)
	NEG_DIAG_13 = Bitboard(0b100000010000000000000000000000000000000000000000000000000000000)
	NEG_DIAG_14 = Bitboard(0b1000000000000000000000000000000000000000000000000000000000000000)
)

const N_RANKS = 8
const N_FILES = 8
const N_DIAGS = 15

func BBWithHighBitsAt(idxs ...int) Bitboard {
	var rtn Bitboard
	for _, idx := range idxs {
		rtn |= 1 << idx
	}
	return rtn
}

func BBWithSquares(squares ...Square) Bitboard {
	var rtn Bitboard
	for _, sq := range squares {
		rtn |= 1 << sq
	}
	return rtn
}

func BBWithRank(rank, bits uint8) Bitboard {
	if DEBUG {
		if rank < 1 || rank > N_RANKS {
			log.Fatalf("invalid rank %d, expected within range [1, %d]", rank, N_RANKS)
		}
	}
	return Bitboard(bits) << (8 * (rank - 1))
}

func BBWithFile(file, bits uint8) Bitboard {
	if DEBUG {
		if file < 1 || file > N_FILES {
			log.Fatalf("invalid file %d, expected within range [1, %d]", file, N_FILES)
		}
	}
	var bb Bitboard
	for r := uint8(0); r < 8; r++ {
		rBit := (bits >> r) & 0b1
		bbIdx := 8*r + file - 1
		bb |= Bitboard(rBit) << bbIdx
	}
	return bb
}

// BBWithPosDiag writes all the bits in order on the positive diagonal
// The diagonal index starts from H1 -> A1 -> A8. The less significant the bit,
// the closer to the bottom (rank=1).
func BBWithPosDiag(diagIdx, bits uint8) Bitboard {
	if DEBUG {
		if diagIdx < 0 || diagIdx >= N_DIAGS {
			log.Fatalf("invalid diag idx %d, expected within range [0, %d]", diagIdx, N_DIAGS)
		}
	}
	var bb Bitboard
	var bbIdx uint8
	var nDiagSq uint8
	if diagIdx <= 7 {
		bbIdx = 7 - diagIdx
		nDiagSq = diagIdx + 1
	} else {
		bbIdx = 8 * (diagIdx - 7)
		nDiagSq = 15 - diagIdx
	}

	for i := uint8(0); i < nDiagSq; i++ {
		iBit := (bits >> i) & 0b1
		bb |= Bitboard(iBit) << bbIdx
		bbIdx += 9
	}

	return bb
}

// BBWithNegDiag writes all bits in order on the negative diagonal
// The diagonal index starts from A1 -> H1 -> H8. The less significant the bit,
// the closer to the bottom (rank=1).
func BBWithNegDiag(diagIdx, bits uint8) Bitboard {
	if DEBUG {
		if diagIdx < 0 || diagIdx >= N_DIAGS {
			log.Fatalf("invalid diag idx %d, expected within range [0, %d]", diagIdx, N_DIAGS)
		}
	}
	var bb Bitboard
	var bbIdx uint8
	var nDiagSq uint8
	if diagIdx <= 7 {
		bbIdx = diagIdx
		nDiagSq = diagIdx + 1
	} else {
		bbIdx = 8*(diagIdx-7) + 7
		nDiagSq = 15 - diagIdx
	}

	for i := uint8(0); i < nDiagSq; i++ {
		iBit := (bits >> i) & 0b1
		bb |= Bitboard(iBit) << bbIdx
		bbIdx += 7
	}

	return bb
}

func (b Bitboard) String() string {
	var rtnBuilder = strings.Builder{}
	for rank := 8; rank > 0; rank-- {
		rtnBuilder.WriteString(fmt.Sprintf("%d ", rank))
		var mask Bitboard
		for file := 1; file < 9; file++ {
			idx := 8*(rank-1) + (file - 1)
			mask = 1 << idx
			if mask&b > 0 {
				rtnBuilder.WriteString("██")
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

// FirstSq returns the square of the least significant bit in the bitboard
// This essentially iterations in row-major order starting from
// A1 -> ... -> H1 -> ... -> H8
func (b Bitboard) FirstSq() Square {
	lsb := bits.TrailingZeros64(uint64(b))
	if lsb == 64 {
		return NULL_SQ
	} else {
		return Square(lsb)
	}
}
