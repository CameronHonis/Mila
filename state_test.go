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
	Describe("::ToFEN", func() {

	})
})

var _ = Describe("Position", func() {
	Describe("::String", func() {
		It("represents the init board clearly", func() {
			pos := main.InitPos()
			expStr := "" +
				"8 ♖ ♘ ♗ ♕ ♔ ♗ ♘ ♖ \n" +
				"7 ♙ ♙ ♙ ♙ ♙ ♙ ♙ ♙ \n" +
				"6   ░░  ░░  ░░  ░░\n" +
				"5 ░░  ░░  ░░  ░░  \n" +
				"4   ░░  ░░  ░░  ░░\n" +
				"3 ░░  ░░  ░░  ░░  \n" +
				"2 ♟︎ ♟︎ ♟︎ ♟︎ ♟︎ ♟︎ ♟︎ ♟︎ \n" +
				"1 ♜ ♞ ♝ ♛ ♚ ♝ ♞ ♜ \n" +
				"  1 2 3 4 5 6 7 8 "
			Expect(pos.String()).To(Equal(expStr))
		})
	})
	Describe("#PositionFromFEN", func() {
		When("the FEN is invalid", func() {
			When("the FEN does not have enough segments", func() {
				It("returns an error", func() {
					fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0"
					Expect(main.PosFromFEN(fen)).Error().To(HaveOccurred())
				})
			})
			When("the FEN has too many segments", func() {
				It("returns an error", func() {
					fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1 1"
					Expect(main.PosFromFEN(fen)).Error().To(HaveOccurred())
				})
			})
			When("the FEN pieces contain too many (9) rows", func() {
				It("returns an error", func() {
					fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR/8 w KQkq - 0 1"
					Expect(main.PosFromFEN(fen)).Error().To(HaveOccurred())
				})
			})
			When("the FEN pieces contain too few (7) rows", func() {
				It("returns an error", func() {
					fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP w KQkq - 0 1"
					Expect(main.PosFromFEN(fen)).Error().To(HaveOccurred())
				})
			})
			When("a row in the FEN pieces contain too many (9) files", func() {
				It("returns an error", func() {
					fen := "rnbqkbnrr/ppppppppp/9/9/9/9/PPPPPPPPP/RNBQKBNRR w KQkq - 0 1"
					Expect(main.PosFromFEN(fen)).Error().To(HaveOccurred())
				})
			})
			When("a row in the FEN pieces contain too few (7) files", func() {
				It("returns an error", func() {
					fen := "rnbqkbn/ppppppp/7/7/7/7/PPPPPPP/RNBQKBN w KQkq - 0 1"
					Expect(main.PosFromFEN(fen)).Error().To(HaveOccurred())
				})
			})
		})
		When("the FEN is valid", func() {
			It("returns a Position with the pieces as described in the FEN", func() {
				fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
				pos, posErr := main.PosFromFEN(fen)
				Expect(posErr).To(Succeed())
				Expect(pos).To(Equal(main.InitPos()))
			})
		})
	})
	Describe("::ToFEN", func() {
		It("serializes the position into the FEN format", func() {
			fen := main.InitPos().FEN()
			Expect(fen).To(Equal("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w"))
		})
	})
})
