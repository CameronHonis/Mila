package main_test

import (
	"github.com/CameronHonis/Mila"
	"github.com/CameronHonis/chess"
	"github.com/CameronHonis/set"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"log"
	"strings"
)

func legacyFullSearchToDepth(pos *chess.Board, depth int, fensByHash map[uint64]*set.Set[string]) int {
	hash := main.ZobristHashOnLegacyBoard(pos)
	trimmedFen := strings.Join(strings.Split(pos.ToFEN(), " ")[:4], " ")
	if _, ok := fensByHash[hash]; !ok {
		fensByHash[hash] = set.EmptySet[string]()
	}
	fensByHash[hash].Add(trimmedFen)
	if depth == 0 {
		return 1
	}
	moves, movesErr := chess.GetLegalMoves(pos)
	nodeCnt := 0
	Expect(movesErr).ToNot(HaveOccurred())
	for _, move := range moves {
		newPos := chess.GetBoardFromMove(pos, move)
		nodeCnt += legacyFullSearchToDepth(newPos, depth-1, fensByHash)
	}
	return nodeCnt
}

var _ = Describe("ZobristHashOnLegacyBoard", func() {
	It("returns unique hashes for all boards depth 4 from init board", func() {
		fensByHash := make(map[uint64]*set.Set[string])
		legacyFullSearchToDepth(chess.GetInitBoard(), 3, fensByHash)
		for hash, fens := range fensByHash {
			if fens.Size() > 1 {
				log.Fatalf("%d has collisions, %s", hash, strings.Join(fens.Flatten(), ", "))
			}
		}
	})
})

var _ = Describe("ZobristHash", func() {
	var pos1 *main.Position
	var pos2 *main.Position
	When("two positions differ by pieces", func() {
		BeforeEach(func() {
			var posErr error
			pos1, posErr = main.FromFEN("rnbqk2r/pppp1ppp/3bpn2/8/2B5/4PN2/PPPP1PPP/RNBQK2R w KQkq - 2 4")
			Expect(posErr).ToNot(HaveOccurred())
			pos2, posErr = main.FromFEN("rnbqk2r/pppp1ppp/3bpn2/8/2B5/4PN2/PPPP1pPP/RNBQK2R w KQkq - 2 4")
			Expect(posErr).ToNot(HaveOccurred())
		})
		It("returns a unique hash for each position", func() {
			Expect(main.NewZHash(pos1)).ToNot(Equal(main.NewZHash(pos2)))
		})
	})
	When("two positions differ by en passant square", func() {
		BeforeEach(func() {
			var posErr error
			pos1, posErr = main.FromFEN("rnbqk2r/pppp1ppp/3bpn2/8/2BP4/4PN2/PPP2PPP/RNBQK2R b KQkq d3 0 4")
			Expect(posErr).ToNot(HaveOccurred())
			pos2, posErr = main.FromFEN("rnbqk2r/pppp1ppp/3bpn2/8/2BP4/4PN2/PPP2PPP/RNBQK2R b KQkq - 0 4")
			Expect(posErr).ToNot(HaveOccurred())
		})
		It("returns a unique hash for each position", func() {
			Expect(main.NewZHash(pos1)).ToNot(Equal(main.NewZHash(pos2)))
		})
	})
	When("two positions differ by castle rights", func() {
		BeforeEach(func() {
			var posErr error
			pos1, posErr = main.FromFEN("rnbqk2r/pppp1ppp/3bpn2/8/2BP4/4PN2/PPP2PPP/RNBQK2R b KQkq d3 0 4")
			Expect(posErr).ToNot(HaveOccurred())
			pos2, posErr = main.FromFEN("rnbqk2r/pppp1ppp/3bpn2/8/2BP4/4PN2/PPP2PPP/RNBQK2R b Kkq d3 0 4")
			Expect(posErr).ToNot(HaveOccurred())
		})
		It("returns a unique hash for each position", func() {
			Expect(main.NewZHash(pos1)).ToNot(Equal(main.NewZHash(pos2)))
		})
	})
})
