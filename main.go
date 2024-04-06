package main

import (
	"fmt"
)

func main() {
	fmt.Println("Mila v0.0.0 - a lightweight chess AI written in go by Cameron Honis")
	NewUci(NewTranspTable()).Start()
}
