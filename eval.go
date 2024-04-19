package main

import "github.com/CameronHonis/chess"

const (
	PAWN_VAL   = 100
	KNIGHT_VAL = 280
	BISHOP_VAL = 320
	ROOK_VAL   = 500
	QUEEN_VAL  = 950
	MATE_VAL   = 10000
	DRAW_VAL   = -50
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
	// TODO: attacks (direct and/or discoveries) - Maybe too slow to justify?

	// TODO: positional eval (square value by piece lookups and castling)

	// TODO: passed pawns
}

func ExpMoves(pos *Position) int {
	// just arbitrary numbers at this point
	// TODO: heavily rely on material, instead of prior moves made
	return MaxInt(80-int(pos.ply), 30)
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
		return MATE_VAL
	} else {
		return 0
	}
}

func SortMoves(pos *chess.Board, moves []*chess.Move, anticipated *chess.Move) []*chess.Move {
	// quicksort impl
	if len(moves) == 1 {
		return moves
	} else if len(moves) == 2 {
		if anticipated != nil {
			if anticipated.Equal(moves[0]) {
				return moves
			} else if anticipated.Equal(moves[1]) {
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
	rtn := make([]*chess.Move, len(left)+len(right))
	for leftIdx < len(left) && rightIdx < len(right) {
		leftMove := left[leftIdx]
		rightMove := right[rightIdx]
		if anticipated != nil && anticipated.Equal(leftMove) {
			rtn[rtnIdx] = leftMove
			leftIdx++
		} else if anticipated != nil && anticipated.Equal(rightMove) {
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
