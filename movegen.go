package main

import "log"

type LegalMoveIter struct {
	state  *State
	pMoves []Move
	idx    int
}

func NewLegalMoveIter(state *State) *LegalMoveIter {
	return &LegalMoveIter{
		state:  state,
		pMoves: GenPseudoLegalMoves(state),
		idx:    0,
	}
}

func (iter *LegalMoveIter) Next() (move Move, done bool) {
	for iter.idx < len(iter.pMoves) {
		pMove := iter.pMoves[iter.idx]
		iter.idx++

		if iter.state.Pos.IsLegalMove(pMove) {
			return pMove, false
		}
	}
	return NULL_MOVE, true
}

func GenPseudoLegalMoves(state *State) []Move {
	pos := state.Pos
	rtn := make([]Move, 0)
	for rank := 1; rank < 9; rank++ {
		for file := 1; file < 9; file++ {
			sq := SqFromCoords(rank, file)
			piece := pos.pieces[sq]
			pt := piece.Type()
			if pt == PAWN {
				moves := GenPseudoLegalPawnMoves(pos, sq, state.EnPassantSq)
				rtn = append(rtn, moves...)
			} else if pt == KING {
				isWhite := piece.IsWhite()
				var canCastleKS bool
				var canCastleQS bool
				if isWhite {
					canCastleKS = state.CastleRights[W_CASTLE_KINGSIDE_RIGHT]
					canCastleQS = state.CastleRights[W_CASTLE_QUEENSIDE_RIGHT]
				} else {
					canCastleKS = state.CastleRights[B_CASTLE_KINGSIDE_RIGHT]
					canCastleQS = state.CastleRights[B_CASTLE_QUEENSIDE_RIGHT]
				}
				moves := GenPseudoLegalKingMoves(pos, sq, canCastleKS, canCastleQS)
				rtn = append(rtn, moves...)
			} else {
				moves := GenPseudoLegalNimblePieceMoves(pos, sq)
				rtn = append(rtn, moves...)
			}
		}
	}
	return rtn
}

func GenPseudoLegalPawnMoves(pos *Position, sq, epSq Square) []Move {
	piece := pos.pieces[sq]
	if DEBUG {
		if piece.Type() != PAWN {
			log.Fatalf("cannot generate pawn move for non-pawn piece on %s in pos:\n%s", sq, pos)
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

func GenPseudoLegalNimblePieceMoves(pos *Position, sq Square) []Move {
	piece := pos.pieces[sq]
	pt := piece.Type()
	if DEBUG {
		if pt != KNIGHT && pt != BISHOP && pt != ROOK && pt != QUEEN {
			log.Fatalf("cannot generate nimble piece move for piece on %d in pos:\n%s", sq, pos)
		}
	}
	var attackBB Bitboard
	if pt == KNIGHT {
		attackBB = KnightAttacksBB(sq)
	} else {
		attackBB = SlidingAttacksBB(pos.OccupiedBB(), sq, pt)
	}
	isWhite := piece.IsWhite()

	rtn := make([]Move, 0)
	for attackBB > 0 {
		var attackSq Square
		attackSq, attackBB = attackBB.PopFirstSq()
		attackedPiece := pos.pieces[attackSq]
		if attackedPiece == EMPTY || attackedPiece.IsWhite() != isWhite {
			rtn = append(rtn, NewNormalMove(sq, attackSq))
		}
	}
	return rtn
}

func GenPseudoLegalKingMoves(pos *Position, sq Square, canCastleKS, canCastleQS bool) []Move {
	piece := pos.pieces[sq]
	if DEBUG {
		if piece.Type() != KING {
			log.Fatalf("cannot generate king moves for non-king piece on %d in pos:\n%s", sq, pos)
		}
	}
	attackBB := KingAttacksBB(sq)
	isWhite := piece.IsWhite()

	rtn := make([]Move, 0)
	for attackBB > 0 {
		var attackSq Square
		attackSq, attackBB = attackBB.PopFirstSq()
		attackedPiece := pos.pieces[attackSq]
		if attackedPiece == EMPTY || attackedPiece.IsWhite() != isWhite {
			rtn = append(rtn, NewNormalMove(sq, attackSq))
		}
	}

	occupiedBB := pos.OccupiedBB()
	if canCastleKS {
		castlePathMask := BBWithSquares(sq+1, sq+2)
		if occupiedBB&castlePathMask == 0 {
			rtn = append(rtn, NewMove(sq, sq+2, NULL_SQ, EMPTY_PIECE_TYPE, true))
		}
	}
	if canCastleQS {
		castlePathMask := BBWithSquares(sq-1, sq-2)
		if occupiedBB&castlePathMask == 0 {
			rtn = append(rtn, NewMove(sq, sq-2, NULL_SQ, EMPTY_PIECE_TYPE, true))
		}
	}
	return rtn
}
