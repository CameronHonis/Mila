package main_test

import (
	"github.com/CameronHonis/Mila"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Position", func() {
	Describe("#FromFEN", func() {

	})
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
	Describe("::ToFEN", func() {

	})
})
