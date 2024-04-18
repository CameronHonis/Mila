package main

import "log"

type LegalMoveIter struct {
	pos    *Position
	pMoves []Move
	idx    int
}

func NewLegalMoveIter(pos *Position) *LegalMoveIter {
	return &LegalMoveIter{
		pMoves: GenPseudoLegalMoves(pos),
		idx:    0,
	}
}

func (iter *LegalMoveIter) Next() (move Move, done bool) {
	for iter.idx < len(iter.pMoves) {
		pMove := iter.pMoves[iter.idx]
		iter.idx++

		if iter.pos.IsLegalMove(pMove) {
			return pMove, false
		}
	}
	return NULL_MOVE, true
}

func GenPseudoLegalMoves(pos *Position) []Move {
	rtn := make([]Move, 0)
	for sq := SQ_A1; sq < N_SQUARES; sq++ {
		piece := pos.pieces[sq]
		pt := piece.Type()
		if pt == PAWN {
			moves := GenPseudoLegalPawnMoves(pos, sq)
			rtn = append(rtn, moves...)
		} else if pt == KING {
			moves := GenPseudoLegalKingMoves(pos, sq)
			rtn = append(rtn, moves...)
		} else {
			moves := GenPseudoLegalNimblePieceMoves(pos, sq)
			rtn = append(rtn, moves...)
		}
	}
	return rtn
}

func GenPseudoLegalPawnMoves(pos *Position, sq Square) []Move {
	piece := pos.pieces[sq]
	if DEBUG {
		if piece.Type() != PAWN {
			log.Fatalf("cannot generate pawn move for non-pawn piece on %s in pos:\n%s", sq, pos)
		}
	}

	rtn := make([]Move, 0)

	var attacksBB = PawnAttacksBB(sq, piece.Color())
	isWhite := piece.IsWhite()
	for attacksBB > 0 {
		var attackSq Square
		attackSq, attacksBB = attacksBB.PopFirstSq()
		attackedPiece := pos.pieces[attackSq]
		if attackedPiece != EMPTY && attackedPiece.IsWhite() != isWhite {
			rank := attackSq.Rank()
			if rank == 1 || rank == 8 {
				rtn = append(rtn, NewPromoteMoves(sq, attackSq)...)
			} else {
				rtn = append(rtn, NewNormalMove(sq, attackSq))
			}
		} else if attackSq == pos.frozenPos.EnPassantSq {
			rtn = append(rtn, NewEnPassantMove(sq, attackSq))
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

func GenPseudoLegalKingMoves(pos *Position, sq Square) []Move {
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
	castleRights := pos.frozenPos.CastleRights
	var canCastleKS bool
	var canCastleQS bool
	if isWhite {
		canCastleKS = castleRights[W_CASTLE_KINGSIDE_RIGHT]
		canCastleQS = castleRights[W_CASTLE_QUEENSIDE_RIGHT]
	} else {
		canCastleKS = castleRights[B_CASTLE_KINGSIDE_RIGHT]
		canCastleQS = castleRights[B_CASTLE_QUEENSIDE_RIGHT]
	}
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

func SlidingAttacksBB(occupied Bitboard, sq Square, pt PieceType) Bitboard {
	initAttacks()
	var rtn Bitboard
	rank := sq.Rank()
	if pt == ROOK || pt == QUEEN {
		file := sq.File()
		rankMask := BBWithRank(rank, 0b11111111)
		occupiedRankBB := occupied & rankMask
		rtn |= rankAttacks[occupiedRankBB][file-1]

		fileMask := BBWithFile(file, 0b11111111)
		occupiedFileBB := occupied & fileMask
		rtn |= fileAttacks[occupiedFileBB][rank-1]
	}
	if pt == BISHOP || pt == QUEEN {
		posDiagMask := BBWithPosDiag(sq.PosDiagIdx(), 0b11111111)
		occupiedPosDiagBB := occupied & posDiagMask
		rtn |= posDiagAttacks[occupiedPosDiagBB][rank-1]

		negDiagMask := BBWithNegDiag(sq.NegDiagIdx(), 0b11111111)
		occupiedNegDiagBB := occupied & negDiagMask
		rtn |= negDiagAttacks[occupiedNegDiagBB][rank-1]
	}
	return rtn
}

func KnightAttacksBB(sq Square) Bitboard {
	initAttacks()
	return knightAttacks[sq]
}

func KingAttacksBB(sq Square) Bitboard {
	initAttacks()
	return kingAttacks[sq]
}

func PawnAttacksBB(sq Square, color Color) Bitboard {
	initAttacks()
	return pawnAttacks[sq][color]
}
