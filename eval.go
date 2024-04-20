package main

const (
	PAWN_VAL   = int16(100)
	KNIGHT_VAL = int16(280)
	BISHOP_VAL = int16(320)
	ROOK_VAL   = int16(500)
	QUEEN_VAL  = int16(950)
	MATE_VAL   = int16(10000)
	DRAW_VAL   = int16(-50)
)

func EvalPos(pos *Position) int16 {
	if pos.result != RESULT_IN_PROGRESS {
		return DRAW_VAL
	}
	mat := pos.material
	eval := PAWN_VAL*mat.pawnDiff() + KNIGHT_VAL*mat.knightDiff() + BISHOP_VAL*mat.bishopDiff() +
		ROOK_VAL*mat.rookDiff() + QUEEN_VAL*mat.queenDiff()
	if pos.isWhiteTurn {
		return eval
	} else {
		return -eval
	}

	// TODO: consider material placement - towards center usually is good
	// TODO: consider passed pawns
}

func ExpMoves(pos *Position) int {
	// just arbitrary numbers at this point
	// TODO: heavily rely on material, instead of prior moves made
	return MaxInt(80-int(pos.ply), 30)
}

func SortMoves(pos *Position, moves []Move, anticipated Move) []Move {
	// quicksort impl
	if len(moves) < 2 {
		return moves
	} else if len(moves) == 2 {
		if anticipated != NULL_MOVE {
			if anticipated == moves[0] {
				return moves
			} else if anticipated == moves[1] {
				Swap(moves, 0, 1)
				return moves
			}
		}
		if !compareMoves(pos, moves[0], moves[1]) {
			Swap(moves, 0, 1)
		}
		return moves
	}

	mid := len(moves) / 2
	left := SortMoves(pos, moves[0:mid], anticipated)
	right := SortMoves(pos, moves[mid:], anticipated)
	var leftIdx = 0
	var rightIdx = 0
	var rtnIdx = 0
	rtn := make([]Move, len(left)+len(right))
	for leftIdx < len(left) && rightIdx < len(right) {
		leftMove := left[leftIdx]
		rightMove := right[rightIdx]
		if anticipated != NULL_MOVE && anticipated == leftMove {
			rtn[rtnIdx] = leftMove
			leftIdx++
		} else if anticipated != NULL_MOVE && anticipated == rightMove {
			rtn[rtnIdx] = rightMove
			rightIdx++
		} else if compareMoves(pos, leftMove, rightMove) {
			rtn[rtnIdx] = leftMove
			leftIdx++
		} else {
			rtn[rtnIdx] = rightMove
			rightIdx++
		}
		rtnIdx++
	}
	for leftIdx < len(left) {
		rtn[rtnIdx] = left[leftIdx]
		rtnIdx++
		leftIdx++
	}
	for rightIdx < len(right) {
		rtn[rtnIdx] = right[rightIdx]
		rtnIdx++
		rightIdx++
	}
	return rtn
}

// compareMoves returns true if the first move (moveA) is expected to be better than the second (moveB)
func compareMoves(pos *Position, moveA, moveB Move) bool {
	return EvalMove(pos, moveA) >= EvalMove(pos, moveB)
}

const (
	CHECK_VALUE = 900
)

func EvalMove(pos *Position, move Move) int16 {
	var moveVal int16
	moveVal += PieceTypeToVal(move.PromotedTo())

	// TODO: factor in checks once captures are fast to determine (pins)
	//captPiece, lastFrozenPos := pos.MakeMove(move)
	//defer pos.UnmakeMove(move, lastFrozenPos, captPiece)
	//if pos.IsKingChecked() {
	//	moveVal += CHECK_VALUE
	//}

	captPiece := pos.pieces[move.EndSq()]
	if move.Type() == CAPTURES_EN_PASSANT {
		captPiece = NewPiece(PAWN, NewColor(!pos.isWhiteTurn))
	}
	moveVal += PieceTypeToVal(captPiece.Type())
	return moveVal
}

func PieceTypeToVal(pt PieceType) int16 {
	if pt == PAWN {
		return PAWN_VAL
	} else if pt == KNIGHT {
		return KNIGHT_VAL
	} else if pt == BISHOP {
		return BISHOP_VAL
	} else if pt == ROOK {
		return ROOK_VAL
	} else if pt == QUEEN {
		return QUEEN_VAL
	}
	return 0
}
