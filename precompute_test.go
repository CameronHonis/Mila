package main

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Precompute", func() {
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
