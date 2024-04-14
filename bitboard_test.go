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

	Describe("#BBWithHighBitsAt", func() {
		When("the input is only a single high bit", func() {
			It("places the single high bit as expected", func() {
				Expect(main.BBWithHighBitsAt(7)).To(Equal(main.Bitboard(128)))
			})
		})
		When("the input is multiple high bits", func() {
			It("places high bits at all the specified locations", func() {
				highBitIdxs := []int{0, 23, 12, 44, 61}
				var expBitboard = main.Bitboard(0b00100000_00000000_00010000_00000000_00000000_10000000_00010000_00000001)
				Expect(main.BBWithHighBitsAt(highBitIdxs...)).To(Equal(expBitboard))
			})
		})
	})

	Describe("#BBWithRank", func() {
		When("the rank is 1", func() {
			It("returns the expected bitboard", func() {
				bb := main.BBWithRank(1, 0b10101001)
				expBB := main.BBWithSquares(main.SQ_A1, main.SQ_D1, main.SQ_F1, main.SQ_H1)
				Expect(bb).To(Equal(expBB))
			})
		})
		When("the rank is 3", func() {
			It("returns the expected bitboard", func() {
				bb := main.BBWithRank(3, 0b01000001)
				expBB := main.BBWithSquares(main.SQ_A3, main.SQ_G3)
				Expect(bb).To(Equal(expBB))
			})
		})
		When("the rank is 8", func() {
			It("returns the expected bitboard", func() {
				bb := main.BBWithRank(8, 0b10100000)
				expBB := main.BBWithSquares(main.SQ_F8, main.SQ_H8)
				Expect(bb).To(Equal(expBB))
			})
		})
	})

	Describe("#BBWithFile", func() {
		When("the file is 1", func() {
			It("returns the expected bitboard", func() {
				bb := main.BBWithFile(1, 0b10000101)
				expBB := main.BBWithSquares(main.SQ_A1, main.SQ_A3, main.SQ_A8)
				Expect(bb).To(Equal(expBB))
			})
		})
		When("the file is 4", func() {
			It("returns the expected bitboard", func() {
				bb := main.BBWithFile(4, 0b00010100)
				expBB := main.BBWithSquares(main.SQ_D3, main.SQ_D5)
				Expect(bb).To(Equal(expBB))
			})
		})
		When("the file is 8", func() {
			It("returns the expected bitboard", func() {
				bb := main.BBWithFile(8, 0b01010001)
				expBB := main.BBWithSquares(main.SQ_H1, main.SQ_H5, main.SQ_H7)
				Expect(bb).To(Equal(expBB))
			})
		})
	})

	Describe("#BBWithPosDiag", func() {
		When("the diagIdx is 0", func() {
			It("returns the expected bitboard", func() {
				bb := main.BBWithPosDiag(0, 0b1)
				expBB := main.BBWithSquares(main.SQ_H1)
				Expect(bb).To(Equal(expBB))
			})
		})
		When("the diagIdx is 6", func() {
			It("returns the expected bitboard", func() {
				bb := main.BBWithPosDiag(6, 0b1010101)
				expBB := main.BBWithSquares(main.SQ_B1, main.SQ_D3, main.SQ_F5, main.SQ_H7)
				Expect(bb).To(Equal(expBB))
			})
		})
		When("the diagIdx is 7", func() {
			It("returns the expected bitboard", func() {
				bb := main.BBWithPosDiag(7, 0b01010101)
				expBB := main.BBWithSquares(main.SQ_A1, main.SQ_C3, main.SQ_E5, main.SQ_G7)
				Expect(bb).To(Equal(expBB))
			})
		})
		When("the diagIdx is 8", func() {
			It("returns the expected bitboard", func() {
				bb := main.BBWithPosDiag(8, 0b1010101)
				expBB := main.BBWithSquares(main.SQ_A2, main.SQ_C4, main.SQ_E6, main.SQ_G8)
				Expect(bb).To(Equal(expBB))
			})
		})
		When("the diagIdx is 14", func() {
			It("returns the expected bitboard", func() {
				bb := main.BBWithPosDiag(14, 0b1)
				expBB := main.BBWithSquares(main.SQ_A8)
				Expect(bb).To(Equal(expBB))
			})
		})
	})

	Describe("#BBWithNegDiag", func() {
		When("the diagIdx is 0", func() {
			It("returns the expected bitboard", func() {
				bb := main.BBWithNegDiag(0, 0b1)
				expBB := main.BBWithSquares(main.SQ_A1)
				Expect(bb).To(Equal(expBB))
			})
		})
		When("the diagIdx is 6", func() {
			It("returns the expected bitboard", func() {
				bb := main.BBWithNegDiag(6, 0b1010101)
				expBB := main.BBWithSquares(main.SQ_G1, main.SQ_E3, main.SQ_C5, main.SQ_A7)
				Expect(bb).To(Equal(expBB))
			})
		})
		When("the diagIdx is 7", func() {
			It("returns the expected bitboard", func() {
				bb := main.BBWithNegDiag(7, 0b01010101)
				expBB := main.BBWithSquares(main.SQ_H1, main.SQ_F3, main.SQ_D5, main.SQ_B7)
				Expect(bb).To(Equal(expBB))
			})
		})
		When("the diagIdx is 8", func() {
			It("returns the expected bitboard", func() {
				bb := main.BBWithNegDiag(8, 0b1010101)
				expBB := main.BBWithSquares(main.SQ_H2, main.SQ_F4, main.SQ_D6, main.SQ_B8)
				Expect(bb).To(Equal(expBB))
			})
		})
		When("the diagIdx is 14", func() {
			It("returns the expected bitboard", func() {
				bb := main.BBWithNegDiag(14, 0b1)
				expBB := main.BBWithSquares(main.SQ_H8)
				Expect(bb).To(Equal(expBB))
			})
		})
	})
})
