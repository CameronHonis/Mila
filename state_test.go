package main_test

import (
	"github.com/CameronHonis/Mila"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("State", func() {
	Describe("#StateFromFEN", func() {
		When("the FEN is not valid", func() {
			When("the FEN contains an invalid turn specifier", func() {
				It("returns an error", func() {
					fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR x KQkq - 0 1"
					Expect(main.StateFromFEN(fen)).Error().To(HaveOccurred())
				})
			})
			When("the FEN contains an invalid en passant square", func() {
				It("returns an error", func() {
					fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq a9 0 1"
					Expect(main.StateFromFEN(fen)).Error().To(HaveOccurred())
				})
			})
			When("the FEN contains an invalid castle rights char (X)", func() {
				It("returns an error", func() {
					fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w X a9 0 1"
					Expect(main.StateFromFEN(fen)).Error().To(HaveOccurred())
				})
			})
		})
		When("the FEN is valid", func() {
			It("returns a State with all fields populated based on the FEN", func() {
				fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq b3 4 12"
				state, stateErr := main.StateFromFEN(fen)
				Expect(stateErr).To(Succeed())
				Expect(state.Pos).ToNot(BeNil())
				Expect(state.Repetitions).ToNot(BeNil())
				Expect(state.NMoves).To(BeEquivalentTo(12))
				Expect(state.Result).To(Equal(main.RESULT_IN_PROGRESS))
				Expect(state.EnPassantSq).To(Equal(main.SQ_B3))
				Expect(state.Rule50).To(BeEquivalentTo(4))
				expCastleRights := [4]bool{true, true, true, true}
				Expect(state.CastleRights).To(Equal(expCastleRights))
			})
			When("no players have castle rights", func() {
				It("returns a State with no castle rights", func() {
					fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w - - 0 1"
					state, stateErr := main.StateFromFEN(fen)
					Expect(stateErr).To(Succeed())
					expCastleRights := [4]bool{false, false, false, false}
					Expect(state.CastleRights).To(Equal(expCastleRights))
				})
			})
			When("only white has castle rights", func() {
				It("returns a State with only white has castle rights", func() {
					fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQ - 0 1"
					state, stateErr := main.StateFromFEN(fen)
					Expect(stateErr).To(Succeed())
					expCastleRights := [4]bool{}
					expCastleRights[main.W_CAN_CASTLE_KINGSIDE] = true
					expCastleRights[main.W_CAN_CASTLE_QUEENSIDE] = true
					Expect(state.CastleRights).To(Equal(expCastleRights))
				})
			})
			When("the en passant square is null", func() {
				It("returns a State with no en passant square", func() {
					fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
					state, stateErr := main.StateFromFEN(fen)
					Expect(stateErr).To(Succeed())
					Expect(state.EnPassantSq).To(Equal(main.NULL_SQ))
				})
			})
		})
	})
	Describe("::FEN", func() {
		It("serializes the position into the FEN format", func() {
			expFen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
			Expect(main.InitState().FEN()).To(Equal(expFen))
		})
	})
})
