package main_test

import (
	"github.com/CameronHonis/Mila"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Move", func() {
	Describe("#NewMove", func() {
		It("builds a move with the expected bits given the start and end square", func() {
			move := main.NewMove(main.SQ_A1, main.SQ_D4, main.NULL_SQ, main.EMPTY_PIECE_TYPE, false)
			Expect(move).To(BeEquivalentTo(0b000000_011011_0000))
		})
	})
	Describe("::StartSquare", func() {
		It("returns the correct start square", func() {
			move := main.NewMove(main.SQ_D2, main.SQ_D4, main.NULL_SQ, main.EMPTY_PIECE_TYPE, false)
			Expect(move.StartSq()).To(BeEquivalentTo(main.SQ_D2))
		})
	})
	Describe("::EndSquare", func() {
		It("returns the correct end square", func() {
			move := main.NewMove(main.SQ_A1, main.SQ_H6, main.NULL_SQ, main.EMPTY_PIECE_TYPE, false)
			Expect(move.EndSq()).To(BeEquivalentTo(main.SQ_H6))
		})
	})
	Describe("::PromotedTo", func() {
		It("returns the correct captured piece", func() {
			move := main.NewMove(main.SQ_B7, main.SQ_B8, main.NULL_SQ, main.ROOK, false)
			Expect(move.PromotedTo()).To(BeEquivalentTo(main.ROOK))
		})
	})
	Describe("::Type", func() {
		When("the move is a normal move", func() {
			It("returns the normal move type", func() {
				move := main.NewMove(main.SQ_A1, main.SQ_H6, main.NULL_SQ, main.EMPTY_PIECE_TYPE, false)
				Expect(move.Type()).To(Equal(main.NORMAL_MOVE))
			})
		})
		When("the move is captures en passant", func() {
			It("returns the captures en passant type", func() {
				move := main.NewMove(main.SQ_D5, main.SQ_E6, main.SQ_E6, main.EMPTY_PIECE_TYPE, false)
				Expect(move.Type()).To(Equal(main.CAPTURES_EN_PASSANT))
			})
		})
		When("the move is a pawn promotion", func() {
			It("returns the pawn promotion type", func() {
				move := main.NewMove(main.SQ_B7, main.SQ_B8, main.NULL_SQ, main.BISHOP, false)
				Expect(move.Type()).To(Equal(main.PAWN_PROMOTION))
			})
		})
		When("the move is castles", func() {
			It("returns the castles type", func() {
				move := main.NewMove(main.SQ_E1, main.SQ_G1, main.NULL_SQ, main.EMPTY_PIECE_TYPE, true)
				Expect(move.Type()).To(Equal(main.CASTLING))
			})
		})
	})
})
