package main_test

import (
	"github.com/CameronHonis/Mila"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SortMoves", func() {
	var pos *main.Position
	var moves []main.Move
	var anticipated main.Move
	BeforeEach(func() {
		var posErr error
		pos, posErr = main.FromFEN("1r1q3r/pBP2pbp/1p2p1pn/4P2k/4QP2/B4N1P/P5P1/R4RK1 w - - 1 19")
		Expect(posErr).ToNot(HaveOccurred())
		moves = make([]main.Move, 0)
		iter := main.NewLegalMoveIter(pos)
		for {
			move, done := iter.Next()
			if done {
				break
			}
			moves = append(moves, move)
		}
	})
	When("there is not anticipated move", func() {
		BeforeEach(func() {
			anticipated = main.NULL_MOVE
		})
		It("sorts the moves by expected value", func() {
			sortedMoves := main.SortMoves(pos, moves, anticipated)
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
	When("there is an anticipated move", func() {
		BeforeEach(func() {
			anticipated = main.NewNormalMove(main.SQ_A1, main.SQ_B1)
		})
		It("sorts the moves by expected value with the anticipated move first", func() {
			sortedMoves := main.SortMoves(pos, moves, anticipated)
			Expect(sortedMoves).To(HaveLen(len(moves)))
			Expect(sortedMoves[0]).To(Equal(anticipated))
			for i := 1; i < len(moves)-1; i++ {
				prevMove := sortedMoves[i]
				currMove := sortedMoves[i+1]
				prevMoveVal := main.EvalMove(pos, prevMove)
				currMoveVal := main.EvalMove(pos, currMove)
				Expect(prevMoveVal).To(BeNumerically(">=", currMoveVal))
			}
		})
	})
})
