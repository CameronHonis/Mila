package main_test

import (
	"github.com/CameronHonis/Mila"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Move", func() {
	Describe("#NewMove", func() {
		It("builds a move with the expected bits given the start and end square", func() {
			move := main.NewMove(main.SQ_A1, main.SQ_D4, main.InitState())
			Expect(move).To(BeEquivalentTo(0b000000_011011_0000))
		})
	})
	Describe("::StartSquare", func() {
		It("returns the correct start square", func() {
			move := main.NewMove(main.SQ_D2, main.SQ_D4, main.InitState())
			Expect(move.StartSq()).To(BeEquivalentTo(main.SQ_D2))
		})
	})
	Describe("::EndSquare", func() {
		It("returns the correct end square", func() {
			move := main.NewMove(main.SQ_A1, main.SQ_H6, main.InitState())
			Expect(move.EndSq()).To(BeEquivalentTo(main.SQ_H6))
		})
	})
	Describe("::CapturedPiece", func() {
		It("returns the correct captured piece", func() {
			move := main.NewMove(main.SQ_A1, main.SQ_D8, main.InitState())
			Expect(move.CapturedPiece()).To(BeEquivalentTo(main.B_QUEEN))
		})
	})
})
