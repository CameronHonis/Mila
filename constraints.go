package main

import (
	"math"
)

type SearchConstraints struct {
	moves       []Move
	whiteMs     int
	blackMs     int
	whiteIncrMs int
	blackIncrMs int
	maxDepth    uint8
	maxNodes    int
	maxMs       int
}

func (sc *SearchConstraints) NodeCntLmt() int {
	if sc.maxNodes > 0 {
		return sc.maxNodes
	}
	return math.MaxInt
}

func (sc *SearchConstraints) DepthLmt() uint8 {
	if sc.maxDepth > 0 {
		return sc.maxDepth
	}
	return math.MaxUint8
}

func (sc *SearchConstraints) MaxSearchMs() int {
	var maxSearchMs = 100_000_000_000
	if sc.maxMs > 0 {
		maxSearchMs = MinInt(maxSearchMs, sc.maxMs)
	}
	return maxSearchMs
}
