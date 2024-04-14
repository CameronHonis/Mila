package main_test

import (
	"github.com/CameronHonis/Mila"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Square", func() {
	Describe("#SqFromAlg", func() {
		When("the square algebraic notation is invalid", func() {
			It("returns an error", func() {
				Expect(main.SqFromAlg("")).Error().To(HaveOccurred())
				Expect(main.SqFromAlg("a0")).Error().To(HaveOccurred())
				Expect(main.SqFromAlg("j4")).Error().To(HaveOccurred())
			})
		})
		When("the square algebraic notation is valid", func() {
			It("returns a square value", func() {
				Expect(main.SqFromAlg("a1")).To(Equal(main.SQ_A1))
				Expect(main.SqFromAlg("h1")).To(Equal(main.SQ_H1))
				Expect(main.SqFromAlg("d4")).To(Equal(main.SQ_D4))
			})
		})
	})
	Describe("::Rank", func() {
		It("returns the correct rank", func() {
			Expect(main.SQ_A1.Rank()).To(BeEquivalentTo(1))
			Expect(main.SQ_B1.Rank()).To(BeEquivalentTo(1))
			Expect(main.SQ_A2.Rank()).To(BeEquivalentTo(2))
			Expect(main.SQ_H8.Rank()).To(BeEquivalentTo(8))
		})
	})
	Describe("::File", func() {
		It("returns the correct file", func() {
			Expect(main.SQ_A1.File()).To(BeEquivalentTo(1))
			Expect(main.SQ_A2.File()).To(BeEquivalentTo(1))
			Expect(main.SQ_B1.File()).To(BeEquivalentTo(2))
			Expect(main.SQ_H8.File()).To(BeEquivalentTo(8))
		})
	})
	Describe("::PosDiagIdx", func() {
		It("returns the correct pos diag idx", func() {
			Expect(main.SQ_H1.PosDiagIdx()).To(BeEquivalentTo(0))
			Expect(main.SQ_G1.PosDiagIdx()).To(BeEquivalentTo(1))
			Expect(main.SQ_H2.PosDiagIdx()).To(BeEquivalentTo(1))
			Expect(main.SQ_A1.PosDiagIdx()).To(BeEquivalentTo(7))
			Expect(main.SQ_D4.PosDiagIdx()).To(BeEquivalentTo(7))
			Expect(main.SQ_A2.PosDiagIdx()).To(BeEquivalentTo(8))
			Expect(main.SQ_G8.PosDiagIdx()).To(BeEquivalentTo(8))
			Expect(main.SQ_A7.PosDiagIdx()).To(BeEquivalentTo(13))
			Expect(main.SQ_A8.PosDiagIdx()).To(BeEquivalentTo(14))
		})
	})
	Describe("::NegDiagIdx", func() {
		It("returns the correct diag idx", func() {
			Expect(main.SQ_A1.NegDiagIdx()).To(BeEquivalentTo(0))
			Expect(main.SQ_A2.NegDiagIdx()).To(BeEquivalentTo(1))
			Expect(main.SQ_B1.NegDiagIdx()).To(BeEquivalentTo(1))
			Expect(main.SQ_G1.NegDiagIdx()).To(BeEquivalentTo(6))
			Expect(main.SQ_A7.NegDiagIdx()).To(BeEquivalentTo(6))
			Expect(main.SQ_H1.NegDiagIdx()).To(BeEquivalentTo(7))
			Expect(main.SQ_E4.NegDiagIdx()).To(BeEquivalentTo(7))
			Expect(main.SQ_H2.NegDiagIdx()).To(BeEquivalentTo(8))
			Expect(main.SQ_B8.NegDiagIdx()).To(BeEquivalentTo(8))
			Expect(main.SQ_H7.NegDiagIdx()).To(BeEquivalentTo(13))
			Expect(main.SQ_G8.NegDiagIdx()).To(BeEquivalentTo(13))
			Expect(main.SQ_H8.NegDiagIdx()).To(BeEquivalentTo(14))
		})
	})
})
