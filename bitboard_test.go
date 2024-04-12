package main_test

import (
	"fmt"
	main "github.com/CameronHonis/Mila"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bitboard", func() {
	Describe("::String", func() {
		It("Generates a multi-line string of a board that represents the bitboard", func() {
			var bitb main.Bitboard = 0b00000000_00100000_00001000_00000100_00000100_00001000_00100001_00010010
			expStr := "" +
				"8   ░░  ░░  ░░  ░░\n" +
				"7 ░░  ░░  ░░██░░  \n" +
				"6   ░░  ██  ░░  ░░\n" +
				"5 ░░  ██  ░░  ░░  \n" +
				"4   ░░██░░  ░░  ░░\n" +
				"3 ░░  ░░██░░  ░░  \n" +
				"2 ██░░  ░░  ██  ░░\n" +
				"1 ░░██░░  ██  ░░  \n" +
				"  1 2 3 4 5 6 7 8 "
			Expect(bitb.String()).To(Equal(expStr), fmt.Sprintf("%s\nis not\n%s", bitb.String(), expStr))
		})
	})

	Describe("#WithHighBitsAt", func() {
		When("the input is only a single high bit", func() {
			It("places the single high bit as expected", func() {
				Expect(main.WithHighBitsAt(7)).To(Equal(main.Bitboard(128)))
			})
		})
		When("the input is multiple high bits", func() {
			It("places high bits at all the specified locations", func() {
				highBitIdxs := []int{0, 23, 12, 44, 61}
				var expBitboard = main.Bitboard(0b00100000_00000000_00010000_00000000_00000000_10000000_00010000_00000001)
				Expect(main.WithHighBitsAt(highBitIdxs...)).To(Equal(expBitboard))
			})
		})
	})
})
