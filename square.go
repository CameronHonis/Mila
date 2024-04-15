package main

import (
	"fmt"
	"math"
)

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

func (s Square) Rank() uint8 {
	return uint8(s)/8 + 1
}

func (s Square) File() uint8 {
	return uint8(s)%8 + 1
}

func (s Square) String() string {
	if s == NULL_SQ {
		return "0"
	}
	rank := (s / 8) + 1
	file := (s % 8) + 1
	fileChar := byte('a' + file - 1)
	return fmt.Sprintf("%c%d", fileChar, rank)
}

func (s Square) PosDiagIdx() uint8 {
	if s.File() >= s.Rank() {
		return 7 - uint8(s%9)
	} else {
		return 16 - uint8(s%9)
	}
}

func (s Square) NegDiagIdx() uint8 {
	if s.File()+s.Rank() < 9 {
		return uint8(s % 7)
	} else {
		if s == SQ_H8 {
			return 14
		} else {
			return 7 + uint8(s%7)
		}
	}
}
