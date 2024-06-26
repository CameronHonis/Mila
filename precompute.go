package main

// N_ROW_PERMUS is the number of possible permutations of occupied squares
// for any row (2^8). This number happens to be the permutations of occupied
// squares for columns or diagonals as well.
const N_ROW_PERMUS = 1 << 8

// Splitting the precomputed rook moves into rows and columns as separate look tables
// reduces the overall amount of entries in the set. This works out since a rooks moves
// on a given row have no dependence to its moves on a file, and vice versa.

var rankAttacks map[Bitboard][N_FILES]Bitboard
var fileAttacks map[Bitboard][N_RANKS]Bitboard
var posDiagAttacks map[Bitboard][N_RANKS]Bitboard
var negDiagAttacks map[Bitboard][N_RANKS]Bitboard
var pawnAttacks [N_SQUARES][N_COLORS]Bitboard
var knightAttacks [N_SQUARES]Bitboard
var kingAttacks [N_SQUARES]Bitboard

var isInitted = false

func initAttackPrecomputes() {
	if isInitted {
		return
	}
	initRankAttacks()
	initFileAttacks()
	initPosDiagAttacks()
	initNegDiagAttacks()
	initPawnAttacks()
	initKnightAttacks()
	initKingAttacks()
	isInitted = true
}

func initRankAttacks() {
	rankAttacks = make(map[Bitboard][N_FILES]Bitboard)
	for _, occupiedBB := range genRankOccupiedBBs() {
		rank := occupiedBB.FirstSq().Rank()
		attackBBs := [N_FILES]Bitboard{}
		for file := uint8(1); file <= N_FILES; file++ {
			var attackBB Bitboard
			sq := SqFromCoords(int(rank), int(file))
			leftProbe := sq
			for leftProbe.File() > 1 {
				leftProbe--
				mask := BBWithHighBitsAt(int(leftProbe))
				attackBB |= mask
				if occupiedBB&mask > 0 {
					break
				}
			}
			rightProbe := sq
			for rightProbe.File() < 8 {
				rightProbe++
				mask := BBWithHighBitsAt(int(rightProbe))
				attackBB |= mask
				if occupiedBB&mask > 0 {
					break
				}
			}
			attackBBs[file-1] = attackBB
		}
		rankAttacks[occupiedBB] = attackBBs
	}
}

func initFileAttacks() {
	fileAttacks = make(map[Bitboard][N_RANKS]Bitboard)
	for _, occupiedBB := range genFileOccupiedBBs() {
		file := occupiedBB.FirstSq().File()
		attackBBs := [N_RANKS]Bitboard{}
		for rank := uint8(1); rank <= N_RANKS; rank++ {
			var attackBB Bitboard
			sq := SqFromCoords(int(rank), int(file))
			upProbe := sq
			for upProbe.Rank() < 8 {
				upProbe += 8
				mask := BBWithHighBitsAt(int(upProbe))
				attackBB |= mask
				if occupiedBB&mask > 0 {
					break
				}
			}
			downProbe := sq
			for downProbe.Rank() > 1 {
				downProbe -= 8
				mask := BBWithHighBitsAt(int(downProbe))
				attackBB |= mask
				if occupiedBB&mask > 0 {
					break
				}
			}
			attackBBs[rank-1] = attackBB
		}
		fileAttacks[occupiedBB] = attackBBs
	}
}

func initPosDiagAttacks() {
	posDiagAttacks = make(map[Bitboard][N_RANKS]Bitboard)
	for _, occupiedBB := range genPosDiagOccupiedBBs() {
		diagIdx := occupiedBB.FirstSq().PosDiagIdx()
		diagMask := BBWithPosDiag(diagIdx, 0b11111111)
		attackBBs := [N_RANKS]Bitboard{}
		for rank := uint8(1); rank <= N_RANKS; rank++ {
			var attackBB Bitboard
			rankMask := BBWithRank(rank, 0b11111111)
			sq := (diagMask & rankMask).FirstSq()
			if sq == NULL_SQ {
				attackBBs[rank-1] = attackBB
				continue
			}
			downLeftProbe := sq
			for downLeftProbe.Rank() > 1 && downLeftProbe.File() > 1 {
				downLeftProbe -= 9
				mask := BBWithHighBitsAt(int(downLeftProbe))
				attackBB |= mask
				if occupiedBB&mask > 0 {
					break
				}
			}
			upRightProbe := sq
			for upRightProbe.Rank() < 8 && upRightProbe.File() < 8 {
				upRightProbe += 9
				mask := BBWithHighBitsAt(int(upRightProbe))
				attackBB |= mask
				if occupiedBB&mask > 0 {
					break
				}
			}
			attackBBs[rank-1] = attackBB
		}
		posDiagAttacks[occupiedBB] = attackBBs
	}
}

func initNegDiagAttacks() {
	negDiagAttacks = make(map[Bitboard][N_RANKS]Bitboard)
	for _, occupiedBB := range genNegDiagOccupiedBBs() {
		diagIdx := occupiedBB.FirstSq().NegDiagIdx()
		diagMask := BBWithNegDiag(diagIdx, 0b11111111)
		attackBBs := [N_RANKS]Bitboard{}
		for rank := uint8(1); rank <= N_RANKS; rank++ {
			var attackBB Bitboard
			rankMask := BBWithRank(rank, 0b11111111)
			sq := (diagMask & rankMask).FirstSq()
			if sq == NULL_SQ {
				attackBBs[rank-1] = attackBB
				continue
			}
			downRightProbe := sq
			for downRightProbe.Rank() > 1 && downRightProbe.File() < 8 {
				downRightProbe -= 7
				mask := BBWithHighBitsAt(int(downRightProbe))
				attackBB |= mask
				if occupiedBB&mask > 0 {
					break
				}
			}
			upLeftProbe := sq
			for upLeftProbe.Rank() < 8 && upLeftProbe.File() > 1 {
				upLeftProbe += 7
				mask := BBWithHighBitsAt(int(upLeftProbe))
				attackBB |= mask
				if occupiedBB&mask > 0 {
					break
				}
			}
			attackBBs[rank-1] = attackBB
		}
		negDiagAttacks[occupiedBB] = attackBBs
	}
}

func initPawnAttacks() {
	pawnAttacks = [N_SQUARES][N_COLORS]Bitboard{}
	for sq := SQ_A1; sq < SQ_A8; sq++ {
		var bb Bitboard
		file := sq.File()
		if file > 1 {
			bb |= BBWithSquares(sq + 7)
		}
		if file < 8 {
			bb |= BBWithSquares(sq + 9)
		}
		pawnAttacks[sq][WHITE] = bb
	}
	for sq := SQ_A2; sq <= SQ_H8; sq++ {
		var bb Bitboard
		file := sq.File()
		if file > 1 {
			bb |= BBWithSquares(sq - 9)
		}
		if file < 8 {
			bb |= BBWithSquares(sq - 7)
		}
		pawnAttacks[sq][BLACK] = bb
	}
}

func initKnightAttacks() {
	knightAttacks = [N_SQUARES]Bitboard{}
	for sq := SQ_A1; sq < N_SQUARES; sq++ {
		file := sq.File()
		rank := sq.Rank()
		var bb Bitboard
		if file > 1 && rank < 7 {
			bb |= BBWithSquares(sq + 15)
		}
		if file > 1 && rank > 2 {
			bb |= BBWithSquares(sq - 17)
		}
		if file > 2 && rank < 8 {
			bb |= BBWithSquares(sq + 6)
		}
		if file > 2 && rank > 1 {
			bb |= BBWithSquares(sq - 10)
		}
		if file < 7 && rank < 8 {
			bb |= BBWithSquares(sq + 10)
		}
		if file < 7 && rank > 1 {
			bb |= BBWithSquares(sq - 6)
		}
		if file < 8 && rank < 7 {
			bb |= BBWithSquares(sq + 17)
		}
		if file < 8 && rank > 2 {
			bb |= BBWithSquares(sq - 15)
		}
		knightAttacks[sq] = bb
	}
}

func initKingAttacks() {
	kingAttacks = [N_SQUARES]Bitboard{}
	for sq := SQ_A1; sq < N_SQUARES; sq++ {
		file := sq.File()
		rank := sq.Rank()
		var bb Bitboard
		if file > 1 {
			bb |= BBWithSquares(sq - 1) // left
			if rank > 1 {
				bb |= BBWithSquares(sq - 9) // left-down
			}
			if rank < 8 {
				bb |= BBWithSquares(sq + 7) // left-up
			}
		}
		if file < 8 {
			bb |= BBWithSquares(sq + 1) // right
			if rank > 1 {
				bb |= BBWithSquares(sq - 7) // right-down
			}
			if rank < 8 {
				bb |= BBWithSquares(sq + 9) // right-up
			}
		}
		if rank > 1 {
			bb |= BBWithSquares(sq - 8) // down
		}
		if rank < 8 {
			bb |= BBWithSquares(sq + 8) // up
		}
		kingAttacks[sq] = bb
	}
}

func genRankOccupiedBBs() []Bitboard {
	rtn := make([]Bitboard, 0)
	for rank := uint8(1); rank <= N_RANKS; rank++ {
		for occupied := 1; occupied < N_ROW_PERMUS; occupied++ {
			bb := BBWithRank(rank, uint8(occupied))
			rtn = append(rtn, bb)
		}
	}
	return rtn
}

func genFileOccupiedBBs() []Bitboard {
	rtn := make([]Bitboard, 0)
	for file := uint8(1); file <= N_FILES; file++ {
		for occupied := 1; occupied < N_ROW_PERMUS; occupied++ {
			bb := BBWithFile(file, uint8(occupied))
			rtn = append(rtn, bb)
		}
	}
	return rtn
}

func genPosDiagOccupiedBBs() []Bitboard {
	rtn := make([]Bitboard, 0)
	for posDiagIdx := uint8(0); posDiagIdx < N_DIAGS; posDiagIdx++ {
		var diagSize uint8
		if posDiagIdx <= 7 {
			diagSize = 8 - (7 - posDiagIdx)
		} else {
			diagSize = 8 - (posDiagIdx - 7)
		}
		nDiagPerms := 1 << diagSize
		for occupied := 1; occupied < nDiagPerms; occupied++ {
			bb := BBWithPosDiag(posDiagIdx, uint8(occupied))
			rtn = append(rtn, bb)
		}
	}
	return rtn
}

func genNegDiagOccupiedBBs() []Bitboard {
	rtn := make([]Bitboard, 0)
	for negDiagIdx := uint8(0); negDiagIdx < N_DIAGS; negDiagIdx++ {
		var diagSize uint8
		if negDiagIdx <= 7 {
			diagSize = 8 - (7 - negDiagIdx)
		} else {
			diagSize = 8 - (negDiagIdx - 7)
		}
		nDiagPerms := 1 << diagSize
		for occupied := 1; occupied < nDiagPerms; occupied++ {
			bb := BBWithNegDiag(negDiagIdx, uint8(occupied))
			rtn = append(rtn, bb)
		}
	}
	return rtn
}
