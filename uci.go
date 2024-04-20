package main

import (
	"bufio"
	"fmt"
	"github.com/CameronHonis/chess"
	"os"
	"strconv"
	"strings"
)

type Uci struct {
	pos *Position
	tt  *TranspTable
}

func NewUci(tt *TranspTable) *Uci {
	return &Uci{
		pos: InitPos(),
		tt:  tt,
	}
}

func (uci *Uci) Start() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		uci.handleInput(line)
	}
	if err := scanner.Err(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "reading stdin: ", err)
	}
}

func (uci *Uci) handleInput(s string) {
	toks := strings.Split(s, " ")
	cmd := toks[0]
	if cmd == "uci" {
		fmt.Println("uciok")
	} else if cmd == "position" {
		pos, err := handlePositionCmd(toks)
		if err != nil {
			fmt.Println(err)
		} else if pos != nil {
			fmt.Println("set position to", pos.FEN())
			uci.pos = pos
		}
	} else if cmd == "go" {
		constraints, err := handleGoCmd(toks, uci.pos)
		if err != nil {
			fmt.Println(err)
		} else if constraints != nil {
			fmt.Println("starting search")
			go NewSearch(uci.pos, constraints, uci.tt).Start()
		}
	} else if cmd == "isready" {
		fmt.Println("readyok")
	} else {
		fmt.Println("unknown command:", cmd)
	}
}

func handlePositionCmd(toks []string) (*Position, error) {
	if len(toks) < 2 || toks[1] == "--help" || toks[1] == "help" {
		printPositionCmdHelp()
		return nil, nil
	}
	tokIdx := 1
	var pos *Position
	if toks[tokIdx] == "fen" {
		if len(toks) < 8 {
			return nil, fmt.Errorf("incorrect number of FEN segments, got %d, expected 6", 8-len(toks))
		}
		fen := strings.Join(toks[2:8], " ")
		var err error
		pos, err = FromFEN(fen)
		if err != nil {
			return nil, fmt.Errorf("could not parse FEN %s: %s", fen, err)
		}
		tokIdx = 8
	} else if toks[tokIdx] == "startpos" {
		pos = InitPos()
		tokIdx = 2
	}

	var isMoveToks = false
	for ; tokIdx < len(toks); tokIdx++ {
		if toks[tokIdx] == "moves" {
			isMoveToks = true
			continue
		}
		if isMoveToks {
			board, boardErr := chess.BoardFromFEN(pos.FEN())
			if boardErr != nil {
				return nil, fmt.Errorf("could not convert from pos to board: %s", boardErr)
			}
			move, moveErr := chess.MoveFromAlgebraic(toks[tokIdx], board)
			if moveErr != nil {
				return nil, fmt.Errorf("could not parse move: %s", moveErr)
			}

			board = chess.GetBoardFromMove(board, move)
			var posErr error
			pos, posErr = FromFEN(board.ToFEN())
			if posErr != nil {
				return nil, fmt.Errorf("could not convert from board to pos: %s", posErr)
			}
		}
	}
	return pos, nil
}

func printPositionCmdHelp() {
	fmt.Println("UCI position: set up the position on the internal board")
	fmt.Println("Usage:")
	fmt.Println("    position [fen | startpos] [moves] ...")
	fmt.Println(Bold("position fen") + " {FEN} ...")
	fmt.Println("    set up the position from the given fen")
	fmt.Println("")
	fmt.Println(Bold("position startpos") + " ...")
	fmt.Println("    set up the position from the start pos")
	fmt.Println("")
	fmt.Println("position ... " + Bold("moves") + " {move1} {move2} ...")
	fmt.Println("   an optional argument that describes the moves that follow the given board position.")
	fmt.Println("   moves must be formatted in UCI-compliant 'Long-Algebraic' format, for more")
	fmt.Println("   information on this format visit: ")
	fmt.Println("   https://www.chessprogramming.org/Algebraic_Chess_Notation#LAN")
}

func handleGoCmd(toks []string, pos *Position) (*SearchConstraints, error) {
	if len(toks) == 2 && (toks[1] == "--help" || toks[1] == "help") {
		printGoCmdHelp()
	}

	opts := &SearchConstraints{}

	var tokIdx = 1
	for tokIdx < len(toks) {
		currTok := toks[tokIdx]
		if currTok == "searchmoves" {
			var moveIdx = 0
			for ; moveIdx+tokIdx+1 < len(toks); moveIdx++ {
				moveStr := toks[moveIdx+tokIdx+1]
				board, boardErr := chess.BoardFromFEN(pos.FEN())
				if boardErr != nil {
					return nil, fmt.Errorf("could not convert from pos to board: %s", boardErr)
				}
				legacyMove, legacyMoveErr := chess.MoveFromLongAlgebraic(moveStr, board)
				if legacyMoveErr != nil {
					return nil, fmt.Errorf("could not parse move %s: %s", moveStr, legacyMoveErr)
				}
				if opts.moves == nil {
					opts.moves = make([]Move, 0)
				}
				move, moveErr := LegalMoveFromLegacyMove(legacyMove, pos)
				if moveErr != nil {
					return nil, fmt.Errorf("could not find legal move match from legacy move: %s", moveErr)
				}
				opts.moves = append(opts.moves, move)
			}
			tokIdx += moveIdx + 1
		} else if currTok == "wtime" {
			if tokIdx+1 >= len(toks) {
				return nil, fmt.Errorf("missing argument for wtime")
			}
			wtimeStr := toks[tokIdx+1]
			wtime, parseErr := strconv.Atoi(wtimeStr)
			if parseErr != nil {
				return nil, fmt.Errorf("could not parse %s as wtime: %s", wtimeStr, parseErr)
			}
			opts.whiteMs = wtime
			tokIdx += 2
		} else if currTok == "btime" {
			if tokIdx+1 >= len(toks) {
				return nil, fmt.Errorf("missing argument for btime")
			}
			btimeStr := toks[tokIdx+1]
			btime, parseErr := strconv.Atoi(btimeStr)
			if parseErr != nil {
				return nil, fmt.Errorf("could not parse %s as btime: %s", btimeStr, parseErr)
			}
			opts.blackMs = btime
			tokIdx += 2
		} else if currTok == "winc" {
			if tokIdx+1 >= len(toks) {
				return nil, fmt.Errorf("missing argument for winc")
			}
			wincStr := toks[tokIdx+1]
			winc, parseErr := strconv.Atoi(wincStr)
			if parseErr != nil {
				return nil, fmt.Errorf("could not parse %s as winc: %s", wincStr, parseErr)
			}
			opts.whiteIncrMs = winc
			tokIdx += 2
		} else if currTok == "binc" {
			if tokIdx+1 >= len(toks) {
				return nil, fmt.Errorf("missing argument for binc")
			}
			bincStr := toks[tokIdx+1]
			binc, parseErr := strconv.Atoi(bincStr)
			if parseErr != nil {
				return nil, fmt.Errorf("could not parse %s as winc: %s", bincStr, parseErr)
			}
			opts.blackIncrMs = binc
			tokIdx += 2
		} else if currTok == "depth" {
			if tokIdx+1 >= len(toks) {
				return nil, fmt.Errorf("missing argument for depth")
			}
			depthStr := toks[tokIdx+1]
			depth, parseErr := strconv.Atoi(depthStr)
			if parseErr != nil {
				return nil, fmt.Errorf("could not parse %s as depth: %s", depthStr, parseErr)
			}
			opts.maxDepth = uint8(depth)
			tokIdx += 2
		} else if currTok == "nodes" {
			if tokIdx+1 >= len(toks) {
				return nil, fmt.Errorf("missing argument for nodes")
			}
			nodesStr := toks[tokIdx+1]
			nodes, parseErr := strconv.Atoi(nodesStr)
			if parseErr != nil {
				return nil, fmt.Errorf("could not parse %s as nodes: %s", nodesStr, parseErr)
			}
			opts.maxNodes = nodes
			tokIdx += 2
		} else if currTok == "movetime" {
			if tokIdx+1 >= len(toks) {
				return nil, fmt.Errorf("missing argument for movetime")
			}
			movetimeStr := toks[tokIdx+1]
			movetime, parseErr := strconv.Atoi(movetimeStr)
			if parseErr != nil {
				return nil, fmt.Errorf("could not parse %s as movetime: %s", movetimeStr, parseErr)
			}
			opts.maxMs = movetime
			tokIdx += 2
		} else {
			return nil, fmt.Errorf("unknown argument: %s", currTok)
		}
	}
	return opts, nil
}

func printGoCmdHelp() {
	fmt.Println("UCI go: start search on the current internal position")
	fmt.Println("Usage:")
	fmt.Println("    go [arguments] ...")
	fmt.Println("")
	fmt.Println("The arguments are:")
	fmt.Println(Tabbed(1, Bold("searchmoves")) + " {move1} ...")
	fmt.Println(Tabbed(2, "restrict search to these moves only"))
	fmt.Println(Tabbed(2, "moves denoted in UCI-compliant 'Long-Algebraic' Notation"))
	fmt.Println(Tabbed(2, "for more information on this format, refer to:"))
	fmt.Println(Tabbed(2, "https://www.chessprogramming.org/Algebraic_Chess_Notation#LAN"))
	fmt.Println(Tabbed(2, "E.g. `go searchmoves e2e4 d2d4"))
	fmt.Println(Tabbed(1, Bold("wtime")+" {mSec}"))
	fmt.Println(Tabbed(2, "informs the engine that white has x msec left on the clock"))
	fmt.Println(Tabbed(1, Bold("btime")+" {mSec}"))
	fmt.Println(Tabbed(2, "informs the engine that black has x msec left on the clock"))
	fmt.Println(Tabbed(1, Bold("winc")+" {mSec}"))
	fmt.Println(Tabbed(2, "informs the engine that white receives x msec after each move"))
	fmt.Println(Tabbed(1, Bold("binc")+" {mSec}"))
	fmt.Println(Tabbed(2, "informs the engine that black receives x msec after each move"))
	fmt.Println(Tabbed(1, Bold("depth")+" {depth}"))
	fmt.Println(Tabbed(2, "limits the search to x plies max"))
	fmt.Println(Tabbed(1, Bold("nodes")+" {nodes}"))
	fmt.Println(Tabbed(2, "limits the search only x nodes"))
	fmt.Println(Tabbed(1, Bold("movetime")+" {mSec}"))
	fmt.Println(Tabbed(2, "requires the search to last exactly x msec"))
}
