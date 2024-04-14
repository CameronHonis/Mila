package main

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Precompute", func() {
	Describe("#SlidingAttacksBB", func() {
		When("the piece is a rook", func() {
			When("no obstructing pieces are on the board", func() {
				It("returns all possible attacks", func() {
					occupiedBB := BBWithSquares(SQ_D4)
					expAttackBB := BBWithSquares(SQ_D1, SQ_D2, SQ_D3, SQ_D5, SQ_D6, SQ_D7, SQ_D8,
						SQ_A4, SQ_B4, SQ_C4, SQ_E4, SQ_F4, SQ_G4, SQ_H4)
					attacksBB := SlidingAttacksBB(occupiedBB, SQ_D4, ROOK)
					Expect(attacksBB).To(Equal(expAttackBB))
				})
			})
			When("one obstructing piece per direction exists", func() {
				It("returns all possible attacks", func() {
					occupiedBB := BBWithSquares(SQ_D4, SQ_B4, SQ_F4, SQ_D3, SQ_D8)
					expAttackBB := BBWithSquares(SQ_B4, SQ_C4, SQ_E4, SQ_F4, SQ_D3, SQ_D5, SQ_D6, SQ_D7, SQ_D8)
					attacksBB := SlidingAttacksBB(occupiedBB, SQ_D4, ROOK)
					Expect(attacksBB).To(Equal(expAttackBB))
				})
			})
			When("two obstructing pieces per direction exists", func() {
				It("returns all possible attacks", func() {
					occupiedBB := BBWithSquares(SQ_D4, SQ_A4, SQ_B4, SQ_F4, SQ_H4, SQ_D3, SQ_D1, SQ_D7, SQ_D8)
					expAttackBB := BBWithSquares(SQ_B4, SQ_C4, SQ_D3, SQ_E4, SQ_F4, SQ_D5, SQ_D6, SQ_D7)
					attacksBB := SlidingAttacksBB(occupiedBB, SQ_D4, ROOK)
					Expect(attacksBB).To(Equal(expAttackBB))
				})
			})
		})
		When("the piece is a queen", func() {
			When("no obstructing pieces are on the board", func() {
				It("returns all possible attacks", func() {
					occupiedBB := BBWithSquares(SQ_D4)
					expAttackBB := BBWithSquares(SQ_D1, SQ_D2, SQ_D3, SQ_D5, SQ_D6, SQ_D7, SQ_D8,
						SQ_A4, SQ_B4, SQ_C4, SQ_E4, SQ_F4, SQ_G4, SQ_H4,
						SQ_A1, SQ_B2, SQ_C3, SQ_E5, SQ_F6, SQ_G7, SQ_H8,
						SQ_A7, SQ_B6, SQ_C5, SQ_E3, SQ_F2, SQ_G1)
					attacksBB := SlidingAttacksBB(occupiedBB, SQ_D4, QUEEN)
					Expect(attacksBB).To(Equal(expAttackBB))
				})
			})
			When("one obstructing piece per direction exists", func() {
				It("returns all possible attacks", func() {
					occupiedBB := BBWithSquares(SQ_D4, SQ_F4, SQ_H8, SQ_D5, SQ_B6, SQ_B4, SQ_C3, SQ_D2, SQ_F2)
					expAttackBB := BBWithSquares(SQ_E4, SQ_F4, SQ_E5, SQ_F6, SQ_G7, SQ_H8, SQ_D5, SQ_C5,
						SQ_B6, SQ_C4, SQ_B4, SQ_C3, SQ_D3, SQ_D2, SQ_E3, SQ_F2)
					attacksBB := SlidingAttacksBB(occupiedBB, SQ_D4, QUEEN)
					Expect(attacksBB).To(Equal(expAttackBB))
				})
			})
			When("two obstructing pieces per direction exists", func() {
				It("returns all possible attacks", func() {
					occupiedBB := BBWithSquares(SQ_D4, SQ_F4, SQ_H4, SQ_G7, SQ_H8, SQ_D5, SQ_D8, SQ_B6,
						SQ_A7, SQ_B4, SQ_A4, SQ_C3, SQ_B2, SQ_D3, SQ_D1, SQ_F2, SQ_G1)
					expAttackBB := BBWithSquares(SQ_E4, SQ_F4, SQ_E5, SQ_F6, SQ_G7, SQ_D5, SQ_C5, SQ_B6, SQ_C4,
						SQ_B4, SQ_C3, SQ_D3, SQ_E3, SQ_F2)
					attacksBB := SlidingAttacksBB(occupiedBB, SQ_D4, QUEEN)
					Expect(attacksBB).To(Equal(expAttackBB))
				})
			})
		})
		When("the piece is a bishop", func() {
			When("no obstructing pieces are on the board", func() {
				It("returns all possible attacks", func() {
					occupiedBB := BBWithSquares(SQ_D4)
					expAttackBB := BBWithSquares(SQ_E5, SQ_F6, SQ_G7, SQ_H8, SQ_C5, SQ_B6, SQ_A7, SQ_C3, SQ_B2,
						SQ_A1, SQ_E3, SQ_F2, SQ_G1)
					attacksBB := SlidingAttacksBB(occupiedBB, SQ_D4, BISHOP)
					Expect(attacksBB).To(Equal(expAttackBB))
				})
			})
			When("one obstructing piece per direction exists", func() {
				It("returns all possible attacks", func() {
					occupiedBB := BBWithSquares(SQ_D4, SQ_H8, SQ_A7, SQ_B2, SQ_E3)
					expAttackBB := BBWithSquares(SQ_E5, SQ_F6, SQ_G7, SQ_H8, SQ_C5, SQ_B6, SQ_A7, SQ_C3, SQ_B2, SQ_E3)
					attacksBB := SlidingAttacksBB(occupiedBB, SQ_D4, BISHOP)
					Expect(attacksBB).To(Equal(expAttackBB))
				})
			})
			When("two obstructing pieces per direction exists", func() {
				It("returns all possible attacks", func() {
					occupiedBB := BBWithSquares(SQ_E5, SQ_H8, SQ_C5, SQ_B6, SQ_B2, SQ_A1, SQ_E3, SQ_G1)
					expAttackBB := BBWithSquares(SQ_E5, SQ_C5, SQ_C3, SQ_B2, SQ_E3)
					attacksBB := SlidingAttacksBB(occupiedBB, SQ_D4, BISHOP)
					Expect(attacksBB).To(Equal(expAttackBB))
				})
			})
		})
	})

	Describe("#genRankOccupiedBBs", func() {
		It("generates the appropriate number of bitboards", func() {
			bbs := genRankOccupiedBBs()
			Expect(bbs).To(HaveLen(8 * 255))
		})
	})
	Describe("#genFileOccupiedBBs", func() {
		It("generates the appropriate number of bitboards", func() {
			bbs := genFileOccupiedBBs()
			Expect(bbs).To(HaveLen(8 * 255))
		})
	})
	Describe("#genPosDiagOccupiedBBs", func() {
		It("generates the appropriate number of bitboards", func() {
			bbs := genPosDiagOccupiedBBs()
			var expNBBs int
			expNBBs += 2 * ((1 << 1) - 1) // corners
			expNBBs += 2 * ((1 << 2) - 1) // diags 1 & 13
			expNBBs += 2 * ((1 << 3) - 1) // diags 2 & 12
			expNBBs += 2 * ((1 << 4) - 1) // diags 3 & 11
			expNBBs += 2 * ((1 << 5) - 1) // diags 4 & 10
			expNBBs += 2 * ((1 << 6) - 1) // diags 5 & 9
			expNBBs += 2 * ((1 << 7) - 1) // diags 6 & 8
			expNBBs += (1 << 8) - 1       // diag 7
			Expect(bbs).To(HaveLen(expNBBs))
		})
	})
	Describe("#genNegDiagOccupiedBBs", func() {
		It("generates the appropriate number of bitboards", func() {
			bbs := genNegDiagOccupiedBBs()
			var expNBBs int
			expNBBs += 2 * ((1 << 1) - 1) // corners
			expNBBs += 2 * ((1 << 2) - 1) // diags 1 & 13
			expNBBs += 2 * ((1 << 3) - 1) // diags 2 & 12
			expNBBs += 2 * ((1 << 4) - 1) // diags 3 & 11
			expNBBs += 2 * ((1 << 5) - 1) // diags 4 & 10
			expNBBs += 2 * ((1 << 6) - 1) // diags 5 & 9
			expNBBs += 2 * ((1 << 7) - 1) // diags 6 & 8
			expNBBs += (1 << 8) - 1       // diag 7
			Expect(bbs).To(HaveLen(expNBBs))
		})
	})

})
