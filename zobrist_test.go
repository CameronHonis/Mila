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

func fullSearchToDepth(pos *chess.Board, depth int, fensByHash map[int64]*set.Set[string]) int {
	hash := main.ZobristHash(pos)
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
		nodeCnt += fullSearchToDepth(newPos, depth-1, fensByHash)
	}
	return nodeCnt
}

var _ = PDescribe("Zobrist", func() {
	It("returns unique hashes for all boards depth 4 from init board", func() {
		fensByHash := make(map[int64]*set.Set[string])
		fullSearchToDepth(chess.GetInitBoard(), 4, fensByHash)
		for hash, fens := range fensByHash {
			if fens.Size() > 1 {
				log.Fatalf("%d has collisions, %s", hash, strings.Join(fens.Flatten(), ", "))
			}
		}
	})
})
