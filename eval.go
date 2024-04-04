package main

import "github.com/CameronHonis/chess"

const (
	PAWN_VAL   = 100
	KNIGHT_VAL = 280
	BISHOP_VAL = 320
	ROOK_VAL   = 500
	QUEEN_VAL  = 950
	KING_VAL   = 10000
)

func EvalPos(pos *chess.Board) float64 {
	if pos.Result == chess.BOARD_RESULT_WHITE_WINS_BY_CHECKMATE {
		return KING_VAL / 100
	} else if pos.Result == chess.BOARD_RESULT_BLACK_WINS_BY_CHECKMATE {
		return -KING_VAL / 100
	} else if pos.Result != chess.BOARD_RESULT_IN_PROGRESS {
		return (-0.5 * PAWN_VAL) / 100
	}
	matCount := pos.ComputeMaterialCount()
	pawnDiff := int(matCount.WhitePawnCount) - int(matCount.BlackPawnCount)
	knightDiff := int(matCount.WhiteKnightCount) - int(matCount.BlackKnightCount)
	bishopDiff := int(matCount.WhiteLightBishopCount) + int(matCount.WhiteDarkBishopCount) -
		int(matCount.BlackLightBishopCount) - int(matCount.BlackDarkBishopCount)
	rookDiff := int(matCount.WhiteRookCount) - int(matCount.BlackRookCount)
	queenDiff := int(matCount.WhiteQueenCount) - int(matCount.BlackQueenCount)
	rawScore := PAWN_VAL*pawnDiff + KNIGHT_VAL*knightDiff + BISHOP_VAL*bishopDiff +
		ROOK_VAL*rookDiff + QUEEN_VAL*queenDiff
	return float64(rawScore) / 100
	// TODO: attacks (direct and/or discoveries) - Maybe too slow to justify?

	// TODO: positional eval (square value by piece lookups and castling)

	// TODO: passed pawns
}

func ExpMoves(pos *chess.Board) int {
	// just arbitrary numbers at this point
	return MaxInt(80-int(pos.FullMoveCount), 30)
}

func PieceValue(piece chess.Piece) int {
	if piece.IsPawn() {
		return PAWN_VAL
	} else if piece.IsKnight() {
		return KNIGHT_VAL
	} else if piece.IsBishop() {
		return BISHOP_VAL
	} else if piece.IsRook() {
		return ROOK_VAL
	} else if piece.IsQueen() {
		return QUEEN_VAL
	} else if piece.IsKing() {
		return KING_VAL
	} else {
		return 0
	}
}

func SortMoves(pos *chess.Board, moves []*chess.Move) []*chess.Move {
	// quicksort impl
	if len(moves) == 1 {
		return moves
	} else if len(moves) == 2 {
		if !compareMoves(pos, moves[0], moves[1]) {
			Swap(moves, 0, 1)
		}
		return moves
	}

	mid := len(moves) / 2
	left := SortMoves(pos, moves[0:mid])
	right := SortMoves(pos, moves[mid:])
	var leftIdx = 0
	var rightIdx = 0
	var rtnIdx = 0
	rtn := make([]*chess.Move, len(left)+len(right))
	for leftIdx < len(left) && rightIdx < len(right) {
		leftMove := left[leftIdx]
		rightMove := right[rightIdx]
		if compareMoves(pos, leftMove, rightMove) {
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
func compareMoves(pos *chess.Board, moveA, moveB *chess.Move) bool {
	return EvalMove(pos, moveA) >= EvalMove(pos, moveB)
}

const (
	CHECK_VALUE = 900
)

func EvalMove(pos *chess.Board, move *chess.Move) int {
	// uses the age-old adage "always look for checks, captures, then attacks"
	var moveVal = 0
	if move.PawnUpgradedTo != chess.EMPTY {
		moveVal += PieceValue(move.PawnUpgradedTo) - PAWN_VAL
	}
	moveVal += CHECK_VALUE * len(move.KingCheckingSquares)
	moveVal += PieceValue(move.CapturedPiece)
	return moveVal
}
