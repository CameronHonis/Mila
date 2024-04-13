package main_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/CameronHonis/Mila"
)

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
