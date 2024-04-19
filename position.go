package main

import (
	"fmt"
	"strconv"
	"strings"
)

// FrozenPos encapsulates the bits about a position that cannot be re-generated
// after making + unmaking a move. This struct is meant to be provided alongside
// a move to unmake a move and restore a position.
type FrozenPos struct {
	EnPassantSq  Square
	Rule50       Ply
	CastleRights [N_CASTLE_RIGHTS]bool
}

func (fp *FrozenPos) Copy() *FrozenPos {
	fpCopy := *fp
	return &fpCopy
}

type Position struct {
	pieces         [N_SQUARES]Piece
	pieceBitboards [N_PIECES]Bitboard
	colorBitboards [N_COLORS]Bitboard
	material       Material
	repetitions    map[ZHash]uint8
	ply            Ply
	hash           ZHash
	result         Result // only covers non-checkmate/stalemate positions
	isWhiteTurn    bool
	isKingChecked  bool

	frozenPos *FrozenPos
}

func InitPos() *Position {
	pos := &Position{
		pieces: [N_SQUARES]Piece{
			W_ROOK, W_KNIGHT, W_BISHOP, W_QUEEN, W_KING, W_BISHOP, W_KNIGHT, W_ROOK,
			W_PAWN, W_PAWN, W_PAWN, W_PAWN, W_PAWN, W_PAWN, W_PAWN, W_PAWN,
			EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY,
			EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY,
			EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY,
			EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY,
			B_PAWN, B_PAWN, B_PAWN, B_PAWN, B_PAWN, B_PAWN, B_PAWN, B_PAWN,
			B_ROOK, B_KNIGHT, B_BISHOP, B_QUEEN, B_KING, B_BISHOP, B_KNIGHT, B_ROOK,
		},
		pieceBitboards: [N_PIECES]Bitboard{
			0b00000000_00000000_11111111_11111111_11111111_11111111_00000000_00000000,
			0b00000000_00000000_00000000_00000000_00000000_00000000_11111111_00000000,
			0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_01000010,
			0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00100100,
			0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_10000001,
			0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00001000,
			0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00010000,
			0b00000000_11111111_00000000_00000000_00000000_00000000_00000000_00000000,
			0b01000010_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
			0b00100100_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
			0b10000001_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
			0b00001000_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
			0b00010000_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
		},
		colorBitboards: [N_COLORS]Bitboard{
			0b00000000_00000000_00000000_00000000_00000000_00000000_11111111_11111111,
			0b11111111_11111111_00000000_00000000_00000000_00000000_00000000_00000000,
		},
		material:    InitMaterial(),
		repetitions: make(map[ZHash]uint8),
		ply:         0,
		result:      RESULT_IN_PROGRESS,
		isWhiteTurn: true,
		frozenPos: &FrozenPos{
			EnPassantSq:  NULL_SQ,
			CastleRights: [N_CASTLE_RIGHTS]bool{true, true, true, true},
		},
	}
	pos.hash = NewZHash(pos)
	return pos
}

func FromFEN(fen string) (*Position, error) {
	var pos = &Position{
		repetitions: make(map[ZHash]uint8),
		frozenPos:   &FrozenPos{},
	}
	fenSegs := strings.Split(fen, " ")
	if len(fenSegs) != 6 {
		return nil, fmt.Errorf("invalid number of fen segments %d, expected 6", len(fenSegs))
	}

	piecesFen := fenSegs[0]
	fenPiecesRows := strings.Split(piecesFen, "/")
	if len(fenPiecesRows) != 8 {
		return nil, fmt.Errorf("invalid number of rows in FEN pieces %d, expected 8", len(fenPiecesRows))
	}
	for fenRowIdx, fenPiecesRow := range fenPiecesRows {
		rank := 8 - fenRowIdx
		var file = 1
		for _, fenPiece := range []byte(fenPiecesRow) {
			if file > 8 {
				return nil, fmt.Errorf("too many pieces on rank %d in fen %s", rank, fen)
			}
			if fenPiece >= '1' && fenPiece <= '8' {
				for i := 0; i < int(fenPiece-'0'); i++ {
					sq := SqFromCoords(rank, file)
					pos.pieceBitboards[EMPTY] |= BBWithSquares(sq)
					file++
				}
				continue
			}

			sq := SqFromCoords(rank, file)
			piece := PieceFromChar(fenPiece)
			pos.pieces[sq] = piece
			pos.pieceBitboards[piece] |= BBWithSquares(sq)
			if piece.IsWhite() {
				pos.colorBitboards[WHITE] |= BBWithSquares(sq)
			} else {
				pos.colorBitboards[BLACK] |= BBWithSquares(sq)
			}
			if piece.Type() != KING {
				pos.material.AddPiece(piece, sq)
			}
			file++
		}
		if file < 9 {
			return nil, fmt.Errorf("not enough pieces on rank %d in fen %s", rank, fen)
		}
	}

	turnSpecifier := fenSegs[1]
	if turnSpecifier == "w" {
		pos.isWhiteTurn = true
	} else if turnSpecifier == "b" {
		pos.isWhiteTurn = false
	} else {
		return nil, fmt.Errorf("unexpected length for turn specifier %s in fen %s", fenSegs[1], fen)
	}

	castleRightsSpecifier := fenSegs[2]
	if castleRightsSpecifier != "-" {
		for _, castleRightChar := range []byte(castleRightsSpecifier) {
			if castleRightChar == 'K' {
				pos.frozenPos.CastleRights[W_CASTLE_KINGSIDE_RIGHT] = true
			} else if castleRightChar == 'Q' {
				pos.frozenPos.CastleRights[W_CASTLE_QUEENSIDE_RIGHT] = true
			} else if castleRightChar == 'k' {
				pos.frozenPos.CastleRights[B_CASTLE_KINGSIDE_RIGHT] = true
			} else if castleRightChar == 'q' {
				pos.frozenPos.CastleRights[B_CASTLE_QUEENSIDE_RIGHT] = true
			} else {
				return nil, fmt.Errorf("could not parse castle rights, unknown char %c in %s", castleRightChar, castleRightsSpecifier)
			}
		}
	}

	epSpecifier := fenSegs[3]
	if epSpecifier == "-" {
		pos.frozenPos.EnPassantSq = NULL_SQ
	} else {
		epSq, epSqErr := SqFromAlg(epSpecifier)
		if epSqErr != nil {
			return nil, fmt.Errorf("could not parse en passant square in fen %s: %s", fen, epSqErr)
		}
		pos.frozenPos.EnPassantSq = epSq
	}

	halfmovesSpecifier := fenSegs[4]
	if halfmoves, halfmoveErr := strconv.Atoi(halfmovesSpecifier); halfmoveErr != nil {
		return nil, fmt.Errorf("could not parse halfmoves: %s", halfmoveErr)
	} else {
		pos.frozenPos.Rule50 = Ply(halfmoves)
	}

	movesSpecifier := fenSegs[5]
	if nMoves, nMovesErr := strconv.Atoi(movesSpecifier); nMovesErr != nil {
		return nil, fmt.Errorf("could not parse number of moves: %s", nMovesErr)
	} else {
		pos.ply = PlyFromNMoves(uint(nMoves), pos.isWhiteTurn)
	}

	pos.hash = NewZHash(pos)

	return pos, nil
}

func (p *Position) String() string {
	var rtnBuilder = strings.Builder{}
	fenPieces := strings.Split(p.FEN(), " ")
	shortFen := strings.Join(fenPieces[1:], " ")
	rtnBuilder.WriteString(shortFen)
	rtnBuilder.WriteByte('\n')

	for rank := 8; rank > 0; rank-- {
		rtnBuilder.WriteString(fmt.Sprintf("%d ", rank))
		for file := 1; file < 9; file++ {
			idx := 8*(rank-1) + (file - 1)
			piece := p.pieces[idx]
			pieceStr := piece.String()
			if piece > 0 {
				rtnBuilder.WriteString(fmt.Sprintf("%s ", pieceStr))
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

func (p *Position) FEN() string {
	var rtnBuilder strings.Builder
	for rank := 8; rank > 0; rank-- {
		var consecSpaces int
		for file := 1; file < 9; file++ {
			sq := SqFromCoords(rank, file)
			if sq == 64 {
				fmt.Println(rank, file)
			}
			piece := p.pieces[sq]
			if piece == EMPTY {
				consecSpaces++
			} else {
				if consecSpaces > 0 {
					rtnBuilder.WriteByte(byte('0' + consecSpaces))
				}
				rtnBuilder.WriteByte(piece.Char())
			}
		}
		if consecSpaces > 0 {
			rtnBuilder.WriteByte(byte('0' + consecSpaces))
		}
		if rank != 1 {
			rtnBuilder.WriteByte('/')
		}
	}
	rtnBuilder.WriteByte(' ')
	if p.isWhiteTurn {
		rtnBuilder.WriteByte('w')
	} else {
		rtnBuilder.WriteByte('b')
	}
	rtnBuilder.WriteByte(' ')

	var anyCastleRights bool
	if p.frozenPos.CastleRights[W_CASTLE_KINGSIDE_RIGHT] {
		rtnBuilder.WriteByte('K')
		anyCastleRights = true
	}
	if p.frozenPos.CastleRights[W_CASTLE_QUEENSIDE_RIGHT] {
		rtnBuilder.WriteByte('Q')
		anyCastleRights = true
	}
	if p.frozenPos.CastleRights[B_CASTLE_KINGSIDE_RIGHT] {
		rtnBuilder.WriteByte('k')
		anyCastleRights = true
	}
	if p.frozenPos.CastleRights[B_CASTLE_QUEENSIDE_RIGHT] {
		rtnBuilder.WriteByte('q')
		anyCastleRights = true
	}
	if !anyCastleRights {
		rtnBuilder.WriteByte('-')
	}
	rtnBuilder.WriteByte(' ')

	if p.frozenPos.EnPassantSq.IsNull() {
		rtnBuilder.WriteByte('-')
	} else {
		rtnBuilder.WriteString(p.frozenPos.EnPassantSq.String())
	}
	rtnBuilder.WriteByte(' ')

	rtnBuilder.WriteString(strconv.Itoa(int(p.frozenPos.Rule50)))
	rtnBuilder.WriteByte(' ')

	rtnBuilder.WriteString(strconv.Itoa(int(NMovesFromPly(p.ply))))
	return rtnBuilder.String()
}

func (p *Position) OccupiedBB() Bitboard {
	var rtn Bitboard
	for _, colorBB := range p.colorBitboards {
		rtn ^= colorBB
	}
	return rtn
}

// IsLegalMove is intended to filter out only valid pseudo-legal moves.
func (p *Position) IsLegalMove(pMove Move) bool {
	piece := p.pieces[pMove.StartSq()]
	pt := piece.Type()
	color := piece.Color()

	if pt == KING && pMove.Type() == CASTLING {
		start := pMove.StartSq()
		end := pMove.EndSq()
		sq0 := start
		var sq1 Square
		if end > start {
			sq1 = sq0 + 1
		} else {
			sq1 = sq0 - 1
		}
		sq2 := end
		return !p.isSquareAttacked(color.Opp(), sq0) &&
			!p.isSquareAttacked(color.Opp(), sq1) &&
			!p.isSquareAttacked(color.Opp(), sq2)
	} else {
		capturedPiece, fPos := p.MakeMove(pMove)
		kingSq := p.pieceBitboards[NewPiece(KING, color)].FirstSq()
		defer p.UnmakeMove(pMove, fPos, capturedPiece)

		return !p.isSquareAttacked(color.Opp(), kingSq)
	}
}

// MakeMove expects the inbound move to be filtered by Position.IsLegalMove
func (p *Position) MakeMove(move Move) (captured Piece, lastFrozenPos *FrozenPos) {
	mt := move.Type()

	p.ply++

	lastFrozenPos = p.frozenPos
	p.updateFrozenPos(move)

	if lastFrozenPos.EnPassantSq != p.frozenPos.EnPassantSq {
		p.hash = p.hash.UpdateEnPassantSq(lastFrozenPos.EnPassantSq, p.frozenPos.EnPassantSq)
	}
	for castleRight := W_CASTLE_KINGSIDE_RIGHT; castleRight < N_CASTLE_RIGHTS; castleRight++ {
		if lastFrozenPos.CastleRights[castleRight] != p.frozenPos.CastleRights[castleRight] {
			p.hash = p.hash.ToggleCastleRight(castleRight)
		}
	}

	if mt == CASTLING {
		p.doCastle(move)
	} else if mt == CAPTURES_EN_PASSANT {
		captured = p.doEnPassant(move)
	} else if mt == PAWN_PROMOTION {
		captured = p.doPromote(move)
	} else { // NORMAL MOVE
		captured = p.movePiece(move.StartSq(), move.EndSq())
	}

	p.isWhiteTurn = !p.isWhiteTurn
	p.hash = p.hash.ToggleTurn()

	if _, ok := p.repetitions[p.hash]; !ok {
		p.repetitions[p.hash] = 0
	}
	p.repetitions[p.hash]++

	if p.repetitions[p.hash] >= 3 {
		p.result = RESULT_DRAW_REPETITION
	} else if p.material.IsForcedDraw() {
		p.result = RESULT_DRAW_MATL
	} else if p.frozenPos.Rule50 >= 50 {
		p.result = RESULT_DRAW_RULE50
	}

	return
}

func (p *Position) UnmakeMove(move Move, fp *FrozenPos, captured Piece) {
	p.result = RESULT_IN_PROGRESS
	p.repetitions[p.hash]--
	if p.repetitions[p.hash] == 0 {
		delete(p.repetitions, p.hash)
	}

	if fp.EnPassantSq != p.frozenPos.EnPassantSq {
		p.hash = p.hash.UpdateEnPassantSq(p.frozenPos.EnPassantSq, fp.EnPassantSq)
	}
	for castleRight := W_CASTLE_KINGSIDE_RIGHT; castleRight < N_CASTLE_RIGHTS; castleRight++ {
		if fp.CastleRights[castleRight] != p.frozenPos.CastleRights[castleRight] {
			p.hash = p.hash.ToggleCastleRight(castleRight)
		}
	}
	p.frozenPos = fp

	p.ply--

	mt := move.Type()
	if mt == CASTLING {
		p.undoCastle(move)
	} else if mt == CAPTURES_EN_PASSANT {
		p.undoEnPassant(move, captured)
	} else if mt == PAWN_PROMOTION {
		p.undoPromote(move, captured)
	} else { // NORMAL
		p.movePiece(move.EndSq(), move.StartSq())
		p.addPiece(move.EndSq(), captured)
	}

	p.isWhiteTurn = !p.isWhiteTurn
	p.hash = p.hash.ToggleTurn()
}

func (p *Position) IsKingChecked() bool {
	color := NewColor(p.isWhiteTurn)
	piece := NewPiece(KING, color)
	sq := p.pieceBitboards[piece].FirstSq()
	return p.isSquareAttacked(color.Opp(), sq)
}

func (p *Position) IsMate() bool {
	return p.IsKingChecked() && !p.HasLegalMoves()
}

func (p *Position) HasLegalMoves() bool {
	_, hasNoLegalMoves := NewLegalMoveIter(p).Next()
	return !hasNoLegalMoves
}

func (p *Position) doCastle(move Move) {
	start := move.StartSq()
	end := move.EndSq()
	endFile := end.File()
	var rookStartSq Square
	var rookEndSq Square
	if endFile == 7 {
		rookStartSq = end + 1
		rookEndSq = end - 1
	} else {
		rookStartSq = end - 2
		rookEndSq = end + 1
	}
	p.movePiece(start, end)
	p.movePiece(rookStartSq, rookEndSq)
}

func (p *Position) undoCastle(move Move) {
	start := move.StartSq()
	end := move.EndSq()
	endFile := end.File()
	var rookStartSq Square
	var rookEndSq Square
	if endFile == 7 {
		rookStartSq = end + 1
		rookEndSq = end - 1
	} else {
		rookStartSq = end - 2
		rookEndSq = end + 1
	}
	p.movePiece(end, start)
	p.movePiece(rookEndSq, rookStartSq)
}

func (p *Position) doEnPassant(move Move) (captured Piece) {
	start := move.StartSq()
	end := move.EndSq()
	epSq := SqFromCoords(int(start.Rank()), int(end.File()))
	p.movePiece(start, end)
	captured = p.removePiece(epSq)
	return
}

func (p *Position) undoEnPassant(move Move, captured Piece) {
	start := move.StartSq()
	end := move.EndSq()
	epSq := SqFromCoords(int(start.Rank()), int(end.File()))
	p.movePiece(end, start)
	p.addPiece(epSq, captured)
}

func (p *Position) doPromote(move Move) (captured Piece) {
	start := move.StartSq()
	isWhite := p.pieces[start].IsWhite()
	end := move.EndSq()

	captured = p.movePiece(start, end)

	prPiece := NewPiece(move.PromotedTo(), NewColor(isWhite))
	p.replacePiece(end, prPiece)
	return
}

func (p *Position) undoPromote(move Move, captured Piece) {
	start := move.StartSq()
	end := move.EndSq()
	isWhite := p.pieces[end].IsWhite()
	p.replacePiece(end, NewPiece(PAWN, NewColor(isWhite)))
	p.movePiece(end, start)
	p.addPiece(end, captured)
}

// removePiece removes the piece and keeps the bitboards consistent
// NOTE: this does not affect the frozen fields
func (p *Position) removePiece(sq Square) (removed Piece) {
	piece := p.pieces[sq]
	if piece != EMPTY {
		color := NewColor(piece.IsWhite())
		mask := BBWithSquares(sq)

		p.pieces[sq] = EMPTY
		p.pieceBitboards[piece] ^= mask
		p.pieceBitboards[EMPTY] ^= mask
		p.colorBitboards[color] ^= mask
		p.material.RemovePiece(piece, sq)
		p.hash = p.hash.UpdatePieceOnSq(piece, EMPTY, sq)
	}
	return piece
}

// movePiece moves the piece and keeps the bitboards consistent
// NOTE: this does not affect the frozen fields
func (p *Position) movePiece(startSq, endSq Square) (captured Piece) {
	captured = p.removePiece(endSq)
	piece := p.pieces[startSq]
	if piece != EMPTY {
		startMask := BBWithSquares(startSq)
		endMask := BBWithSquares(endSq)
		color := NewColor(piece.IsWhite())

		p.pieces[startSq] = EMPTY
		p.pieces[endSq] = piece
		p.pieceBitboards[piece] ^= startMask | endMask
		p.pieceBitboards[EMPTY] ^= startMask | endMask
		p.colorBitboards[color] ^= startMask | endMask
		p.hash = p.hash.UpdatePieceOnSq(piece, EMPTY, startSq)
		p.hash = p.hash.UpdatePieceOnSq(EMPTY, piece, endSq)
	}
	return
}

func (p *Position) replacePiece(sq Square, piece Piece) {
	p.removePiece(sq)
	p.addPiece(sq, piece)
}

func (p *Position) addPiece(sq Square, piece Piece) {
	if piece != EMPTY {
		mask := BBWithSquares(sq)
		p.pieces[sq] = piece
		p.pieceBitboards[EMPTY] ^= mask
		p.pieceBitboards[piece] ^= mask
		p.colorBitboards[NewColor(piece.IsWhite())] ^= mask
		p.material.AddPiece(piece, sq)
		p.hash = p.hash.UpdatePieceOnSq(EMPTY, piece, sq)
	}
}

// updateFrozenPos updates a copy of the immutable (irreversible) state in
// Position. Castle rights and en passant square also bleed into Position.hash,
// but are NOT updated here. This is intended to be called before the Position's
// pieces are updated.
func (p *Position) updateFrozenPos(move Move) {
	fp := p.frozenPos.Copy()
	p.frozenPos = fp
	fp.EnPassantSq = NULL_SQ
	fp.Rule50++

	start := move.StartSq()
	end := move.EndSq()
	mt := move.Type()
	piece := p.pieces[start]
	capturedPiece := p.pieces[end]
	isWhite := piece.IsWhite()
	pt := piece.Type()

	if mt == CASTLING || pt == PAWN || capturedPiece != EMPTY {
		fp.Rule50 = 0
	}

	if pt == KING {
		if isWhite {
			if fp.CastleRights[W_CASTLE_KINGSIDE_RIGHT] {
				fp.CastleRights[W_CASTLE_KINGSIDE_RIGHT] = false
			}
			if fp.CastleRights[W_CASTLE_QUEENSIDE_RIGHT] {
				fp.CastleRights[W_CASTLE_QUEENSIDE_RIGHT] = false
			}
		} else {
			if fp.CastleRights[B_CASTLE_KINGSIDE_RIGHT] {
				fp.CastleRights[B_CASTLE_KINGSIDE_RIGHT] = false
			}
			if fp.CastleRights[B_CASTLE_QUEENSIDE_RIGHT] {
				fp.CastleRights[B_CASTLE_QUEENSIDE_RIGHT] = false
			}
		}
	}
	if fp.CastleRights[W_CASTLE_QUEENSIDE_RIGHT] {
		fp.CastleRights[W_CASTLE_QUEENSIDE_RIGHT] = !(start == SQ_A1 || end == SQ_A1)
	}
	if fp.CastleRights[W_CASTLE_KINGSIDE_RIGHT] {
		fp.CastleRights[W_CASTLE_KINGSIDE_RIGHT] = !(start == SQ_H1 || end == SQ_H1)
	}
	if fp.CastleRights[B_CASTLE_QUEENSIDE_RIGHT] {
		fp.CastleRights[B_CASTLE_QUEENSIDE_RIGHT] = !(start == SQ_A8 || end == SQ_A8)
	}
	if fp.CastleRights[B_CASTLE_KINGSIDE_RIGHT] {
		fp.CastleRights[B_CASTLE_KINGSIDE_RIGHT] = !(start == SQ_H8 || end == SQ_H8)
	}
	if pt == PAWN {
		if start.Rank() == 2 && end.Rank() == 4 {
			fp.EnPassantSq = SqFromCoords(3, int(end.File()))
		} else if start.Rank() == 7 && end.Rank() == 5 {
			fp.EnPassantSq = SqFromCoords(6, int(end.File()))
		}
	}
}

func (p *Position) isSquareAttacked(attackColor Color, sq Square) bool {
	pawnAttacksBB := PawnAttacksBB(sq, attackColor.Opp())
	var attackPiece = NewPiece(PAWN, attackColor)
	if pawnAttacksBB&p.pieceBitboards[attackPiece] != 0 {
		return true
	}
	knightAttacksBB := KnightAttacksBB(sq)
	attackPiece = NewPiece(KNIGHT, attackColor)
	if knightAttacksBB&p.pieceBitboards[attackPiece] != 0 {
		return true
	}
	kingAttacksBB := KingAttacksBB(sq)
	attackPiece = NewPiece(KING, attackColor)
	if kingAttacksBB&p.pieceBitboards[attackPiece] != 0 {
		return true
	}
	occupiedBB := p.OccupiedBB()
	diagAttacksBB := SlidingAttacksBB(occupiedBB, sq, BISHOP)
	attackPiece = NewPiece(BISHOP, attackColor)
	if diagAttacksBB&p.pieceBitboards[attackPiece] != 0 {
		return true
	}
	attackPiece = NewPiece(QUEEN, attackColor)
	if diagAttacksBB&p.pieceBitboards[attackPiece] != 0 {
		return true
	}
	straightAttacks := SlidingAttacksBB(occupiedBB, sq, ROOK)
	attackPiece = NewPiece(ROOK, attackColor)
	if straightAttacks&p.pieceBitboards[attackPiece] != 0 {
		return true
	}
	attackPiece = NewPiece(QUEEN, attackColor)
	if straightAttacks&p.pieceBitboards[attackPiece] != 0 {
		return true
	}

	return false
}

func (p *Position) givesCheck(move Move) bool {
	end := move.EndSq()
	piece := p.pieces[move.StartSq()]
	pt := piece.Type()
	color := piece.Color()

	var attacksBB Bitboard
	if pt == PAWN {
		attacksBB = PawnAttacksBB(end, color)
	} else if pt == KNIGHT {
		attacksBB = KnightAttacksBB(end)
	} else if pt == BISHOP || pt == ROOK || pt == QUEEN {
		attacksBB = SlidingAttacksBB(p.OccupiedBB(), end, pt)
	} else {
		attacksBB = KingAttacksBB(end)
	}

	if color == WHITE {
		return p.pieceBitboards[B_KING]&attacksBB > 0
	} else {
		return p.pieceBitboards[W_KING]&attacksBB > 0
	}
}
