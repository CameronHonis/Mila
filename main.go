package main

import (
	"fmt"
	"github.com/CameronHonis/chess"
)

var Position = chess.GetInitBoard()
var Options *SearchOptions

func main() {
	fmt.Println("Mila v0.0.0 - a lightweight chess AI written in go by Cameron Honis")
	loopStdin(handleInput)
}
