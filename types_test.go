package main_test

import (
	"github.com/CameronHonis/Mila"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Types", func() {
	Describe("SqFromAlg", func() {
		When("the square algebraic notation is invalid", func() {
			It("returns an error", func() {
				Expect(main.SqFromAlg("")).Error().To(HaveOccurred())
				Expect(main.SqFromAlg("a0")).Error().To(HaveOccurred())
				Expect(main.SqFromAlg("j4")).Error().To(HaveOccurred())
			})
		})
		When("the square algebraic notation is valid", func() {
			It("returns a square value", func() {
				Expect(main.SqFromAlg("a1")).To(Equal(main.Square(main.SQ_A1)))
				Expect(main.SqFromAlg("h1")).To(Equal(main.Square(main.SQ_H1)))
				Expect(main.SqFromAlg("d4")).To(Equal(main.Square(main.SQ_D4)))
			})
		})
	})
})
