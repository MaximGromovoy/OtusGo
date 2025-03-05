package main

import (
	fmt "fmt"

	"./chessBoard"
)

func main() {
	var size uint

	fmt.Print("Input board size: ")
	fmt.Scan(&size)

	CreateAndPrintBoard(size)
}

func CreateAndPrintBoard(boardSize uint) {
	chessBoard.New(boardSize).PrintBoard()
}
