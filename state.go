package main

import (
	"fmt"
	"github.com/CameronHonis/marker"
	"strconv"
	"strings"
)

type State struct {
	__mutableOnMakeMove__ marker.Marker
	Pos                   *Position
	Repetitions           map[ZHash]uint8
	NMoves                uint
	Result                Result

	__immutableOnMakeMove__ marker.Marker
	EnPassantSq             Square
	Rule50                  uint
	CastleRights            [N_CASTLE_RIGHTS]bool
}

func InitState() *State {
	return &State{
		Pos:          InitPos(),
		Repetitions:  make(map[ZHash]uint8),
		Result:       RESULT_IN_PROGRESS,
		CastleRights: [N_CASTLE_RIGHTS]bool{true, true, true, true},
		EnPassantSq:  NULL_SQ,
		NMoves:       1,
	}
}

func StateFromFEN(fen string) (*State, error) {
	state := &State{
		Repetitions: make(map[ZHash]uint8),
	}
	fenSegs := strings.Split(fen, " ")
	if len(fenSegs) != 6 {
		return nil, fmt.Errorf("invalid number of fen segments %d, expected 6", len(fenSegs))
	}

	if pos, posErr := PosFromFEN(fen); posErr != nil {
		return nil, posErr
	} else {
		state.Pos = pos
	}

	castleRightsSpecifier := fenSegs[2]
	if castleRightsSpecifier != "-" {
		for _, castleRightChar := range []byte(castleRightsSpecifier) {
			if castleRightChar == 'K' {
				state.CastleRights[W_CAN_CASTLE_KINGSIDE] = true
			} else if castleRightChar == 'Q' {
				state.CastleRights[W_CAN_CASTLE_QUEENSIDE] = true
			} else if castleRightChar == 'k' {
				state.CastleRights[B_CAN_CASTLE_KINGSIDE] = true
			} else if castleRightChar == 'q' {
				state.CastleRights[B_CAN_CASTLE_QUEENSIDE] = true
			} else {
				return nil, fmt.Errorf("could not parse castle rights, unknown char %c in %s", castleRightChar, castleRightsSpecifier)
			}
		}
	}

	epSpecifier := fenSegs[3]
	if epSpecifier == "-" {
		state.EnPassantSq = NULL_SQ
	} else {
		epSq, epSqErr := SqFromAlg(epSpecifier)
		if epSqErr != nil {
			return nil, fmt.Errorf("could not parse en passant square in fen %s: %s", fen, epSqErr)
		}
		state.EnPassantSq = epSq
	}

	halfmovesSpecifier := fenSegs[4]
	if halfmoves, halfmoveErr := strconv.Atoi(halfmovesSpecifier); halfmoveErr != nil {
		return nil, fmt.Errorf("could not parse halfmoves: %s", halfmoveErr)
	} else {
		state.Rule50 = uint(halfmoves)
	}

	movesSpecifier := fenSegs[5]
	if nMoves, nMovesErr := strconv.Atoi(movesSpecifier); nMovesErr != nil {
		return nil, fmt.Errorf("could not parse number of moves: %s", nMovesErr)
	} else {
		state.NMoves = uint(nMoves)
	}

	return state, nil
}

func (s *State) FEN() string {
	var rtnBuilder strings.Builder
	rtnBuilder.WriteString(s.Pos.FEN())
	rtnBuilder.WriteByte(' ')

	var anyCastleRights bool
	if s.CastleRights[W_CAN_CASTLE_KINGSIDE] {
		rtnBuilder.WriteByte('K')
		anyCastleRights = true
	}
	if s.CastleRights[W_CAN_CASTLE_QUEENSIDE] {
		rtnBuilder.WriteByte('Q')
		anyCastleRights = true
	}
	if s.CastleRights[B_CAN_CASTLE_KINGSIDE] {
		rtnBuilder.WriteByte('k')
		anyCastleRights = true
	}
	if s.CastleRights[B_CAN_CASTLE_QUEENSIDE] {
		rtnBuilder.WriteByte('q')
		anyCastleRights = true
	}
	if !anyCastleRights {
		rtnBuilder.WriteByte('-')
	}
	rtnBuilder.WriteByte(' ')

	if s.EnPassantSq.IsNull() {
		rtnBuilder.WriteByte('-')
	} else {
		rtnBuilder.WriteString(s.EnPassantSq.String())
	}
	rtnBuilder.WriteByte(' ')

	rtnBuilder.WriteString(strconv.Itoa(int(s.Rule50)))
	rtnBuilder.WriteByte(' ')

	rtnBuilder.WriteString(strconv.Itoa(int(s.NMoves)))
	return rtnBuilder.String()
}

//func (s *State) LegalMoves() []Move {
//
//}
