package main_test

import (
	"github.com/CameronHonis/Mila"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("PlyFromNMoves", func() {
	It("converts different moves to their corresponding ply", func() {
		Expect(main.PlyFromNMoves(1, true)).To(Equal(main.Ply(0)))
		Expect(main.PlyFromNMoves(1, false)).To(Equal(main.Ply(1)))
		Expect(main.PlyFromNMoves(2, true)).To(Equal(main.Ply(2)))
		Expect(main.PlyFromNMoves(2, false)).To(Equal(main.Ply(3)))
	})
})

var _ = Describe("NMovesFromPly", func() {
	It("converts different plies to their corresponding move", func() {
		Expect(main.NMovesFromPly(0)).To(Equal(uint(1)))
		Expect(main.NMovesFromPly(1)).To(Equal(uint(1)))
		Expect(main.NMovesFromPly(2)).To(Equal(uint(2)))
		Expect(main.NMovesFromPly(3)).To(Equal(uint(2)))
		Expect(main.NMovesFromPly(4)).To(Equal(uint(3)))
	})
})
