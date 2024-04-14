package main

import "log"

//func GenLegalMoves(state *State) []Move {
//
//}
//
//func GenPseudoLegalMoves(state *State) []Move {
//
//}

func GenPseudoLegalPawnMoves(pos *Position, sq, epSq Square) []Move {
	piece := pos.pieces[sq]
	if DEBUG {
		if piece.Type() != PAWN {
			log.Fatalf("cannot generate pawn move for non-pawn piece in pos:\n%s", pos)
		}
	}
	file := sq.File()
	isWhite := piece.IsWhite()

	rtn := make([]Move, 0)
	if file < 8 {
		var rAttackSq Square
		if isWhite {
			rAttackSq = sq + 9
		} else {
			rAttackSq = sq - 7
		}
		attackedPiece := pos.pieces[rAttackSq]
		if attackedPiece != EMPTY {
			if attackedPiece.IsWhite() != isWhite {
				if rAttackSq >= SQ_A8 || rAttackSq <= SQ_H1 {
					rtn = append(rtn, NewPromoteMoves(sq, rAttackSq)...)
				} else {
					rtn = append(rtn, NewNormalMove(sq, rAttackSq))
				}
			}
		} else {
			if rAttackSq == epSq {
				rtn = append(rtn, NewEnPassantMove(sq, rAttackSq))
			}
		}
	}
	if file > 1 {
		var lAttackSq Square
		if isWhite {
			lAttackSq = sq + 7
		} else {
			lAttackSq = sq - 9
		}
		attackedPiece := pos.pieces[lAttackSq]
		if attackedPiece != EMPTY {
			if attackedPiece.IsWhite() != isWhite {
				if lAttackSq >= SQ_A8 || lAttackSq <= SQ_H1 {
					rtn = append(rtn, NewPromoteMoves(sq, lAttackSq)...)
				} else {
					rtn = append(rtn, NewNormalMove(sq, lAttackSq))
				}
			}
		} else {
			if lAttackSq == epSq {
				rtn = append(rtn, NewEnPassantMove(sq, lAttackSq))
			}
		}
	}

	var sqInFront Square
	if isWhite {
		sqInFront = sq + 8
	} else {
		sqInFront = sq - 8
	}
	pieceInFront := pos.pieces[sqInFront]
	if pieceInFront != EMPTY {
		return rtn
	}
	if sqInFront >= SQ_A8 || sqInFront <= SQ_H1 {
		rtn = append(rtn, NewPromoteMoves(sq, sqInFront)...)
		return rtn
	} else {
		rtn = append(rtn, NewNormalMove(sq, sqInFront))
	}
	if isWhite && sq >= SQ_A2 && sq <= SQ_H2 {
		sqTwoInFront := sq + 16
		if pos.pieces[sqTwoInFront] == EMPTY {
			rtn = append(rtn, NewNormalMove(sq, sqTwoInFront))
		}
	} else if !isWhite && sq >= SQ_A7 && sq <= SQ_H7 {
		sqTwoInFront := sq - 16
		if pos.pieces[sqTwoInFront] == EMPTY {
			rtn = append(rtn, NewNormalMove(sq, sqTwoInFront))
		}
	}
	rtn = append(rtn, NewNormalMove(sq, sqInFront))
	return rtn
}

//func GenPseudoLegalKnightMoves(state *State, sq Square) []Move {
//
//}
//
//func GenPseudoLegalBishopMoves(state *State, sq Square) []Move {
//
//}
//
//func GenPseudoLegalRookMoves(state *State, sq Square) []Move {
//
//}
//
//func GenPseudoLegalQueenMoves(state *State, sq Square) []Move {
//
//}
//
//func GenPseudoLegalKingMoves(state *State, sq Square) []Move {
//
//}
