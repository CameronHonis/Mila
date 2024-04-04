package main_test

import (
	"github.com/CameronHonis/Mila"
	"github.com/CameronHonis/chess"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SortMoves", func() {
	It("sorts the moves by expected value", func() {
		pos, posErr := chess.BoardFromFEN("1r1q3r/pBP2pbp/1p2p1pn/4P2k/4QP2/B4N1P/P5P1/R4RK1 w - - 1 19")
		Expect(posErr).ToNot(HaveOccurred())
		moves, movesErr := chess.GetLegalMoves(pos)
		Expect(movesErr).ToNot(HaveOccurred())
		sortedMoves := main.SortMoves(pos, moves)
		Expect(sortedMoves).To(HaveLen(len(moves)))
		for i := 0; i < len(moves)-1; i++ {
			prevMove := sortedMoves[i]
			currMove := sortedMoves[i+1]
			prevMoveVal := main.EvalMove(pos, prevMove)
			currMoveVal := main.EvalMove(pos, currMove)
			Expect(prevMoveVal).To(BeNumerically(">=", currMoveVal))
		}
	})
})
