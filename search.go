package main

import (
	"fmt"
	"github.com/CameronHonis/chess"
	"math"
	"time"
)

const ALPHA_BETA_PRUNING_ENABLED = false

var searchHalt = false

func StartSearch() {
	searchHalt = false
	maxSearchMs := Options.MaxSearchMs()
	go func() {
		time.Sleep(time.Duration(maxSearchMs) * time.Millisecond)
		searchHalt = true
	}()
	var depth int
	var bestMove *chess.Move
	var nodeCount int
	for {
		depth++

		if searchHalt {
			fmt.Println("search halt")
			break
		}
		if Options.MaxDepth > 0 && depth > Options.MaxDepth {
			fmt.Println("max depth reached")
			break
		}

		var score float64
		bestMove, score, nodeCount = searchToDepth(Position, depth, -math.MaxFloat64, math.MaxFloat64)
		fmt.Printf("info depth %d score %f move %s nodes %d\n", depth, score, bestMove.ToLongAlgebraic(), nodeCount)
	}
	fmt.Printf("bestmove %s\n", bestMove.ToLongAlgebraic())

}

func searchToDepth(pos *chess.Board, depth int, alpha float64, beta float64) (bestMove *chess.Move, score float64, nodeCount int) {
	if depth == 0 || pos.Result != chess.BOARD_RESULT_IN_PROGRESS {
		return nil, EvalPos(pos), 1
	}

	moves, err := chess.GetLegalMoves(pos)
	if err != nil || len(moves) == 0 {
		panic(fmt.Sprintf("could not get legal moves from pos %s: %s", pos, err))
	}
	var bestScore float64
	for _, move := range moves {
		if searchHalt {
			break
		}
		newPos := chess.GetBoardFromMove(pos, move)
		_, subScore, subNodeCount := searchToDepth(newPos, depth-1, alpha, beta)
		nodeCount += subNodeCount
		if Options.MaxNodes > 0 && nodeCount > Options.MaxNodes {
			searchHalt = true
		}
		if pos.IsWhiteTurn {
			if subScore > alpha {
				alpha = subScore
				bestScore = subScore
				bestMove = move
				if ALPHA_BETA_PRUNING_ENABLED && subScore > beta {
					break
				}
			}
		} else { // black turn
			if subScore < beta {
				beta = subScore
				bestScore = subScore
				bestMove = move
				if ALPHA_BETA_PRUNING_ENABLED && subScore < alpha {
					break
				}
			}
		}
	}

	return bestMove, bestScore, nodeCount
}

func msForSearch(pos *chess.Board, bankMs int, incrMs int) int {
	expMoves := ExpMoves(pos)
	return incrMs + bankMs/expMoves
}
