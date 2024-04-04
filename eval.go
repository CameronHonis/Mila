package main

import "github.com/CameronHonis/chess"

const (
	PAWN_VAL   = 100
	KNIGHT_VAL = 280
	BISHOP_VAL = 320
	ROOK_VAL   = 500
	QUEEN_VAL  = 950
)

func Eval(pos *chess.Board) float64 {
	if pos.Result == chess.BOARD_RESULT_WHITE_WINS_BY_CHECKMATE {
		return 30
	} else if pos.Result == chess.BOARD_RESULT_BLACK_WINS_BY_CHECKMATE {
		return -30
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
}

func ExpMoves(pos *chess.Board) int {
	// just arbitrary numbers at this point
	return MaxInt(80-int(pos.FullMoveCount), 30)
}
